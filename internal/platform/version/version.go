package version

import (
	"strconv"
	"strings"
)

// Current application semantic version
const Current = "1.5.0"

// DefaultGitHubRepo represents the default repository for release queries
const DefaultGitHubRepo = "zfhassaan/school-management-system"

// IsNewerVersion returns true if remoteVer is strictly greater than localVer (semver).
func IsNewerVersion(localVer, remoteVer string) bool {
	cleanLocal := strings.TrimPrefix(strings.TrimSpace(localVer), "v")
	cleanRemote := strings.TrimPrefix(strings.TrimSpace(remoteVer), "v")

	localParts := strings.Split(cleanLocal, ".")
	remoteParts := strings.Split(cleanRemote, ".")

	maxLen := len(localParts)
	if len(remoteParts) > maxLen {
		maxLen = len(remoteParts)
	}

	for i := 0; i < maxLen; i++ {
		var lVal, rVal int
		if i < len(localParts) {
			lVal, _ = strconv.Atoi(localParts[i])
		}
		if i < len(remoteParts) {
			rVal, _ = strconv.Atoi(remoteParts[i])
		}

		if rVal > lVal {
			return true
		}
		if rVal < lVal {
			return false
		}
	}

	return false
}

