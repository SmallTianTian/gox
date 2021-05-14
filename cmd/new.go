package cmd

import (
	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/grpc"
	"github.com/SmallTianTian/fresh-go/internal/project"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "创建一个新的项目。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prepare()

		newProject(args[0])

		finally()
	},
}

func init() {
	newCmd.PersistentFlags().StringVarP(&config.DefaultConfig.Project.Path, "project_path", "p",
		config.DefaultConfig.Project.Path, "设置该项目创建的路径，归功于 go.mod，我们可以在任何地方创建项目。")
	newCmd.PersistentFlags().StringVarP(&config.DefaultConfig.Project.Org, "organization", "o",
		config.DefaultConfig.Project.Org, "项目的组织名，将在初始化 go.mod 的时候将其写入，并在导入本项目代码中体现。")
	newCmd.PersistentFlags().BoolVar(&config.DefaultConfig.Project.Vendor, "vendor", false, "是否使用 vendor，不推荐。")

	newCmd.PersistentFlags().BoolVar(
		&config.DefaultConfig.FrameEnable.HTTP, "http", true,
		"项目是否提供 HTTP 服务，默认开启。如需关闭，请使用 -http=false")
	newCmd.PersistentFlags().BoolVar(
		&config.DefaultConfig.FrameEnable.GRPC, "grpc", true,
		"项目是否提供 GRPC 服务，默认开启。如需关闭，请使用 -grpc=false")
	newCmd.PersistentFlags().BoolVar(
		&config.DefaultConfig.FrameEnable.Proxy, "proxy", true,
		"项目是否提供 GRPC Proxy 服务，默认开启。如需关闭，请使用 -proxy=false，注意，grpc 关闭后，proxy 也无法开启。")

	newCmd.PersistentFlags().StringVar(&config.DefaultConfig.Adapter.Logger, "log", "zap",
		"目前支持的日志框架：zap、logrus，可以使用小写的框架名来指定使用的日志框架。")
}

func newProject(name string) {
	// 设置项目名称
	config.DefaultConfig.Project.Name = name

	project.NewProject()

	if config.DefaultConfig.FrameEnable.GRPC {
		grpc.NewDemo()

		if config.DefaultConfig.FrameEnable.Proxy {
			grpc.NewDemoProxy()
		}
	}
}
