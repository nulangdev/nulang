package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Semver represents a parsed semantic version
type Semver struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
	Raw        string
}

// semver regex pattern
var semverRegex = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// ParseSemver parses a semantic version string
func ParseSemver(version string) (*Semver, error) {
	version = strings.TrimSpace(version)
	
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid semver: %s", version)
	}

	major, _ := strconv.Atoi(matches[1])
	minor := 0
	if matches[2] != "" {
		minor, _ = strconv.Atoi(matches[2])
	}
	patch := 0
	if matches[3] != "" {
		patch, _ = strconv.Atoi(matches[3])
	}

	return &Semver{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: matches[4],
		Build:      matches[5],
		Raw:        version,
	}, nil
}

// String returns the string representation of the version
func (s *Semver) String() string {
	version := fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
	if s.Prerelease != "" {
		version += "-" + s.Prerelease
	}
	if s.Build != "" {
		version += "+" + s.Build
	}
	return version
}

// Compare compares two versions
// Returns: -1 if s < other, 0 if s == other, 1 if s > other
func (s *Semver) Compare(other *Semver) int {
	if s.Major != other.Major {
		if s.Major < other.Major {
			return -1
		}
		return 1
	}
	if s.Minor != other.Minor {
		if s.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if s.Patch != other.Patch {
		if s.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Prerelease comparison
	if s.Prerelease != "" && other.Prerelease == "" {
		return -1 // Prerelease is always less than release
	}
	if s.Prerelease == "" && other.Prerelease != "" {
		return 1
	}
	if s.Prerelease != other.Prerelease {
		return comparePrerelease(s.Prerelease, other.Prerelease)
	}

	return 0
}

func comparePrerelease(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aNum, aIsNum := strconv.Atoi(aParts[i])
		bNum, bIsNum := strconv.Atoi(bParts[i])

		if aIsNum == nil && bIsNum == nil {
			if aNum < bNum {
				return -1
			}
			if aNum > bNum {
				return 1
			}
		} else if aIsNum == nil {
			return -1 // Numeric < string
		} else if bIsNum == nil {
			return 1
		} else {
			if aParts[i] < bParts[i] {
				return -1
			}
			if aParts[i] > bParts[i] {
				return 1
			}
		}
	}

	if len(aParts) < len(bParts) {
		return -1
	}
	if len(aParts) > len(bParts) {
		return 1
	}

	return 0
}

// LessThan returns true if s < other
func (s *Semver) LessThan(other *Semver) bool {
	return s.Compare(other) < 0
}

// LessThanOrEqual returns true if s <= other
func (s *Semver) LessThanOrEqual(other *Semver) bool {
	return s.Compare(other) <= 0
}

// GreaterThan returns true if s > other
func (s *Semver) GreaterThan(other *Semver) bool {
	return s.Compare(other) > 0
}

// GreaterThanOrEqual returns true if s >= other
func (s *Semver) GreaterThanOrEqual(other *Semver) bool {
	return s.Compare(other) >= 0
}

// Equal returns true if s == other
func (s *Semver) Equal(other *Semver) bool {
	return s.Compare(other) == 0
}

// SemverRange represents a version range constraint
type SemverRange struct {
	Raw         string
	Comparators []SemverComparator
	IsOR        bool   // If true, comparators are ORed, otherwise ANDed
	SubRanges   []*SemverRange // For complex OR ranges
}

// SemverComparator represents a single version comparator
type SemverComparator struct {
	Operator string  // "", "=", ">", ">=", "<", "<=", "^", "~"
	Version  *Semver
}

// ParseSemverRange parses a semver range string
// Supports: exact, ^, ~, >, >=, <, <=, ||, -, x, *
func ParseSemverRange(rangeStr string) (*SemverRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	
	if rangeStr == "" || rangeStr == "*" || rangeStr == "latest" {
		return &SemverRange{
			Raw:         rangeStr,
			Comparators: []SemverComparator{{Operator: ">=", Version: &Semver{Major: 0, Minor: 0, Patch: 0}}},
		}, nil
	}

	// Handle OR ranges (||)
	if strings.Contains(rangeStr, "||") {
		parts := strings.Split(rangeStr, "||")
		subRanges := make([]*SemverRange, 0, len(parts))
		for _, part := range parts {
			subRange, err := ParseSemverRange(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			subRanges = append(subRanges, subRange)
		}
		return &SemverRange{
			Raw:       rangeStr,
			IsOR:      true,
			SubRanges: subRanges,
		}, nil
	}

	// Handle hyphen ranges (1.0.0 - 2.0.0)
	if strings.Contains(rangeStr, " - ") {
		parts := strings.Split(rangeStr, " - ")
		if len(parts) == 2 {
			lower, err := ParseSemver(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, err
			}
			upper, err := ParseSemver(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, err
			}
			return &SemverRange{
				Raw: rangeStr,
				Comparators: []SemverComparator{
					{Operator: ">=", Version: lower},
					{Operator: "<=", Version: upper},
				},
			}, nil
		}
	}

	// Handle space-separated AND ranges
	parts := strings.Fields(rangeStr)
	comparators := make([]SemverComparator, 0)

	for _, part := range parts {
		comp, err := parseComparator(part)
		if err != nil {
			return nil, err
		}
		comparators = append(comparators, comp...)
	}

	return &SemverRange{
		Raw:         rangeStr,
		Comparators: comparators,
	}, nil
}

func parseComparator(part string) ([]SemverComparator, error) {
	part = strings.TrimSpace(part)
	if part == "" {
		return nil, nil
	}

	// Handle operators
	var operator string
	var versionStr string

	if strings.HasPrefix(part, ">=") {
		operator = ">="
		versionStr = part[2:]
	} else if strings.HasPrefix(part, "<=") {
		operator = "<="
		versionStr = part[2:]
	} else if strings.HasPrefix(part, ">") {
		operator = ">"
		versionStr = part[1:]
	} else if strings.HasPrefix(part, "<") {
		operator = "<"
		versionStr = part[1:]
	} else if strings.HasPrefix(part, "^") {
		operator = "^"
		versionStr = part[1:]
	} else if strings.HasPrefix(part, "~") {
		operator = "~"
		versionStr = part[1:]
	} else if strings.HasPrefix(part, "=") {
		operator = "="
		versionStr = part[1:]
	} else {
		operator = "="
		versionStr = part
	}

	versionStr = strings.TrimSpace(versionStr)

	// Handle x-ranges (1.x, 1.2.x)
	if strings.Contains(versionStr, "x") || strings.Contains(versionStr, "X") || strings.Contains(versionStr, "*") {
		return parseXRange(versionStr)
	}

	version, err := ParseSemver(versionStr)
	if err != nil {
		return nil, err
	}

	// Expand ^ and ~ operators
	switch operator {
	case "^":
		return expandCaret(version), nil
	case "~":
		return expandTilde(version), nil
	default:
		return []SemverComparator{{Operator: operator, Version: version}}, nil
	}
}

// expandCaret expands ^1.2.3 to >=1.2.3 <2.0.0
func expandCaret(v *Semver) []SemverComparator {
	upper := &Semver{}
	
	if v.Major != 0 {
		// ^1.2.3 := >=1.2.3 <2.0.0
		upper.Major = v.Major + 1
		upper.Minor = 0
		upper.Patch = 0
	} else if v.Minor != 0 {
		// ^0.2.3 := >=0.2.3 <0.3.0
		upper.Major = 0
		upper.Minor = v.Minor + 1
		upper.Patch = 0
	} else {
		// ^0.0.3 := >=0.0.3 <0.0.4
		upper.Major = 0
		upper.Minor = 0
		upper.Patch = v.Patch + 1
	}

	return []SemverComparator{
		{Operator: ">=", Version: v},
		{Operator: "<", Version: upper},
	}
}

// expandTilde expands ~1.2.3 to >=1.2.3 <1.3.0
func expandTilde(v *Semver) []SemverComparator {
	upper := &Semver{
		Major: v.Major,
		Minor: v.Minor + 1,
		Patch: 0,
	}

	return []SemverComparator{
		{Operator: ">=", Version: v},
		{Operator: "<", Version: upper},
	}
}

// parseXRange parses x-ranges like 1.x, 1.2.x
func parseXRange(versionStr string) ([]SemverComparator, error) {
	// Replace x, X, * with 0 for parsing
	parts := strings.Split(versionStr, ".")
	
	if len(parts) == 1 || parts[1] == "x" || parts[1] == "X" || parts[1] == "*" {
		// 1.x or 1.* -> >=1.0.0 <2.0.0
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		return []SemverComparator{
			{Operator: ">=", Version: &Semver{Major: major, Minor: 0, Patch: 0}},
			{Operator: "<", Version: &Semver{Major: major + 1, Minor: 0, Patch: 0}},
		}, nil
	}

	if len(parts) >= 2 && (len(parts) < 3 || parts[2] == "x" || parts[2] == "X" || parts[2] == "*") {
		// 1.2.x -> >=1.2.0 <1.3.0
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		return []SemverComparator{
			{Operator: ">=", Version: &Semver{Major: major, Minor: minor, Patch: 0}},
			{Operator: "<", Version: &Semver{Major: major, Minor: minor + 1, Patch: 0}},
		}, nil
	}

	return nil, fmt.Errorf("invalid x-range: %s", versionStr)
}

// Match tests if a version matches the range
func (r *SemverRange) Match(version string) bool {
	v, err := ParseSemver(version)
	if err != nil {
		return false
	}
	return r.MatchVersion(v)
}

// MatchVersion tests if a Semver matches the range
func (r *SemverRange) MatchVersion(v *Semver) bool {
	// Handle OR sub-ranges
	if r.IsOR && len(r.SubRanges) > 0 {
		for _, subRange := range r.SubRanges {
			if subRange.MatchVersion(v) {
				return true
			}
		}
		return false
	}

	// All comparators must match (AND)
	for _, comp := range r.Comparators {
		if !comp.Match(v) {
			return false
		}
	}

	return true
}

// Match tests if a version matches the comparator
func (c *SemverComparator) Match(v *Semver) bool {
	switch c.Operator {
	case "", "=":
		return v.Equal(c.Version)
	case ">":
		return v.GreaterThan(c.Version)
	case ">=":
		return v.GreaterThanOrEqual(c.Version)
	case "<":
		return v.LessThan(c.Version)
	case "<=":
		return v.LessThanOrEqual(c.Version)
	default:
		return false
	}
}

// MaxSatisfying returns the highest version in the list that satisfies the range
func (r *SemverRange) MaxSatisfying(versions []string) string {
	var max *Semver
	var maxStr string

	for _, vStr := range versions {
		v, err := ParseSemver(vStr)
		if err != nil {
			continue
		}
		if !r.MatchVersion(v) {
			continue
		}
		if max == nil || v.GreaterThan(max) {
			max = v
			maxStr = vStr
		}
	}

	return maxStr
}

// MinSatisfying returns the lowest version in the list that satisfies the range
func (r *SemverRange) MinSatisfying(versions []string) string {
	var min *Semver
	var minStr string

	for _, vStr := range versions {
		v, err := ParseSemver(vStr)
		if err != nil {
			continue
		}
		if !r.MatchVersion(v) {
			continue
		}
		if min == nil || v.LessThan(min) {
			min = v
			minStr = vStr
		}
	}

	return minStr
}

// ValidSemver checks if a string is a valid semver
func ValidSemver(version string) bool {
	_, err := ParseSemver(version)
	return err == nil
}

// CleanVersion removes leading 'v' and extra whitespace
func CleanVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")
	return version
}

// IncrementMajor returns a new version with major incremented
func (s *Semver) IncrementMajor() *Semver {
	return &Semver{
		Major: s.Major + 1,
		Minor: 0,
		Patch: 0,
	}
}

// IncrementMinor returns a new version with minor incremented
func (s *Semver) IncrementMinor() *Semver {
	return &Semver{
		Major: s.Major,
		Minor: s.Minor + 1,
		Patch: 0,
	}
}

// IncrementPatch returns a new version with patch incremented
func (s *Semver) IncrementPatch() *Semver {
	return &Semver{
		Major: s.Major,
		Minor: s.Minor,
		Patch: s.Patch + 1,
	}
}
