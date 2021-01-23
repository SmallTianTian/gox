package http

import (
	"os"
	"testing"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/project"
	"github.com/SmallTianTian/fresh-go/test"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewHTTP(t *testing.T) {
	Convey("", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir)

		config.DefaultConfig.Project.Name = "fresh"
		config.DefaultConfig.Project.Org = "github.com"
		config.DefaultConfig.Project.Vendor = false
		config.DefaultConfig.Project.Path = dir
		project.NewProject()
		Convey("正常初始化 HTTP 服务", func() {
			
		})
		Convey("2", func() {})
		Convey("3", func() {})
	})
}
