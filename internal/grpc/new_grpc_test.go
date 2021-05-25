package grpc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/project"
	"github.com/SmallTianTian/fresh-go/test"
	"github.com/SmallTianTian/fresh-go/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewDemo(t *testing.T) {
	Convey("", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir)

		// 准备工作
		config.DefaultConfig.Project.Path = dir
		config.DefaultConfig.Project.Name = "fresh"
		config.DefaultConfig.Project.Remote = "github.com"
		config.DefaultConfig.Project.Vendor = false
		path := filepath.Join(config.DefaultConfig.Project.Path, config.DefaultConfig.Project.Name)
		project.NewProject()

		Convey("能正常新建 demo 示例", func() {
			demoPath := filepath.Join(path, "api", "helloworld", "v1")
			demoImplPath := filepath.Join(path, "internal", "ui", "grpc", "helloworld_v1.go")
			grpcWirePath := filepath.Join(path, "internal", "ui", "grpc", "wire.go")
			serverWirePath := filepath.Join(path, "internal", "server", "wire.go")
			grpcServerPath := filepath.Join(path, "internal", "server", "grpc.go")
			bufPath := filepath.Join(path, "buf.yaml")
			bufGenPath := filepath.Join(path, "buf.gen.yaml")
			certSh := filepath.Join(path, "scripts", "create-cert.sh")

			So(utils.IsExist(demoPath), ShouldBeFalse)
			So(utils.IsExist(demoImplPath), ShouldBeFalse)
			So(utils.IsExist(grpcWirePath), ShouldBeFalse)
			So(utils.IsExist(serverWirePath), ShouldBeFalse)
			So(utils.IsExist(grpcServerPath), ShouldBeFalse)
			So(utils.IsExist(bufPath), ShouldBeFalse)
			So(utils.IsExist(bufGenPath), ShouldBeFalse)
			So(utils.IsExist(certSh), ShouldBeFalse)

			NewDemo()

			So(utils.IsExist(demoPath), ShouldBeTrue)
			So(utils.IsExist(demoImplPath), ShouldBeTrue)
			So(utils.IsExist(grpcWirePath), ShouldBeTrue)
			So(utils.IsExist(serverWirePath), ShouldBeTrue)
			So(utils.IsExist(grpcServerPath), ShouldBeTrue)
			So(utils.IsExist(bufPath), ShouldBeTrue)
			So(utils.IsExist(bufGenPath), ShouldBeTrue)
			So(utils.IsExist(certSh), ShouldBeTrue)
			utils.Exec(path, "go", "mod", "tidy")
			So(utils.Exec(path, "go", "vet", "./..."), ShouldEqual, nil)
		})
	})
}
