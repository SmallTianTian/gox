package cmd

import (
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/internal/grpc"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Create a fresh project.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		InitGrpc()
	},
}

func init() {
	grpcCmd.PersistentFlags().IntVar(&GrpcPort, "port", 50051, "Set grpc server port.")
	grpcCmd.PersistentFlags().IntVar(&GrpcProxyPort, "proxy-port", 0, "Create grpc port server port.")
}

func InitGrpc() {
	if GrpcPort != 0 && !utils.IsExist(filepath.Join(ProjectPath, "ui/grpc")) {
		grpc.NewGrpc(ProjectPath, Organization, GrpcPort, GrpcProxyPort)
		if !IsNewProject {
			utils.GoModRebuild(ProjectPath)
		}
	}
}
