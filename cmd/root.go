package cmd

import (
	"github.com/spf13/cobra"
	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/pkg/logger"
)

var (
	commands = []*cobra.Command{versionCmd, newCmd}
)

var rootCmd = &cobra.Command{
	Use:   "gox",
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

func setDebug(conf *config.Config) {
	if debug {
		conf.Logger.Level = "debug"
	}
}
