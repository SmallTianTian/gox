package cmd

import (
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/internal/http"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Create a fresh project.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		InitHttp()
	},
}

func init() {
	httpCmd.PersistentFlags().IntVar(&HttpPort, "port", 8080, "Set http server port.")
}

func InitHttp() {
	if HttpPort != 0 && !utils.IsExist(filepath.Join(ProjectPath, "ui/http")) {
		http.NewHttp(ProjectPath, Organization, HttpPort)
		if !IsNewProject {
			utils.GoModRebuild(ProjectPath)
		}
	}
}
