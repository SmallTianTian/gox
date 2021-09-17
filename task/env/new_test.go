package env

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"tianxu.xin/gox/internal/test"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task/env/mod"
)

func TestFreshProjectTask_Check(t *testing.T) {
	Convey("测试检查", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		ctx := test.GetTestCtx(filepath.Join(dir, "tt"))
		fpt := NewFreshProjectTask()

		Convey("目录不存在的情况下，检查通过", func() {
			flag := fpt.PreCheck(ctx)
			So(flag, ShouldBeTrue)
		})

		Convey("目录已存在的情况下，检查未通过", func() {
			err := os.Mkdir(util.MustExtractProjectPath(ctx), fs.ModeDir)
			So(err, ShouldBeNil)

			flag := fpt.PreCheck(ctx)
			So(flag, ShouldBeFalse)
		})
	})
}

func TestFreshProjectTask_Exec(t *testing.T) {
	Convey("测试新建项目", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		ctx := test.GetTestCtx(filepath.Join(dir, "demo"))
		config := util.MustExtractConf(ctx)
		config.GoEnv.Module = defaultModule

		fpt := NewFreshProjectTask()

		Convey("正常执行不报错，但无法通过检查", func() {
			pp := util.MustExtractProjectPath(ctx)
			So(util.IsExist(pp), ShouldBeFalse)

			err := fpt.Exec(ctx)
			So(err, ShouldBeNil)
			So(util.IsExist(pp), ShouldBeTrue)
			err = util.Exec(pp, "go", "vet", "./...")
			So(err, ShouldBeError)
		})

		Convey("执行 mod 任务后，能通过检查", func() {
			pp := util.MustExtractProjectPath(ctx)
			So(util.IsExist(pp), ShouldBeFalse)

			err := fpt.Exec(ctx)
			So(err, ShouldBeNil)
			So(util.IsExist(pp), ShouldBeTrue)

			mt := mod.NewGoModInitTask()
			mt.Exec(ctx)
			err = util.Exec(pp, "go", "vet", "./...")
			So(err, ShouldBeNil)
		})

		Convey("执行 mod、vendor 任务后，能通过检查，且有 vendor 目录", func() {
			pp := util.MustExtractProjectPath(ctx)
			So(util.IsExist(pp), ShouldBeFalse)

			err := fpt.Exec(ctx)
			So(err, ShouldBeNil)
			So(util.IsExist(pp), ShouldBeTrue)

			mt := mod.NewGoModInitTask()
			mt.Exec(ctx)
			err = util.Exec(pp, "go", "vet", "./...")
			So(err, ShouldBeNil)
			So(util.FileInPath("vendor")(ctx), ShouldBeFalse)

			vt := mod.NewGoModVendorTask()
			vt.Exec(ctx)
			err = util.Exec(pp, "go", "vet", "./...")
			So(err, ShouldBeNil)
			So(util.FileInPath("vendor")(ctx), ShouldBeTrue)
		})
	})
}
