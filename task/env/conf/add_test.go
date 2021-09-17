package conf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"tianxu.xin/gox/internal/test"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task/env"
)

func TestAddTask_Check(t *testing.T) {
	Convey("测试检查", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		ctx := test.GetTestCtx(filepath.Join(dir, "tt"))
		task := NewAddTask(&AddConfig{[]string{"Logger"}, "level", "string", "debug"})

		Convey("目录不存在的情况下，检查未通过", func() {
			flag := task.PreCheck(ctx)
			So(flag, ShouldBeFalse)
		})

		Convey("非法参数名称，检查未通过", func() {
			task := NewAddTask(&AddConfig{[]string{"Logger"}, "-bad", "string", "value"})
			flag := task.PreCheck(ctx)
			So(flag, ShouldBeFalse)
		})

		Convey("目录已存在的情况下，不包含 internal/conf/config.go 文件，检查未通过", func() {
			ctx = env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
			// 移除 internal/conf/config.go 文件
			os.Remove(filepath.Join(util.MustExtractProjectPath(ctx), "internal/conf/config.go"))

			flag := task.PreCheck(ctx)
			So(flag, ShouldBeFalse)
		})

		Convey("目录已存在的情况下，不包含 configs/config.yaml.example 文件，检查未通过", func() {
			ctx = env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
			// 移除 configs/config.yaml.example 文件
			os.Remove(filepath.Join(util.MustExtractProjectPath(ctx), "configs/config.yaml.example"))

			flag := task.PreCheck(ctx)
			So(flag, ShouldBeFalse)
		})

		Convey("目录已存在的情况下，不包含 configs/config.yaml 文件，检查未通过", func() {
			ctx = env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
			// 移除 configs/config.yaml 文件
			os.Remove(filepath.Join(util.MustExtractProjectPath(ctx), "configs/config.yaml"))

			flag := task.PreCheck(ctx)
			So(flag, ShouldBeFalse)
		})

		Convey("正常初始化项目，能通过检查", func() {
			ctx = env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())

			flag := task.PreCheck(ctx)
			So(flag, ShouldBeTrue)
		})
	})
}

func TestAddTask_Exec(t *testing.T) {
	Convey("测试新建项目", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
		task := NewAddTask(&AddConfig{[]string{"Logger"}, "level", "string", "debug"})

		Convey("正常执行不报错，能设置内容，且能通过检查", func() {
			pp := util.MustExtractProjectPath(ctx)
			So(util.IsExist(pp), ShouldBeTrue)

			// 刚创建项目，两个 yaml 文件应该是空
			cc := string(util.ReadFile(filepath.Join(pp, "configs/config.yaml")))
			cec := string(util.ReadFile(filepath.Join(pp, "configs/config.yaml.example")))
			So(cc, ShouldBeEmpty)
			So(cec, ShouldBeEmpty)
			// go 文件中也没有 logger 的字样
			gc := string(util.ReadFile(filepath.Join(pp, "internal/conf/config.go")))
			So(strings.Contains(gc, "Logger"), ShouldBeFalse)

			err := task.Exec(ctx)
			So(err, ShouldBeNil)

			// 任务执行完成后，也能通过 vet 检查
			err = util.Exec(pp, "go", "vet", "./...")
			So(err, ShouldBeNil)

			except := `Logger:
  level: debug
`

			// 刚创建项目，两个 yaml 文件应该是空
			cc = string(util.ReadFile(filepath.Join(pp, "configs/config.yaml")))
			cec = string(util.ReadFile(filepath.Join(pp, "configs/config.yaml.example")))
			So(cc, ShouldEqual, except)
			So(cec, ShouldEqual, except)
			// go 文件中也没有 logger 的字样
			gc = string(util.ReadFile(filepath.Join(pp, "internal/conf/config.go")))
			So(strings.Contains(gc, "Logger"), ShouldBeTrue)
		})

		Convey("多次执行不报错，能设置内容，且能通过检查", func() {
			pp := util.MustExtractProjectPath(ctx)
			So(util.IsExist(pp), ShouldBeTrue)

			err := task.Exec(ctx)
			So(err, ShouldBeNil)

			// 再次执行同一个 struct
			task = NewAddTask(&AddConfig{[]string{"Logger"}, "app", "string", "test"})
			task.Exec(ctx)
			So(err, ShouldBeNil)

			// 执行不同的 struct
			task = NewAddTask(&AddConfig{[]string{"HTTP"}, "port", "int", "8080"})
			task.Exec(ctx)
			So(err, ShouldBeNil)

			// 任务执行完成后，也能通过 vet 检查
			err = util.Exec(pp, "go", "vet", "./...")
			So(err, ShouldBeNil)

			except := `HTTP:
  port: "8080"
Logger:
  app: test
  level: debug
`

			// 刚创建项目，两个 yaml 文件应该是空
			cc := string(util.ReadFile(filepath.Join(pp, "configs/config.yaml")))
			cec := string(util.ReadFile(filepath.Join(pp, "configs/config.yaml.example")))
			So(cc, ShouldEqual, except)
			So(cec, ShouldEqual, except)
			// go 文件中也没有 logger 的字样
			gc := string(util.ReadFile(filepath.Join(pp, "internal/conf/config.go")))
			So(strings.Contains(gc, "Logger"), ShouldBeTrue)
			So(strings.Contains(gc, "Level"), ShouldBeTrue)
			So(strings.Contains(gc, "App"), ShouldBeTrue)
			So(strings.Contains(gc, "HTTP"), ShouldBeTrue)
			So(strings.Contains(gc, "Port"), ShouldBeTrue)
		})

	})
}
