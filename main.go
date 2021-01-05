package main

import (
	"github.com/SmallTianTian/fresh-go/cmd"
)

var (
	Version     string
	BuildTime   string
	GoVersion   string
	GitRevision string
)

func main() {
	cmd.Version = Version
	cmd.BuildTime = BuildTime
	cmd.GoVersion = GoVersion
	cmd.GitRevision = GitRevision

	cmd.Execute()
}
