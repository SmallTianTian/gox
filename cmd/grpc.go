package cmd

import (
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/grpc"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Add grpc server.",
	Run: func(cmd *cobra.Command, args []string) {
		prepare()
		if config.DefaultConfig.GRPC.Port <= 0 {
			logger.Error("Grpc port must > 0.")
			return
		}

		if utils.IsExist(filepath.Join(config.DefaultConfig.Project.Path, "ui/grpc")) {
			logger.Error("Couldn't init grpc again.")
			return
		}

		logger.Debug("Will begin create grpc.")
		grpc.NewGrpc()
		utils.GoModRebuild(config.DefaultConfig.Project.Path)
	},
}

func init() {
	grpcCmd.PersistentFlags().IntVar(&config.DefaultConfig.GRPC.Port, "port", 50051, "Set grpc server port.")
	grpcCmd.PersistentFlags().IntVar(&config.DefaultConfig.GRPC.Proxy, "proxy", 0,
		"Set grpc proxy server port. If set, will create grpc proxy. Default not create.")
}
