package evaluator

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NPMRegistry interacts with npm registry
type NPMRegistry struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
	UserAgent  string
}

// PackageMetadata from npm registry
type PackageMetadata struct {
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	DistTags    map[string]string             `json:"dist-tags"`
	Versions    map[string]PackageVersionInfo `json:"versions"`
	Time        map[string]string             `json:"time"`
	Readme      string                        `json:"readme,omitempty"`
	Homepage    string                        `json:"homepage,omitempty"`
	Keywords    []string                      `json:"keywords,omitempty"`
	Repository  interface{}                   `json:"repository,omitempty"`
	Author      interface{}                   `json:"author,omitempty"`
	License     string                        `json:"license,omitempty"`
}

// PackageVersionInfo for a specific version
type PackageVersionInfo struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description,omitempty"`
	Main            string            `json:"main,omitempty"`
	Module          string            `json:"module,omitempty"`
	Types           string            `json:"types,omitempty"`
	Exports         interface{}       `json:"exports,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies map[string]string `json:"peerDependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
	Engines         interface{}       `json:"engines,omitempty"`         // Can be map or array
	Bin             interface{}       `json:"bin,omitempty"`
	Scripts         interface{}       `json:"scripts,omitempty"`         // Can be map or null
	Dist            PackageDist       `json:"dist"`
	License         interface{}       `json:"license,omitempty"`         // Can be string or object
	Deprecated      string            `json:"deprecated,omitempty"`
	HasInstallScript bool             `json:"hasInstallScript,omitempty"`
}

// PackageDist contains distribution info
type PackageDist struct {
	Shasum        string `json:"shasum"`
	Tarball       string `json:"tarball"`
	Integrity     string `json:"integrity"`
	FileCount     int    `json:"fileCount,omitempty"`
	UnpackedSize  int64  `json:"unpackedSize,omitempty"`
	NpmSignature  string `json:"npm-signature,omitempty"`
}

// NewNPMRegistry creates a new NPM registry client
func NewNPMRegistry() *NPMRegistry {
	return &NPMRegistry{
		BaseURL:   "https://registry.npmjs.org",
		UserAgent: "nulang-package-manager/1.0.0",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewNPMRegistryWithToken creates a registry client with authentication
func NewNPMRegistryWithToken(token string) *NPMRegistry {
	registry := NewNPMRegistry()
	registry.Token = token
	return registry
}

// doRequest performs an HTTP request with proper headers
func (r *NPMRegistry) doRequest(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", r.UserAgent)
	req.Header.Set("Accept", "application/json")
	
	if r.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Token)
	}

	return r.HTTPClient.Do(req)
}

// GetPackageMetadata fetches complete package info from registry
func (r *NPMRegistry) GetPackageMetadata(name string) (*PackageMetadata, error) {
	// Handle scoped packages (@scope/package)
	encodedName := name
	if strings.HasPrefix(name, "@") {
		encodedName = strings.Replace(name, "/", "%2F", 1)
	}

	url := fmt.Sprintf("%s/%s", r.BaseURL, encodedName)
	
	resp, err := r.doRequest("GET", url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch package %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("package '%s' not found in registry", name)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("registry returned status %d for package %s", resp.StatusCode, name)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var metadata PackageMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse package metadata: %w", err)
	}

	return &metadata, nil
}

// GetPackageVersion fetches info for a specific version
func (r *NPMRegistry) GetPackageVersion(name, version string) (*PackageVersionInfo, error) {
	// Handle scoped packages
	encodedName := name
	if strings.HasPrefix(name, "@") {
		encodedName = strings.Replace(name, "/", "%2F", 1)
	}

	url := fmt.Sprintf("%s/%s/%s", r.BaseURL, encodedName, version)
	
	resp, err := r.doRequest("GET", url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("version '%s' not found for package '%s'", version, name)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var info PackageVersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse version info: %w", err)
	}

	return &info, nil
}

// GetLatestVersion returns the latest version tag
func (r *NPMRegistry) GetLatestVersion(name string) (string, error) {
	metadata, err := r.GetPackageMetadata(name)
	if err != nil {
		return "", err
	}

	if latest, ok := metadata.DistTags["latest"]; ok {
		return latest, nil
	}

	return "", fmt.Errorf("no latest version found for %s", name)
}

// ResolveVersion resolves a version range to a concrete version
func (r *NPMRegistry) ResolveVersion(name, versionRange string) (string, error) {
	// Handle special cases
	if versionRange == "latest" || versionRange == "" || versionRange == "*" {
		return r.GetLatestVersion(name)
	}

	// Parse the version range
	semverRange, err := ParseSemverRange(versionRange)
	if err != nil {
		return "", fmt.Errorf("invalid version range: %w", err)
	}

	// Get all available versions
	metadata, err := r.GetPackageMetadata(name)
	if err != nil {
		return "", err
	}

	// Collect all versions
	versions := make([]string, 0, len(metadata.Versions))
	for v := range metadata.Versions {
		versions = append(versions, v)
	}

	// Find the best matching version
	resolved := semverRange.MaxSatisfying(versions)
	if resolved == "" {
		return "", fmt.Errorf("no version satisfies range '%s' for package '%s'", versionRange, name)
	}

	return resolved, nil
}

// DownloadPackage downloads and extracts a package tarball
func (r *NPMRegistry) DownloadPackage(name, version, destPath string) (*PackageVersionInfo, error) {
	// Get version info
	info, err := r.GetPackageVersion(name, version)
	if err != nil {
		return nil, err
	}

	if info.Dist.Tarball == "" {
		return nil, fmt.Errorf("no tarball URL for %s@%s", name, version)
	}

	// Download tarball
	resp, err := r.doRequest("GET", info.Dist.Tarball)
	if err != nil {
		return nil, fmt.Errorf("failed to download tarball: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download tarball: HTTP %d", resp.StatusCode)
	}

	// Read tarball content
	tarballContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tarball: %w", err)
	}

	// Verify integrity if available
	if info.Dist.Integrity != "" {
		if !VerifyIntegrity(tarballContent, info.Dist.Integrity) {
			return nil, fmt.Errorf("integrity check failed for %s@%s", name, version)
		}
	}

	// Extract tarball
	if err := extractTarball(tarballContent, destPath); err != nil {
		return nil, fmt.Errorf("failed to extract tarball: %w", err)
	}

	return info, nil
}

// extractTarball extracts a gzipped tarball to the destination path
func extractTarball(content []byte, destPath string) error {
	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	// Create gzip reader
	gzr, err := gzip.NewReader(strings.NewReader(string(content)))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Create tar reader
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// npm tarballs have a 'package/' prefix that we need to strip
		name := header.Name
		if strings.HasPrefix(name, "package/") {
			name = name[8:]
		} else if strings.HasPrefix(name, "package") {
			name = name[7:]
			if strings.HasPrefix(name, "/") {
				name = name[1:]
			}
		}

		if name == "" {
			continue
		}

		target := filepath.Join(destPath, name)

		// Ensure the target is within destPath (security check)
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destPath)) {
			return fmt.Errorf("invalid file path: %s", name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			// Create file
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// Copy content
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		case tar.TypeSymlink:
			// Handle symlinks
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				// Ignore symlink errors on Windows or permission issues
				continue
			}
		}
	}

	return nil
}

// GetAllVersions returns all available versions for a package
func (r *NPMRegistry) GetAllVersions(name string) ([]string, error) {
	metadata, err := r.GetPackageMetadata(name)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(metadata.Versions))
	for v := range metadata.Versions {
		versions = append(versions, v)
	}

	return versions, nil
}

// SearchPackages searches for packages (basic implementation)
func (r *NPMRegistry) SearchPackages(query string, limit int) ([]PackageSearchResult, error) {
	url := fmt.Sprintf("%s/-/v1/search?text=%s&size=%d", r.BaseURL, query, limit)

	resp, err := r.doRequest("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Objects []struct {
			Package PackageSearchResult `json:"package"`
		} `json:"objects"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	packages := make([]PackageSearchResult, len(result.Objects))
	for i, obj := range result.Objects {
		packages[i] = obj.Package
	}

	return packages, nil
}

// PackageSearchResult represents a search result
type PackageSearchResult struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Keywords    []string `json:"keywords,omitempty"`
	Publisher   struct {
		Username string `json:"username"`
	} `json:"publisher,omitempty"`
	Date string `json:"date,omitempty"`
}

// CheckPackageExists checks if a package exists in the registry
func (r *NPMRegistry) CheckPackageExists(name string) bool {
	_, err := r.GetLatestVersion(name)
	return err == nil
}

// GetDependencies returns all dependencies for a package version
func (r *NPMRegistry) GetDependencies(name, version string) (map[string]string, error) {
	info, err := r.GetPackageVersion(name, version)
	if err != nil {
		return nil, err
	}

	return info.Dependencies, nil
}

// ResolveDependencyTree resolves the complete dependency tree
type DependencyNode struct {
	Name         string
	Version      string
	Resolved     string
	Integrity    string
	Dependencies map[string]*DependencyNode
	Dev          bool
	Optional     bool
	Peer         bool
}

// ResolveDependencyTree resolves all dependencies recursively
func (r *NPMRegistry) ResolveDependencyTree(name, versionRange string, visited map[string]bool) (*DependencyNode, error) {
	if visited == nil {
		visited = make(map[string]bool)
	}

	// Prevent circular dependencies
	key := name + "@" + versionRange
	if visited[key] {
		return nil, nil
	}
	visited[key] = true

	// Resolve version
	version, err := r.ResolveVersion(name, versionRange)
	if err != nil {
		return nil, err
	}

	// Get package info
	info, err := r.GetPackageVersion(name, version)
	if err != nil {
		return nil, err
	}

	node := &DependencyNode{
		Name:         name,
		Version:      version,
		Resolved:     info.Dist.Tarball,
		Integrity:    info.Dist.Integrity,
		Dependencies: make(map[string]*DependencyNode),
	}

	// Resolve sub-dependencies
	for depName, depRange := range info.Dependencies {
		depNode, err := r.ResolveDependencyTree(depName, depRange, visited)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s: %w", depName, err)
		}
		if depNode != nil {
			node.Dependencies[depName] = depNode
		}
	}

	return node, nil
}
