package update

import (
	"os"
	"regexp"
	"strings"
)

const RELEASES_URL = "https://github.com/homeland-social/os/releases"
const RELEASE_PATH = "/root/etc/release"

func getReleasePath() string {
	path, exists := os.LookupEnv("RELEASE_PATH")

	if !exists {
		path = RELEASE_PATH
	}

	return path
}

func getArchVersion() (string, string, error) {
	var arch, version string

	path := getReleasePath()
	body, err := os.ReadFile(path)
	if err != nil {
		return arch, version, err
	}

	var matches [][]byte
	versionPattern := regexp.MustCompile(`VERSION=(.*)`)
	archPattern := regexp.MustCompile(`ARCH=(.*)`)
	lines := strings.Split(string(body[:]), "\n")

	for _, line := range lines {
		matches = versionPattern.FindSubmatch([]byte(line))
		if len(matches) > 1 {
			version = string(matches[1][:])
		}
		matches = archPattern.FindSubmatch([]byte(line))
		if len(matches) > 1 {
			arch = string(matches[1][:])
		}
	}

	return arch, version, nil
}

func CheckUpdates() {

}
