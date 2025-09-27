package version

import "fmt"

const (
	// Version is the current version of FitBoisBot
	Version = "2.1.0"
)

// Variables that can be set at build time via ldflags
var (
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// GetVersion returns the current version
func GetVersion() string {
	return Version
}

// GetVersionInfo returns formatted version string for display
func GetVersionInfo() string {
	return fmt.Sprintf("FitBoisBot %s", Version)
}

// GetDetailedVersionInfo returns version with build details
func GetDetailedVersionInfo() string {
	info := fmt.Sprintf("FitBoisBot %s", Version)
	if GitCommit != "unknown" {
		if len(GitCommit) > 7 {
			info += fmt.Sprintf(" (%s)", GitCommit[:7])
		} else {
			info += fmt.Sprintf(" (%s)", GitCommit)
		}
	}
	if BuildTime != "unknown" {
		info += fmt.Sprintf("\nBuilt: %s", BuildTime)
	}
	return info
}