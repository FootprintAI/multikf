package version

import (
	"fmt"

	goversion "github.com/hashicorp/go-version"
)

var (
	BuildTime   = ""
	GitCommitId = ""
	version, _  = goversion.NewVersion("2.3.1")
)

func GetVersion() *goversion.Version {
	return version
}

func GetBuildTime() string {
	return BuildTime
}

func GetCommitHash() string {
	return GitCommitId
}

func Print() {
	fmt.Printf("ver:%s, build time:%s, hashid:%s\n", version.String(), BuildTime, GitCommitId)
}
