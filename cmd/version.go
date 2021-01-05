package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Version     string
	BuildTime   string
	GoVersion   string
	GitRevision string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version.",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func showVersion() {
	println("Version:", Version)
	println("Build time:", BuildTime)
	println("Go version:", GoVersion)
	println("Git commit hash:", GitRevision)
}
