package cmd

import (
	"github.com/SmallTianTian/fresh-go/internal/doctor"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the running environment.",
	Run: func(cmd *cobra.Command, args []string) {
		prepare()

		doctor.CheckingEnv()
		logger.Debug("Will begin check running env.")
	},
}
