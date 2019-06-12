package version

var (
	// Version is the version of containerenv
	Version = "v0.1.0"

	// BuildMetadata is extra build time data
	BuildMetadata = "unreleased"
	// GitCommit is the git sha1
	GitCommit = ""
	// GitTreeState is the state of the git tree
	GitTreeState = ""
)

// GetVersion returns the semver string of the version
func GetVersion() string {
	if BuildMetadata == "" {
		return Version
	}
	return Version + "+" + BuildMetadata
}
