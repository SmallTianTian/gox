package cmd

import (
	"os"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/internal/new"
	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a fresh project.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		ProjectName = args[0]
		ProjectPath = filepath.Join(ProjectPath, ProjectName)
		IsNewProject = true
		NewProject()
	},
}

func init() {
	pwd, _ := os.Getwd()
	newCmd.PersistentFlags().StringVarP(&ProjectPath, "project_path", "p", pwd, "The place where the project was created.")
	newCmd.PersistentFlags().StringVarP(&Organization, "organization", "o", "github.com", "Your project organization.")
	newCmd.PersistentFlags().BoolVar(&Vendor, "vendor", false, "Use the vendor directory.")
	newCmd.PersistentFlags().IntVar(&HttpPort, "http-port", 0, "Add http server with port.")
	newCmd.PersistentFlags().IntVar(&GrpcPort, "grpc-port", 0, "Add grpc server with port.")
	newCmd.PersistentFlags().IntVar(&GrpcProxyPort, "proxy-port", 0, "Add grpc server with port.")
}

func NewProject() {
	err := new.NewProject(ProjectPath, Organization, Vendor)
	utils.MustNotError(err)

	// init another server
	InitHttp()
	InitGrpc()

	utils.FirstMod(ProjectPath, filepath.Join(Organization, ProjectName), Vendor)
}

func notNewSetProjectInfo() {
	ProjectPath, _ = os.Getwd()
	oAp := utils.GetOrganizationAndProjectName(ProjectPath)
	utils.MustNotBlank(oAp)

	ProjectName = filepath.Base(oAp)
	Organization = filepath.Dir(oAp)
}
