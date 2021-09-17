package main

import (
	"tianxu.xin/gox/cmd"
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
