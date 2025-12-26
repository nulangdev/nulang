package evaluator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// PackageConfig represents the nulang.yml configuration
type PackageConfig struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	Description  string            `yaml:"description,omitempty"`
	Main         string            `yaml:"main,omitempty"`
	Dependencies map[string]string `yaml:"dependencies,omitempty"`
}

// LockEntry represents a single package entry in the lock file
type LockEntry struct {
	Name     string    `yaml:"name"`
	URL      string    `yaml:"url"`
	Commit   string    `yaml:"commit"`
	Checksum string    `yaml:"checksum"`
	Resolved time.Time `yaml:"resolved"`
}

// LockFile represents the nulang.lock file structure
type LockFile struct {
	Version  int         `yaml:"version"`
	Packages []LockEntry `yaml:"packages"`
}

// LoadPackageConfig loads and parses the nulang.yml file
func LoadPackageConfig(path string) (*PackageConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read package config: %w", err)
	}

	var config PackageConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse package config: %w", err)
	}

	// Set defaults
	if config.Main == "" {
		config.Main = "index.nu"
	}

	return &config, nil
}

// SavePackageConfig saves the package configuration to nulang.yml
func SavePackageConfig(path string, config *PackageConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal package config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// LoadLockFile loads and parses the nulang.lock file
func LoadLockFile(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &LockFile{Version: 1, Packages: []LockEntry{}}, nil
		}
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock LockFile
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lock, nil
}

// SaveLockFile saves the lock file to nulang.lock
func SaveLockFile(path string, lock *LockFile) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// FindLockEntry finds a package in the lock file
func (lf *LockFile) FindLockEntry(name string) *LockEntry {
	for i := range lf.Packages {
		if lf.Packages[i].Name == name {
			return &lf.Packages[i]
		}
	}
	return nil
}

// UpdateOrAddEntry updates an existing entry or adds a new one
func (lf *LockFile) UpdateOrAddEntry(entry LockEntry) {
	for i := range lf.Packages {
		if lf.Packages[i].Name == entry.Name {
			lf.Packages[i] = entry
			return
		}
	}
	lf.Packages = append(lf.Packages, entry)
}

// InstallDependencies installs all dependencies from nulang.yml
func InstallDependencies(projectPath string) error {
	configPath := filepath.Join(projectPath, "nulang.yml")
	lockPath := filepath.Join(projectPath, "nulang.lock")
	modulesPath := filepath.Join(projectPath, ".nu_modules")

	// Load package config
	config, err := LoadPackageConfig(configPath)
	if err != nil {
		return err
	}

	// Load or create lock file
	lock, err := LoadLockFile(lockPath)
	if err != nil {
		return err
	}

	// Create .nu_modules directory
	if err := os.MkdirAll(modulesPath, 0755); err != nil {
		return fmt.Errorf("failed to create .nu_modules directory: %w", err)
	}

	// Install each dependency
	for name, url := range config.Dependencies {
		fmt.Printf("Installing %s from %s...\n", name, url)

		pkgDir := filepath.Join(modulesPath, name)
		
		// Download the package
		commit, checksum, err := DownloadFromGitHub(url, pkgDir)
		if err != nil {
			return fmt.Errorf("failed to install %s: %w", name, err)
		}

		// Update lock file
		lock.UpdateOrAddEntry(LockEntry{
			Name:     name,
			URL:      url,
			Commit:   commit,
			Checksum: checksum,
			Resolved: time.Now(),
		})

		fmt.Printf("  ✓ Installed %s@%s\n", name, commit[:7])
	}

	// Save lock file
	if err := SaveLockFile(lockPath, lock); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	fmt.Printf("\n✓ Installed %d packages\n", len(config.Dependencies))
	return nil
}

// DownloadFromGitHub downloads a package from GitHub
// URL format: github.com/user/repo
// It first downloads nulang.yml to find the main file, then downloads that file
func DownloadFromGitHub(url, destPath string) (commit string, checksum string, err error) {
	// Clean URL (remove https:// if present)
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Parse GitHub URL: github.com/user/repo
	parts := strings.Split(url, "/")
	if len(parts) < 3 || parts[0] != "github.com" {
		return "", "", fmt.Errorf("invalid GitHub URL: %s (expected github.com/user/repo)", url)
	}

	user := parts[1]
	repo := parts[2]

	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create package directory: %w", err)
	}

	// First, download nulang.yml to find the main file
	mainFile := "index.nu" // default fallback
	configURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/nulang.yml", user, repo)
	
	configResp, err := http.Get(configURL)
	if err == nil && configResp.StatusCode == 200 {
		configContent, err := io.ReadAll(configResp.Body)
		configResp.Body.Close()
		
		if err == nil {
			// Save nulang.yml locally
			configPath := filepath.Join(destPath, "nulang.yml")
			os.WriteFile(configPath, configContent, 0644)
			
			// Parse to get main file
			var pkgConfig PackageConfig
			if yaml.Unmarshal(configContent, &pkgConfig) == nil && pkgConfig.Main != "" {
				mainFile = pkgConfig.Main
			}
		}
	} else if configResp != nil {
		configResp.Body.Close()
	}

	// Download the main file
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", user, repo, mainFile)
	mainPath := filepath.Join(destPath, mainFile)

	resp, err := http.Get(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to download from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to download %s: HTTP %d", rawURL, resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	if err := os.WriteFile(mainPath, content, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write %s: %w", mainFile, err)
	}

	// If main file is not index.nu, create index.nu that re-exports it
	if mainFile != "index.nu" {
		indexContent := fmt.Sprintf("// Auto-generated index that re-exports main file\nexport * from \"./%s\"\n", strings.TrimSuffix(mainFile, ".nu"))
		indexPath := filepath.Join(destPath, "index.nu")
		os.WriteFile(indexPath, []byte(indexContent), 0644)
	}

	// Get latest commit SHA from API
	commit, err = getLatestCommit(user, repo)
	if err != nil {
		// Use timestamp as fallback commit
		commit = fmt.Sprintf("snapshot-%d", time.Now().Unix())
	}

	// Calculate checksum
	hash := sha256.Sum256(content)
	checksum = "sha256:" + hex.EncodeToString(hash[:])

	return commit, checksum, nil
}

// getLatestCommit gets the latest commit SHA for a repo
func getLatestCommit(user, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/main", user, repo)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "nulang-package-manager")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	// Simple JSON parsing for sha field
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Find "sha": "..." in the response
	bodyStr := string(body)
	shaIdx := strings.Index(bodyStr, `"sha":"`)
	if shaIdx == -1 {
		return "", fmt.Errorf("sha not found in response")
	}

	shaStart := shaIdx + 7
	shaEnd := strings.Index(bodyStr[shaStart:], `"`)
	if shaEnd == -1 {
		return "", fmt.Errorf("invalid sha format")
	}

	return bodyStr[shaStart : shaStart+shaEnd], nil
}

// InitProject creates a new nulang.yml file
func InitProject(projectPath string) error {
	configPath := filepath.Join(projectPath, "nulang.yml")

	// Check if already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("nulang.yml already exists")
	}

	// Get project name from directory
	name := filepath.Base(projectPath)
	if name == "." || name == "/" {
		name = "my-project"
	}

	config := PackageConfig{
		Name:         name,
		Version:      "1.0.0",
		Main:         "index.nu",
		Dependencies: make(map[string]string),
	}

	if err := SavePackageConfig(configPath, &config); err != nil {
		return err
	}

	fmt.Printf("Created nulang.yml for %s\n", name)
	return nil
}

// FindNuModulesPath finds the .nu_modules directory by traversing up
func FindNuModulesPath(startPath string) string {
	dir := startPath
	if !isDirectory(dir) {
		dir = filepath.Dir(dir)
	}

	for {
		modulesPath := filepath.Join(dir, ".nu_modules")
		if info, err := os.Stat(modulesPath); err == nil && info.IsDir() {
			return modulesPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
