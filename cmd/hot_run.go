package cmd

import (
	"os"
	"os/signal"

	"github.com/SmallTianTian/fresh-go/internal/hot_run"
	"github.com/spf13/cobra"
)

var hotRunCmd = &cobra.Command{
	Use:   "hotrun",
	Short: "Hot compile and run.",
	Run: func(cmd *cobra.Command, args []string) {
		defer hot_run.ClearRun()
		go hot_run.HotRun(args)
		flag := make(chan os.Signal)
		signal.Notify(flag, os.Interrupt, os.Kill)
		<-flag
	},
}
