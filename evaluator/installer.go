package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Installer handles package installation
type Installer struct {
	Registry   *NPMRegistry
	ProjectDir string
	ModulesDir string
	Verbose    bool
	mu         sync.Mutex
}

// NewInstaller creates a new package installer
func NewInstaller(projectDir string) *Installer {
	return &Installer{
		Registry:   NewNPMRegistry(),
		ProjectDir: projectDir,
		ModulesDir: filepath.Join(projectDir, "node_modules"),
		Verbose:    true,
	}
}

// InstallResult represents the result of an installation
type InstallResult struct {
	Installed   []string
	Updated     []string
	Removed     []string
	Errors      []error
	TotalCount  int
}

// InstallAll installs all dependencies from package.json
func (i *Installer) InstallAll() (*InstallResult, error) {
	result := &InstallResult{}

	pkgPath := filepath.Join(i.ProjectDir, "package.json")
	lockPath := filepath.Join(i.ProjectDir, "package-lock.json")

	// Load package.json
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package.json: %w", err)
	}

	// Load or create package-lock.json
	lock, err := LoadPackageLockJSON(lockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package-lock.json: %w", err)
	}

	// Update lock file metadata
	lock.Name = pkg.Name
	lock.Version = pkg.Version

	// Create node_modules directory
	if err := os.MkdirAll(i.ModulesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create node_modules: %w", err)
	}

	// Collect all dependencies
	allDeps := make(map[string]string)
	devDeps := make(map[string]bool)

	for name, version := range pkg.Dependencies {
		allDeps[name] = version
	}
	for name, version := range pkg.DevDependencies {
		allDeps[name] = version
		devDeps[name] = true
	}

	result.TotalCount = len(allDeps)

	if i.Verbose {
		fmt.Printf("Installing %d packages...\n", len(allDeps))
	}

	// Install each dependency
	for name, versionRange := range allDeps {
		isDev := devDeps[name]
		
		if i.Verbose {
			fmt.Printf("  Installing %s@%s...\n", name, versionRange)
		}

		err := i.installPackage(name, versionRange, isDev, lock)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", name, err))
			fmt.Printf("  ✗ Failed to install %s: %s\n", name, err)
		} else {
			result.Installed = append(result.Installed, name)
			if i.Verbose {
				ver, _ := lock.GetPackageVersion(name)
				fmt.Printf("  ✓ Installed %s@%s\n", name, ver)
			}
		}
	}

	// Update root dependencies in lock file
	lock.UpdateRootDependencies(pkg.Dependencies, pkg.DevDependencies)

	// Save package-lock.json
	if err := SavePackageLockJSON(lockPath, lock); err != nil {
		return result, fmt.Errorf("failed to save package-lock.json: %w", err)
	}

	if i.Verbose {
		fmt.Printf("\n✓ Installed %d packages\n", len(result.Installed))
		if len(result.Errors) > 0 {
			fmt.Printf("✗ %d packages failed\n", len(result.Errors))
		}
	}

	return result, nil
}

// installPackage installs a single package
func (i *Installer) installPackage(name, versionRange string, isDev bool, lock *PackageLockJSON) error {
	// Resolve version
	version, err := i.Registry.ResolveVersion(name, versionRange)
	if err != nil {
		return fmt.Errorf("failed to resolve version: %w", err)
	}

	// Check if already installed with correct version
	if entry, ok := lock.GetPackageByName(name); ok {
		if entry.Version == version {
			// Already installed
			return nil
		}
	}

	// Download and install package
	destPath := filepath.Join(i.ModulesDir, name)
	
	// Remove existing if upgrading
	os.RemoveAll(destPath)

	info, err := i.Registry.DownloadPackage(name, version, destPath)
	if err != nil {
		return err
	}

	// Update lock file
	lock.AddPackage(name, version, info.Dist.Tarball, info.Dist.Integrity, isDev)

	// Install sub-dependencies
	for depName, depRange := range info.Dependencies {
		if err := i.installPackage(depName, depRange, false, lock); err != nil {
			// Continue on sub-dependency errors
			fmt.Printf("    ⚠ Warning: failed to install %s: %s\n", depName, err)
		}
	}

	return nil
}

// InstallPackage installs a specific package and updates package.json
func (i *Installer) InstallPackage(name string, version string, isDev bool) error {
	pkgPath := filepath.Join(i.ProjectDir, "package.json")
	lockPath := filepath.Join(i.ProjectDir, "package-lock.json")

	// Load or create package.json
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		// package.json doesn't exist, create it
		projectName := filepath.Base(i.ProjectDir)
		if projectName == "." || projectName == "/" {
			projectName = "my-project"
		}
		// Sanitize name for npm
		projectName = strings.ToLower(projectName)
		projectName = strings.ReplaceAll(projectName, " ", "-")
		
		pkg = NewPackageJSON(projectName)
		fmt.Printf("Creating package.json for %s...\n", projectName)
	}

	// Load package-lock.json (creates empty one if doesn't exist)
	lock, err := LoadPackageLockJSON(lockPath)
	if err != nil {
		lock = NewPackageLockJSON(pkg.Name, pkg.Version)
	}

	// Resolve version if not specified or is a range
	resolvedVersion := version
	if version == "" || version == "latest" {
		resolvedVersion, err = i.Registry.GetLatestVersion(name)
		if err != nil {
			return err
		}
		version = "^" + resolvedVersion
	} else if !strings.HasPrefix(version, "^") && !strings.HasPrefix(version, "~") && !strings.HasPrefix(version, ">") && !strings.HasPrefix(version, "<") {
		// If exact version specified, resolve it
		info, err := i.Registry.GetPackageVersion(name, version)
		if err != nil {
			return err
		}
		resolvedVersion = info.Version
		version = "^" + resolvedVersion
	} else {
		resolvedVersion, err = i.Registry.ResolveVersion(name, version)
		if err != nil {
			return err
		}
	}

	// Create node_modules if needed
	if err := os.MkdirAll(i.ModulesDir, 0755); err != nil {
		return err
	}

	fmt.Printf("Installing %s@%s...\n", name, resolvedVersion)

	// Install the package
	destPath := filepath.Join(i.ModulesDir, name)
	os.RemoveAll(destPath)

	info, err := i.Registry.DownloadPackage(name, resolvedVersion, destPath)
	if err != nil {
		return err
	}

	// Update package.json
	if isDev {
		pkg.AddDevDependency(name, version)
	} else {
		pkg.AddDependency(name, version)
	}

	// Update lock file
	lock.Name = pkg.Name
	lock.Version = pkg.Version
	lock.AddPackage(name, resolvedVersion, info.Dist.Tarball, info.Dist.Integrity, isDev)
	lock.UpdateRootDependencies(pkg.Dependencies, pkg.DevDependencies)

	// Install sub-dependencies
	for depName, depRange := range info.Dependencies {
		if err := i.installPackage(depName, depRange, false, lock); err != nil {
			fmt.Printf("  ⚠ Warning: failed to install %s: %s\n", depName, err)
		}
	}

	// Save files
	if err := SavePackageJSON(pkgPath, pkg); err != nil {
		return err
	}
	if err := SavePackageLockJSON(lockPath, lock); err != nil {
		return err
	}

	fmt.Printf("✓ Added %s@%s\n", name, resolvedVersion)
	return nil
}

// UninstallPackage removes a package
func (i *Installer) UninstallPackage(name string) error {
	pkgPath := filepath.Join(i.ProjectDir, "package.json")
	lockPath := filepath.Join(i.ProjectDir, "package-lock.json")

	// Load package.json
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		return err
	}

	// Remove from package.json
	if !pkg.RemoveDependency(name) {
		return fmt.Errorf("package '%s' is not installed", name)
	}

	// Load and update lock file
	lock, err := LoadPackageLockJSON(lockPath)
	if err == nil {
		lock.RemovePackage(name)
		lock.UpdateRootDependencies(pkg.Dependencies, pkg.DevDependencies)
		SavePackageLockJSON(lockPath, lock)
	}

	// Remove from node_modules
	pkgDir := filepath.Join(i.ModulesDir, name)
	if err := os.RemoveAll(pkgDir); err != nil {
		return fmt.Errorf("failed to remove package directory: %w", err)
	}

	// Save package.json
	if err := SavePackageJSON(pkgPath, pkg); err != nil {
		return err
	}

	fmt.Printf("✓ Removed %s\n", name)
	return nil
}

// ListInstalled lists all installed packages
func (i *Installer) ListInstalled() ([]InstalledPackage, error) {
	lockPath := filepath.Join(i.ProjectDir, "package-lock.json")
	
	lock, err := LoadPackageLockJSON(lockPath)
	if err != nil {
		return nil, err
	}

	pkgPath := filepath.Join(i.ProjectDir, "package.json")
	pkg, _ := LoadPackageJSON(pkgPath)

	var packages []InstalledPackage
	
	for path, entry := range lock.Packages {
		if path == "" {
			continue // Skip root
		}
		
		name := strings.TrimPrefix(path, "node_modules/")
		
		dep := InstalledPackage{
			Name:    name,
			Version: entry.Version,
			Dev:     entry.Dev,
		}
		
		// Check if it's a direct dependency
		if pkg != nil {
			if _, ok := pkg.Dependencies[name]; ok {
				dep.Direct = true
			}
			if _, ok := pkg.DevDependencies[name]; ok {
				dep.Direct = true
				dep.Dev = true
			}
		}
		
		packages = append(packages, dep)
	}

	// Sort by name
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages, nil
}

// InstalledPackage represents an installed package
type InstalledPackage struct {
	Name    string
	Version string
	Dev     bool
	Direct  bool
}

// Prune removes packages not listed in package.json
func (i *Installer) Prune() ([]string, error) {
	pkgPath := filepath.Join(i.ProjectDir, "package.json")
	
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		return nil, err
	}

	// Get all packages that should be installed
	needed := make(map[string]bool)
	for name := range pkg.Dependencies {
		needed[name] = true
	}
	for name := range pkg.DevDependencies {
		needed[name] = true
	}

	// Check what's installed
	entries, err := os.ReadDir(i.ModulesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var removed []string
	for _, entry := range entries {
		name := entry.Name()
		
		// Skip hidden files and scoped packages (handled separately)
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "@") {
			continue
		}

		if !needed[name] {
			pkgDir := filepath.Join(i.ModulesDir, name)
			if err := os.RemoveAll(pkgDir); err == nil {
				removed = append(removed, name)
			}
		}
	}

	return removed, nil
}

// Update updates packages to their latest versions
func (i *Installer) Update(packages []string) error {
	pkgPath := filepath.Join(i.ProjectDir, "package.json")
	
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		return err
	}

	// If no packages specified, update all
	if len(packages) == 0 {
		for name := range pkg.Dependencies {
			packages = append(packages, name)
		}
		for name := range pkg.DevDependencies {
			packages = append(packages, name)
		}
	}

	for _, name := range packages {
		// Get current version range
		versionRange := pkg.Dependencies[name]
		if versionRange == "" {
			versionRange = pkg.DevDependencies[name]
		}
		if versionRange == "" {
			fmt.Printf("  ⚠ Package %s not found in dependencies\n", name)
			continue
		}

		// Get latest version that satisfies the range
		latest, err := i.Registry.ResolveVersion(name, versionRange)
		if err != nil {
			fmt.Printf("  ⚠ Failed to resolve %s: %s\n", name, err)
			continue
		}

		// Check if update is needed
		installed, err := i.ListInstalled()
		if err == nil {
			for _, inst := range installed {
				if inst.Name == name && inst.Version == latest {
					fmt.Printf("  %s@%s is already up to date\n", name, latest)
					continue
				}
			}
		}

		// Reinstall
		isDev := pkg.DevDependencies[name] != ""
		if err := i.InstallPackage(name, latest, isDev); err != nil {
			fmt.Printf("  ✗ Failed to update %s: %s\n", name, err)
		}
	}

	return nil
}

// CreateBinLinks creates symlinks in node_modules/.bin for package binaries
func (i *Installer) CreateBinLinks() error {
	binDir := filepath.Join(i.ModulesDir, ".bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	// Scan installed packages for bin entries
	entries, err := os.ReadDir(i.ModulesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		pkgDir := filepath.Join(i.ModulesDir, entry.Name())
		pkgJsonPath := filepath.Join(pkgDir, "package.json")

		pkg, err := LoadPackageJSON(pkgJsonPath)
		if err != nil {
			continue
		}

		bins := pkg.GetBin()
		for binName, binPath := range bins {
			source := filepath.Join(pkgDir, binPath)
			target := filepath.Join(binDir, binName)

			// Remove existing link
			os.Remove(target)

			// Create symlink
			if err := os.Symlink(source, target); err != nil {
				continue
			}

			// Make executable
			os.Chmod(source, 0755)
		}
	}

	return nil
}
