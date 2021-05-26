package cmd

import (
	"strings"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/grpc"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:     "api [service]",
	Short:   "创建一个新的 GRPC service。",
	Args:    cobra.ExactArgs(1),
	Example: "fresh-go api user",
	Run: func(cmd *cobra.Command, args []string) {
		prepare()

		newAPI(args[0])

		finally()
	},
}

func init() {
	apiCmd.PersistentFlags().BoolVar(
		&config.DefaultConfig.FrameEnable.Proxy, "proxy", true,
		"项目是否提供 GRPC Proxy 服务，默认开启。如需关闭，请使用 --proxy=false。")
}

func newAPI(serviceName string) {
	dir := config.DefaultConfig.Project.Path
	if !utils.CheckGoProject(dir) {
		logger.Fatalf("`%s` not golang project path.", dir)
	}

	serviceName = strings.TrimSpace(serviceName)
	if !utils.LegalVarName(serviceName) {
		logger.Fatalf("Not legal var name: `%s`.", serviceName)
	}

	grpc.New(serviceName)
	if config.DefaultConfig.FrameEnable.Proxy {
		grpc.NewProxy(serviceName, "")
	}
}
