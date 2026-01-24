package evaluator

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// PackageLockJSON represents the package-lock.json structure (lockfileVersion 3)
type PackageLockJSON struct {
	Name            string                      `json:"name"`
	Version         string                      `json:"version"`
	LockfileVersion int                         `json:"lockfileVersion"`
	Requires        bool                        `json:"requires,omitempty"`
	Packages        map[string]PackageLockEntry `json:"packages"`
}

// PackageLockEntry represents a single package entry in the lock file
type PackageLockEntry struct {
	Name             string            `json:"name,omitempty"`
	Version          string            `json:"version"`
	Resolved         string            `json:"resolved,omitempty"`
	Integrity        string            `json:"integrity,omitempty"`
	Dev              bool              `json:"dev,omitempty"`
	Optional         bool              `json:"optional,omitempty"`
	Peer             bool              `json:"peer,omitempty"`
	HasInstallScript bool              `json:"hasInstallScript,omitempty"`
	License          string            `json:"license,omitempty"`
	Dependencies     map[string]string `json:"dependencies,omitempty"`
	DevDependencies  map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies map[string]string `json:"peerDependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
	Engines          map[string]string `json:"engines,omitempty"`
	Funding          interface{}       `json:"funding,omitempty"`
	Bin              interface{}       `json:"bin,omitempty"`
	CPU              []string          `json:"cpu,omitempty"`
	OS               []string          `json:"os,omitempty"`
}

// LoadPackageLockJSON loads and parses the package-lock.json file
func LoadPackageLockJSON(path string) (*PackageLockJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty lock file if it doesn't exist
			return NewPackageLockJSON("", ""), nil
		}
		return nil, fmt.Errorf("failed to read package-lock.json: %w", err)
	}

	var lock PackageLockJSON
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse package-lock.json: %w", err)
	}

	// Ensure packages map exists
	if lock.Packages == nil {
		lock.Packages = make(map[string]PackageLockEntry)
	}

	return &lock, nil
}

// SavePackageLockJSON saves the lock file to package-lock.json
func SavePackageLockJSON(path string, lock *PackageLockJSON) error {
	// Ensure lockfileVersion is set
	if lock.LockfileVersion == 0 {
		lock.LockfileVersion = 3
	}

	// Sort packages for deterministic output
	sortedPackages := make(map[string]PackageLockEntry)
	keys := make([]string, 0, len(lock.Packages))
	for k := range lock.Packages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sortedPackages[k] = lock.Packages[k]
	}
	lock.Packages = sortedPackages

	// Pretty print with 2-space indentation
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package-lock.json: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	return os.WriteFile(path, data, 0644)
}

// NewPackageLockJSON creates a new empty PackageLockJSON
func NewPackageLockJSON(name, version string) *PackageLockJSON {
	lock := &PackageLockJSON{
		Name:            name,
		Version:         version,
		LockfileVersion: 3,
		Requires:        true,
		Packages:        make(map[string]PackageLockEntry),
	}

	// Add root package entry (empty key represents the root)
	lock.Packages[""] = PackageLockEntry{
		Name:    name,
		Version: version,
	}

	return lock
}

// GetPackageEntry gets a package entry by its path in node_modules
// path format: "node_modules/package-name" or "node_modules/scope/package-name"
func (l *PackageLockJSON) GetPackageEntry(pkgPath string) (*PackageLockEntry, bool) {
	if entry, ok := l.Packages[pkgPath]; ok {
		return &entry, true
	}
	return nil, false
}

// GetPackageByName gets a package entry by package name
func (l *PackageLockJSON) GetPackageByName(name string) (*PackageLockEntry, bool) {
	path := "node_modules/" + name
	return l.GetPackageEntry(path)
}

// SetPackageEntry sets or updates a package entry
func (l *PackageLockJSON) SetPackageEntry(pkgPath string, entry PackageLockEntry) {
	if l.Packages == nil {
		l.Packages = make(map[string]PackageLockEntry)
	}
	l.Packages[pkgPath] = entry
}

// AddPackage adds a new package to the lock file
func (l *PackageLockJSON) AddPackage(name, version, resolved, integrity string, isDev bool) {
	path := "node_modules/" + name
	entry := PackageLockEntry{
		Version:   version,
		Resolved:  resolved,
		Integrity: integrity,
		Dev:       isDev,
	}
	l.SetPackageEntry(path, entry)
}

// RemovePackage removes a package from the lock file
func (l *PackageLockJSON) RemovePackage(name string) bool {
	path := "node_modules/" + name
	if _, ok := l.Packages[path]; ok {
		delete(l.Packages, path)
		return true
	}
	return false
}

// UpdateRootDependencies updates the root package entry with dependencies
func (l *PackageLockJSON) UpdateRootDependencies(deps, devDeps map[string]string) {
	root := l.Packages[""]
	root.Dependencies = deps
	root.DevDependencies = devDeps
	l.Packages[""] = root
}

// GetInstalledPackages returns a list of all installed package names
func (l *PackageLockJSON) GetInstalledPackages() []string {
	var packages []string
	for path := range l.Packages {
		if path == "" {
			continue // Skip root entry
		}
		// Extract package name from path (node_modules/package-name)
		name := filepath.Base(path)
		// Handle scoped packages (node_modules/@scope/package-name)
		dir := filepath.Dir(path)
		if filepath.Base(dir) != "node_modules" {
			scope := filepath.Base(dir)
			name = scope + "/" + name
		}
		packages = append(packages, name)
	}
	sort.Strings(packages)
	return packages
}

// IsPackageInstalled checks if a package is in the lock file
func (l *PackageLockJSON) IsPackageInstalled(name string) bool {
	_, ok := l.GetPackageByName(name)
	return ok
}

// GetPackageVersion returns the installed version of a package
func (l *PackageLockJSON) GetPackageVersion(name string) (string, bool) {
	entry, ok := l.GetPackageByName(name)
	if !ok {
		return "", false
	}
	return entry.Version, true
}

// CalculateIntegrity calculates the SHA-512 integrity hash for content
func CalculateIntegrity(content []byte) string {
	hash := sha512.Sum512(content)
	encoded := base64.StdEncoding.EncodeToString(hash[:])
	return "sha512-" + encoded
}

// VerifyIntegrity verifies content against an integrity hash
func VerifyIntegrity(content []byte, integrity string) bool {
	if len(integrity) < 7 {
		return false
	}
	
	// Only support sha512 for now
	if integrity[:7] != "sha512-" {
		return false
	}
	
	expected := integrity[7:]
	hash := sha512.Sum512(content)
	actual := base64.StdEncoding.EncodeToString(hash[:])
	
	return actual == expected
}

// PackageLockMetadata holds metadata about the lock file
type PackageLockMetadata struct {
	GeneratedAt   time.Time
	NuVersion     string
	PackageCount  int
}

// GetMetadata returns metadata about the lock file
func (l *PackageLockJSON) GetMetadata() PackageLockMetadata {
	count := 0
	for path := range l.Packages {
		if path != "" {
			count++
		}
	}
	
	return PackageLockMetadata{
		PackageCount: count,
	}
}

// Diff compares two lock files and returns the differences
type LockFileDiff struct {
	Added   []string // Packages added
	Removed []string // Packages removed
	Changed []string // Packages with version changes
}

// Diff compares this lock file with another
func (l *PackageLockJSON) Diff(other *PackageLockJSON) LockFileDiff {
	diff := LockFileDiff{
		Added:   []string{},
		Removed: []string{},
		Changed: []string{},
	}
	
	current := l.GetInstalledPackages()
	previous := other.GetInstalledPackages()
	
	currentSet := make(map[string]bool)
	for _, pkg := range current {
		currentSet[pkg] = true
	}
	
	previousSet := make(map[string]bool)
	for _, pkg := range previous {
		previousSet[pkg] = true
	}
	
	// Find added and changed
	for _, pkg := range current {
		if !previousSet[pkg] {
			diff.Added = append(diff.Added, pkg)
		} else {
			// Check if version changed
			currentVersion, _ := l.GetPackageVersion(pkg)
			previousVersion, _ := other.GetPackageVersion(pkg)
			if currentVersion != previousVersion {
				diff.Changed = append(diff.Changed, pkg)
			}
		}
	}
	
	// Find removed
	for _, pkg := range previous {
		if !currentSet[pkg] {
			diff.Removed = append(diff.Removed, pkg)
		}
	}
	
	return diff
}

// Clone creates a deep copy of the lock file
func (l *PackageLockJSON) Clone() *PackageLockJSON {
	clone := &PackageLockJSON{
		Name:            l.Name,
		Version:         l.Version,
		LockfileVersion: l.LockfileVersion,
		Requires:        l.Requires,
		Packages:        make(map[string]PackageLockEntry),
	}
	
	for path, entry := range l.Packages {
		entryCopy := entry
		
		// Deep copy maps
		if entry.Dependencies != nil {
			entryCopy.Dependencies = make(map[string]string)
			for k, v := range entry.Dependencies {
				entryCopy.Dependencies[k] = v
			}
		}
		if entry.DevDependencies != nil {
			entryCopy.DevDependencies = make(map[string]string)
			for k, v := range entry.DevDependencies {
				entryCopy.DevDependencies[k] = v
			}
		}
		if entry.PeerDependencies != nil {
			entryCopy.PeerDependencies = make(map[string]string)
			for k, v := range entry.PeerDependencies {
				entryCopy.PeerDependencies[k] = v
			}
		}
		if entry.Engines != nil {
			entryCopy.Engines = make(map[string]string)
			for k, v := range entry.Engines {
				entryCopy.Engines[k] = v
			}
		}
		
		clone.Packages[path] = entryCopy
	}
	
	return clone
}
