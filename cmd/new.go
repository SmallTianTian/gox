package cmd

import (
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/grpc"
	"github.com/SmallTianTian/fresh-go/internal/http"
	"github.com/SmallTianTian/fresh-go/internal/project"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a fresh project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prepare()

		// 设置项目名称
		config.DefaultConfig.Project.Name = args[0]

		logger.Debug("Will begin create new project.")
		project.NewProject()

		// 初始化其他模块
		logger.Debug("Will begin create grpc.")
		grpc.NewGrpc()
		logger.Debug("Will begin create http.")
		http.NewHTTP()

		// 执行 go mod
		utils.FirstMod(config.DefaultConfig.Project.Path,
			filepath.Join(config.DefaultConfig.Project.Org, config.DefaultConfig.Project.Name),
			config.DefaultConfig.Project.Vendor)
	},
}

func init() {
	newCmd.PersistentFlags().StringVarP(&config.DefaultConfig.Project.Path, "project_path", "p",
		config.DefaultConfig.Project.EnvPath, "The place where the project was created.")
	newCmd.PersistentFlags().StringVarP(&config.DefaultConfig.Project.Org, "organization", "o",
		config.DefaultConfig.Project.EnvOrg, "Your project organization.")
	newCmd.PersistentFlags().BoolVar(&config.DefaultConfig.Project.Vendor, "vendor", false, "Use the vendor directory.")
	newCmd.PersistentFlags().IntVar(&config.DefaultConfig.HTTP.Port, "http-port", 0, "Add http server with port.")
	newCmd.PersistentFlags().IntVar(&config.DefaultConfig.GRPC.Port, "grpc-port", 0, "Add grpc server with port.")
	newCmd.PersistentFlags().IntVar(&config.DefaultConfig.GRPC.Proxy, "proxy-port", 0, "Add grpc proxy server with port.")
}
