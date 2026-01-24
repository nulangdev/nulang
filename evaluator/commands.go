package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CLI command handlers for the package manager

// InitProjectJSON creates a new package.json file (interactive)
func InitProjectJSON(projectPath string) error {
	pkgPath := filepath.Join(projectPath, "package.json")

	// Check if already exists
	if _, err := os.Stat(pkgPath); err == nil {
		return fmt.Errorf("package.json already exists")
	}

	// Get project name from directory
	name := filepath.Base(projectPath)
	if name == "." || name == "/" {
		name = "my-project"
	}

	// Sanitize name for npm
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	// Create default package.json
	pkg := NewPackageJSON(name)

	// Try to detect main file
	if _, err := os.Stat(filepath.Join(projectPath, "index.ts")); err == nil {
		pkg.Main = "index.ts"
	} else if _, err := os.Stat(filepath.Join(projectPath, "index.js")); err == nil {
		pkg.Main = "index.js"
	} else if _, err := os.Stat(filepath.Join(projectPath, "main.ts")); err == nil {
		pkg.Main = "main.ts"
	} else if _, err := os.Stat(filepath.Join(projectPath, "main.js")); err == nil {
		pkg.Main = "main.js"
	}

	// Save package.json
	if err := SavePackageJSON(pkgPath, pkg); err != nil {
		return err
	}

	fmt.Printf("Created package.json\n")
	fmt.Printf("\n")
	fmt.Printf("  name: %s\n", pkg.Name)
	fmt.Printf("  version: %s\n", pkg.Version)
	fmt.Printf("  main: %s\n", pkg.Main)
	fmt.Printf("  license: %s\n", pkg.License)
	fmt.Printf("\n")

	return nil
}

// InitProjectJSONInteractive creates package.json with user prompts
func InitProjectJSONInteractive(projectPath string) error {
	pkgPath := filepath.Join(projectPath, "package.json")

	// Check if already exists
	if _, err := os.Stat(pkgPath); err == nil {
		fmt.Print("package.json already exists. Overwrite? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			return nil
		}
	}

	reader := bufio.NewReader(os.Stdin)

	// Get project name
	defaultName := filepath.Base(projectPath)
	if defaultName == "." || defaultName == "/" {
		defaultName = "my-project"
	}
	defaultName = strings.ToLower(strings.ReplaceAll(defaultName, " ", "-"))

	fmt.Printf("package name: (%s) ", defaultName)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = defaultName
	}

	// Get version
	fmt.Print("version: (1.0.0) ")
	version, _ := reader.ReadString('\n')
	version = strings.TrimSpace(version)
	if version == "" {
		version = "1.0.0"
	}

	// Get description
	fmt.Print("description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// Get entry point
	defaultMain := "index.js"
	if _, err := os.Stat(filepath.Join(projectPath, "index.ts")); err == nil {
		defaultMain = "index.ts"
	}
	fmt.Printf("entry point: (%s) ", defaultMain)
	main, _ := reader.ReadString('\n')
	main = strings.TrimSpace(main)
	if main == "" {
		main = defaultMain
	}

	// Get keywords
	fmt.Print("keywords: ")
	keywordsStr, _ := reader.ReadString('\n')
	keywordsStr = strings.TrimSpace(keywordsStr)
	var keywords []string
	if keywordsStr != "" {
		for _, k := range strings.Split(keywordsStr, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keywords = append(keywords, k)
			}
		}
	}

	// Get author
	fmt.Print("author: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSpace(author)

	// Get license
	fmt.Print("license: (ISC) ")
	license, _ := reader.ReadString('\n')
	license = strings.TrimSpace(license)
	if license == "" {
		license = "ISC"
	}

	// Create package
	pkg := &PackageJSON{
		Name:        name,
		Version:     version,
		Description: description,
		Main:        main,
		Scripts: map[string]string{
			"test":  "echo \"Error: no test specified\" && exit 1",
			"start": "nu " + main,
		},
		Keywords:     keywords,
		Author:       author,
		License:      license,
		Dependencies: make(map[string]string),
	}

	// Show generated package.json
	fmt.Println("\nAbout to write to", pkgPath)
	fmt.Println("{")
	fmt.Printf("  \"name\": \"%s\",\n", pkg.Name)
	fmt.Printf("  \"version\": \"%s\",\n", pkg.Version)
	if pkg.Description != "" {
		fmt.Printf("  \"description\": \"%s\",\n", pkg.Description)
	}
	fmt.Printf("  \"main\": \"%s\",\n", pkg.Main)
	fmt.Printf("  \"scripts\": {\n")
	fmt.Printf("    \"test\": \"%s\",\n", pkg.Scripts["test"])
	fmt.Printf("    \"start\": \"%s\"\n", pkg.Scripts["start"])
	fmt.Printf("  },\n")
	if len(pkg.Keywords) > 0 {
		fmt.Printf("  \"keywords\": [\"%s\"],\n", strings.Join(pkg.Keywords, "\", \""))
	}
	if author != "" {
		fmt.Printf("  \"author\": \"%s\",\n", author)
	}
	fmt.Printf("  \"license\": \"%s\"\n", pkg.License)
	fmt.Println("}")

	fmt.Print("\nIs this OK? (yes) ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "" && confirm != "y" && confirm != "yes" {
		fmt.Println("Aborted.")
		return nil
	}

	// Save
	if err := SavePackageJSON(pkgPath, pkg); err != nil {
		return err
	}

	fmt.Println("Created package.json")
	return nil
}

// HandleInstallCommand handles the install command
func HandleInstallCommand(args []string, flags map[string]bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	installer := NewInstaller(cwd)

	// Check for flags
	isDev := flags["save-dev"] || flags["D"]
	isGlobal := flags["global"] || flags["g"]

	if isGlobal {
		return fmt.Errorf("global install not yet supported")
	}

	if len(args) == 0 {
		// Install all dependencies from package.json
		_, err := installer.InstallAll()
		return err
	}

	// Install specific package(s)
	for _, pkg := range args {
		// Parse package@version
		name := pkg
		version := ""
		if idx := strings.LastIndex(pkg, "@"); idx > 0 {
			name = pkg[:idx]
			version = pkg[idx+1:]
		}

		if err := installer.InstallPackage(name, version, isDev); err != nil {
			return err
		}
	}

	return nil
}

// HandleUninstallCommand handles the uninstall command
func HandleUninstallCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify package(s) to uninstall")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	installer := NewInstaller(cwd)

	for _, name := range args {
		if err := installer.UninstallPackage(name); err != nil {
			fmt.Printf("  ⚠ %s\n", err)
		}
	}

	return nil
}

// HandleListCommand handles the list command
func HandleListCommand() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	installer := NewInstaller(cwd)
	packages, err := installer.ListInstalled()
	if err != nil {
		return err
	}

	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return nil
	}

	// Load package.json for tree display
	pkgPath := filepath.Join(cwd, "package.json")
	pkg, _ := LoadPackageJSON(pkgPath)

	if pkg != nil {
		fmt.Printf("%s@%s\n", pkg.Name, pkg.Version)
	} else {
		fmt.Println("(no package.json)")
	}

	for _, p := range packages {
		prefix := "├──"
		if p.Dev {
			fmt.Printf("%s %s@%s (dev)\n", prefix, p.Name, p.Version)
		} else {
			fmt.Printf("%s %s@%s\n", prefix, p.Name, p.Version)
		}
	}

	return nil
}

// HandleRunCommand handles the run command
func HandleRunCommand(scriptName string, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pkgPath := filepath.Join(cwd, "package.json")
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		return fmt.Errorf("no package.json found")
	}

	// If no script name, list available scripts
	if scriptName == "" {
		if len(pkg.Scripts) == 0 {
			fmt.Println("No scripts defined in package.json")
			return nil
		}

		fmt.Println("Available scripts:")
		for name, cmd := range pkg.Scripts {
			fmt.Printf("  %s: %s\n", name, cmd)
		}
		return nil
	}

	// Get script command
	script, ok := pkg.Scripts[scriptName]
	if !ok {
		return fmt.Errorf("script '%s' not found in package.json", scriptName)
	}

	fmt.Printf("> %s\n", script)
	fmt.Println()

	// Execute script
	return executeScript(script, args, cwd)
}

// executeScript executes a script command
func executeScript(script string, args []string, cwd string) error {
	// Add node_modules/.bin to PATH
	binPath := filepath.Join(cwd, "node_modules", ".bin")
	path := os.Getenv("PATH")
	os.Setenv("PATH", binPath+string(os.PathListSeparator)+path)

	// Parse command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		fullCmd := script
		if len(args) > 0 {
			fullCmd += " " + strings.Join(args, " ")
		}
		cmd = exec.Command("cmd", "/C", fullCmd)
	} else {
		fullCmd := script
		if len(args) > 0 {
			fullCmd += " " + strings.Join(args, " ")
		}
		cmd = exec.Command("sh", "-c", fullCmd)
	}

	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// HandleOutdatedCommand shows outdated packages
func HandleOutdatedCommand() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pkgPath := filepath.Join(cwd, "package.json")
	pkg, err := LoadPackageJSON(pkgPath)
	if err != nil {
		return fmt.Errorf("no package.json found")
	}

	lockPath := filepath.Join(cwd, "package-lock.json")
	lock, err := LoadPackageLockJSON(lockPath)
	if err != nil {
		return fmt.Errorf("no package-lock.json found")
	}

	registry := NewNPMRegistry()

	fmt.Println("Checking for outdated packages...")
	fmt.Println()
	fmt.Printf("%-30s %-15s %-15s %-15s\n", "Package", "Current", "Wanted", "Latest")

	hasOutdated := false

	// Check all dependencies
	allDeps := pkg.GetAllDependencies()
	for name, versionRange := range allDeps {
		// Get installed version
		current, _ := lock.GetPackageVersion(name)
		if current == "" {
			current = "(not installed)"
		}

		// Get wanted version (satisfies range)
		wanted, err := registry.ResolveVersion(name, versionRange)
		if err != nil {
			continue
		}

		// Get latest version
		latest, err := registry.GetLatestVersion(name)
		if err != nil {
			continue
		}

		// Check if outdated
		if current != wanted || wanted != latest {
			hasOutdated = true
			fmt.Printf("%-30s %-15s %-15s %-15s\n", name, current, wanted, latest)
		}
	}

	if !hasOutdated {
		fmt.Println("All packages are up to date!")
	}

	return nil
}

// HandleUpdateCommand updates packages
func HandleUpdateCommand(packages []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	installer := NewInstaller(cwd)
	return installer.Update(packages)
}

// HandlePruneCommand removes extraneous packages
func HandlePruneCommand() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	installer := NewInstaller(cwd)
	removed, err := installer.Prune()
	if err != nil {
		return err
	}

	if len(removed) == 0 {
		fmt.Println("No extraneous packages to remove")
	} else {
		for _, name := range removed {
			fmt.Printf("Removed %s\n", name)
		}
	}

	return nil
}
