package env

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/test"
	"tianxu.xin/gox/internal/util"
	ct "tianxu.xin/gox/pkg/context"
	"tianxu.xin/gox/task/env/mod"
)

func TestLoadProjectTask_Check(t *testing.T) {
	Convey("测试检查", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		lp := NewLoadProjectTask()

		Convey("不存在的地址，将无法通过检查", func() {
			ctx := ProjectCtxHelper(filepath.Join(dir, "tt"))

			So(lp.PreCheck(ctx), ShouldBeFalse)
		})

		Convey("地址存在，但不包含 go.mod 文件，将无法通过检查", func() {
			ctx := ProjectCtxHelper(filepath.Join(dir, "tt"), ResetCtx())

			So(lp.PreCheck(ctx), ShouldBeFalse)
		})

		Convey("地址存在，包含 go.mod 文件，将通过检查", func() {
			ctx := ProjectCtxHelper(filepath.Join(dir, "tt"),
				WithProject(), ResetCtx())

			So(lp.PreCheck(ctx), ShouldBeTrue)
		})
	})
}

func TestLoadProjectTask_Exec(t *testing.T) {
	Convey("测试载入项目", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		ctx := test.GetTestCtx(filepath.Join(dir, "tt"))

		lp := NewLoadProjectTask()

		Convey("正常 load 出 module", func() {
			// 调用 new project 任务创建一个项目
			module := "tianxu.xin/gox"
			util.MustExtractConf(ctx).GoEnv.Module = module
			NewFreshProjectTask().Exec(ctx)
			// 调用 mod init 任务，将生成 go.mod 文件
			mod.NewGoModInitTask().Exec(ctx)
			ct.MustExtract(ctx).Set(constant.Config, &config.Config{}) // reset config

			So(util.MustExtractConf(ctx).GoEnv.Module, ShouldBeEmpty)
			lp.Exec(ctx)
			So(util.MustExtractConf(ctx).GoEnv.Module, ShouldEqual, module)
		})

		Convey("正常 load 出 module 和 vendor", func() {
			// 调用 new project 任务创建一个项目
			module := "tianxu.xin/gox"
			util.MustExtractConf(ctx).GoEnv.Module = module
			NewFreshProjectTask().Exec(ctx)
			// 调用 mod init 任务，将生成 go.mod 文件
			mod.NewGoModInitTask().Exec(ctx)
			mod.NewGoModVendorTask().Exec(ctx)
			ct.MustExtract(ctx).Set(constant.Config, &config.Config{}) // reset config

			So(util.MustExtractConf(ctx).GoEnv.Module, ShouldBeEmpty)
			So(util.MustExtractConf(ctx).GoEnv.Vendor, ShouldBeFalse)
			lp.Exec(ctx)
			So(util.MustExtractConf(ctx).GoEnv.Module, ShouldEqual, module)
			So(util.MustExtractConf(ctx).GoEnv.Vendor, ShouldBeTrue)
		})
	})
}
