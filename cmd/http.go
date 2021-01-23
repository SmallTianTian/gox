package cmd

import (
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/http"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Add http server.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prepare()
		if config.DefaultConfig.HTTP.Port <= 0 {
			logger.Error("Http port must > 0.")
			return
		}

		if utils.IsExist(filepath.Join(config.DefaultConfig.Project.Path, "ui/http")) {
			logger.Error("Couldn't init http again.")
			return
		}

		logger.Debug("Will begin create http.")
		http.NewHTTP()
		utils.GoModRebuild(config.DefaultConfig.Project.Path)
	},
}

func init() {
	httpCmd.PersistentFlags().IntVar(&config.DefaultConfig.HTTP.Port, "port", 8080, "Set http server port.")
}
