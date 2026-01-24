package evaluator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PackageJSON represents the package.json structure (npm compatible)
type PackageJSON struct {
	Name             string                 `json:"name"`
	Version          string                 `json:"version"`
	Description      string                 `json:"description,omitempty"`
	Main             string                 `json:"main,omitempty"`
	Module           string                 `json:"module,omitempty"`
	Types            string                 `json:"types,omitempty"`
	Type             string                 `json:"type,omitempty"` // "module" or "commonjs"
	Scripts          map[string]string      `json:"scripts,omitempty"`
	Dependencies     map[string]string      `json:"dependencies,omitempty"`
	DevDependencies  map[string]string      `json:"devDependencies,omitempty"`
	PeerDependencies map[string]string      `json:"peerDependencies,omitempty"`
	OptionalDependencies map[string]string  `json:"optionalDependencies,omitempty"`
	Keywords         []string               `json:"keywords,omitempty"`
	Author           interface{}            `json:"author,omitempty"`     // string or object {name, email, url}
	Contributors     []interface{}          `json:"contributors,omitempty"`
	License          string                 `json:"license,omitempty"`
	Repository       interface{}            `json:"repository,omitempty"` // string or object {type, url}
	Bugs             interface{}            `json:"bugs,omitempty"`       // string or object {url, email}
	Homepage         string                 `json:"homepage,omitempty"`
	Engines          map[string]string      `json:"engines,omitempty"`
	Private          bool                   `json:"private,omitempty"`
	Exports          interface{}            `json:"exports,omitempty"`    // string, object, or conditional exports
	Bin              interface{}            `json:"bin,omitempty"`        // string or map[string]string
	Files            []string               `json:"files,omitempty"`
	Directories      map[string]string      `json:"directories,omitempty"`
	Man              interface{}            `json:"man,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
	PublishConfig    map[string]interface{} `json:"publishConfig,omitempty"`
	Workspaces       []string               `json:"workspaces,omitempty"`
	OS               []string               `json:"os,omitempty"`
	CPU              []string               `json:"cpu,omitempty"`
	Funding          interface{}            `json:"funding,omitempty"`
}

// Author represents the author field as an object
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// Repository represents the repository field as an object
type Repository struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Directory string `json:"directory,omitempty"`
}

// Bugs represents the bugs field as an object
type Bugs struct {
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// LoadPackageJSON loads and parses the package.json file
func LoadPackageJSON(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Set defaults
	if pkg.Main == "" {
		pkg.Main = "index.js"
	}
	if pkg.Version == "" {
		pkg.Version = "1.0.0"
	}

	return &pkg, nil
}

// SavePackageJSON saves the package configuration to package.json
func SavePackageJSON(path string, pkg *PackageJSON) error {
	// Pretty print with 2-space indentation (npm standard)
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	// Add trailing newline (npm standard)
	data = append(data, '\n')

	return os.WriteFile(path, data, 0644)
}

// FindPackageJSON finds the nearest package.json by traversing up the directory tree
func FindPackageJSON(startPath string) (string, error) {
	dir := startPath
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		dir = filepath.Dir(dir)
	}

	for {
		pkgPath := filepath.Join(dir, "package.json")
		if _, err := os.Stat(pkgPath); err == nil {
			return pkgPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("package.json not found")
}

// GetAuthorString returns the author as a string (handles both string and object formats)
func (p *PackageJSON) GetAuthorString() string {
	if p.Author == nil {
		return ""
	}

	switch v := p.Author.(type) {
	case string:
		return v
	case map[string]interface{}:
		name, _ := v["name"].(string)
		email, _ := v["email"].(string)
		url, _ := v["url"].(string)
		
		result := name
		if email != "" {
			result += " <" + email + ">"
		}
		if url != "" {
			result += " (" + url + ")"
		}
		return result
	default:
		return ""
	}
}

// SetAuthor sets the author field (accepts string or Author struct)
func (p *PackageJSON) SetAuthor(author interface{}) {
	p.Author = author
}

// GetRepositoryURL returns the repository URL
func (p *PackageJSON) GetRepositoryURL() string {
	if p.Repository == nil {
		return ""
	}

	switch v := p.Repository.(type) {
	case string:
		return v
	case map[string]interface{}:
		url, _ := v["url"].(string)
		return url
	default:
		return ""
	}
}

// GetBin returns the bin field as a map
func (p *PackageJSON) GetBin() map[string]string {
	if p.Bin == nil {
		return nil
	}

	switch v := p.Bin.(type) {
	case string:
		// Single binary uses package name
		return map[string]string{p.Name: v}
	case map[string]interface{}:
		result := make(map[string]string)
		for k, val := range v {
			if s, ok := val.(string); ok {
				result[k] = s
			}
		}
		return result
	default:
		return nil
	}
}

// GetExportsMain returns the main export path
func (p *PackageJSON) GetExportsMain() string {
	if p.Exports == nil {
		return p.Main
	}

	switch v := p.Exports.(type) {
	case string:
		return v
	case map[string]interface{}:
		// Try "." first (main entry point)
		if main, ok := v["."]; ok {
			return resolveExportCondition(main)
		}
		// Try "import" or "require" directly
		return resolveExportCondition(v)
	default:
		return p.Main
	}
}

// resolveExportCondition resolves conditional exports
func resolveExportCondition(export interface{}) string {
	switch v := export.(type) {
	case string:
		return v
	case map[string]interface{}:
		// Priority: import > require > default > node
		priorities := []string{"import", "require", "default", "node"}
		for _, key := range priorities {
			if val, ok := v[key]; ok {
				return resolveExportCondition(val)
			}
		}
		// Return first value if no priority match
		for _, val := range v {
			return resolveExportCondition(val)
		}
	}
	return ""
}

// GetExportPath resolves a subpath export (e.g., "lodash/get")
func (p *PackageJSON) GetExportPath(subpath string) (string, bool) {
	if p.Exports == nil {
		// No exports field, use direct file resolution
		return "", false
	}

	exports, ok := p.Exports.(map[string]interface{})
	if !ok {
		return "", false
	}

	// Normalize subpath
	if subpath == "" || subpath == "." {
		subpath = "."
	} else if subpath[0] != '.' {
		subpath = "./" + subpath
	}

	// Direct match
	if val, ok := exports[subpath]; ok {
		resolved := resolveExportCondition(val)
		return resolved, resolved != ""
	}

	// Pattern matching with wildcards (e.g., "./*": "./lib/*.js")
	for pattern, target := range exports {
		if matched, result := matchExportPattern(pattern, subpath, target); matched {
			return result, true
		}
	}

	return "", false
}

// matchExportPattern matches export patterns with wildcards
func matchExportPattern(pattern, subpath string, target interface{}) (bool, string) {
	// Simple wildcard matching for patterns like "./*"
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		if len(subpath) >= len(prefix) && subpath[:len(prefix)] == prefix {
			suffix := subpath[len(prefix):]
			targetStr := resolveExportCondition(target)
			if targetStr != "" && len(targetStr) > 0 && targetStr[len(targetStr)-1] == '*' {
				return true, targetStr[:len(targetStr)-1] + suffix
			}
		}
	}
	return false, ""
}

// HasScript checks if a script exists
func (p *PackageJSON) HasScript(name string) bool {
	if p.Scripts == nil {
		return false
	}
	_, ok := p.Scripts[name]
	return ok
}

// GetScript returns a script by name
func (p *PackageJSON) GetScript(name string) (string, bool) {
	if p.Scripts == nil {
		return "", false
	}
	script, ok := p.Scripts[name]
	return script, ok
}

// AddDependency adds a dependency
func (p *PackageJSON) AddDependency(name, version string) {
	if p.Dependencies == nil {
		p.Dependencies = make(map[string]string)
	}
	p.Dependencies[name] = version
}

// AddDevDependency adds a dev dependency
func (p *PackageJSON) AddDevDependency(name, version string) {
	if p.DevDependencies == nil {
		p.DevDependencies = make(map[string]string)
	}
	p.DevDependencies[name] = version
}

// RemoveDependency removes a dependency from all dependency fields
func (p *PackageJSON) RemoveDependency(name string) bool {
	removed := false
	
	if p.Dependencies != nil {
		if _, ok := p.Dependencies[name]; ok {
			delete(p.Dependencies, name)
			removed = true
		}
	}
	
	if p.DevDependencies != nil {
		if _, ok := p.DevDependencies[name]; ok {
			delete(p.DevDependencies, name)
			removed = true
		}
	}
	
	if p.PeerDependencies != nil {
		if _, ok := p.PeerDependencies[name]; ok {
			delete(p.PeerDependencies, name)
			removed = true
		}
	}
	
	if p.OptionalDependencies != nil {
		if _, ok := p.OptionalDependencies[name]; ok {
			delete(p.OptionalDependencies, name)
			removed = true
		}
	}
	
	return removed
}

// GetAllDependencies returns all dependencies (including dev, peer, optional)
func (p *PackageJSON) GetAllDependencies() map[string]string {
	all := make(map[string]string)
	
	for name, version := range p.Dependencies {
		all[name] = version
	}
	for name, version := range p.DevDependencies {
		all[name] = version
	}
	for name, version := range p.PeerDependencies {
		all[name] = version
	}
	for name, version := range p.OptionalDependencies {
		all[name] = version
	}
	
	return all
}

// NewPackageJSON creates a new PackageJSON with default values
func NewPackageJSON(name string) *PackageJSON {
	return &PackageJSON{
		Name:         name,
		Version:      "1.0.0",
		Description:  "",
		Main:         "index.js",
		Scripts:      map[string]string{
			"test": "echo \"Error: no test specified\" && exit 1",
		},
		Keywords:     []string{},
		Author:       "",
		License:      "ISC",
		Dependencies: make(map[string]string),
	}
}

// FindNodeModulesPath finds the node_modules directory by traversing up
func FindNodeModulesPath(startPath string) string {
	dir := startPath
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		dir = filepath.Dir(dir)
	}

	for {
		modulesPath := filepath.Join(dir, "node_modules")
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
