package cmd

import (
	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var (
	commands = []*cobra.Command{versionCmd, newCmd, doctorCmd, initCmd, apiCmd}
)

var rootCmd = &cobra.Command{
	Use:   "fresh-go",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		logger.InitLogger(config.DefaultConfig.Logger.Level)
	})

	initDebug()
	initRegister()
}

// 注册命令
func initRegister() {
	for _, command := range commands {
		rootCmd.AddCommand(command)
	}
}

var debug bool

// 注册 debug 标志
func initDebug() {
	for _, command := range commands {
		command.PersistentFlags().BoolVar(&debug, "debug", false, "Open debug.")
	}
}

// 最初的工作准备
func prepare() {
	if debug {
		logger.InitLogger("debug")
	}
}

// 最后的整理工作
func finally() {
	dir := config.DefaultConfig.Project.Path
	utils.GoModRebuild(dir)
	utils.GoFmtCode(dir)
}
