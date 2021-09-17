package git

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/test"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task/env"
)

func TestInitTask_Check(t *testing.T) {
	Convey("测试检查", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		task := NewInitTask()

		Convey("不存在的地址，将无法通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"))

			So(task.PreCheck(ctx), ShouldBeFalse)
		})

		Convey("项目存在，不存在 .git 目录，将通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())

			So(task.PreCheck(ctx), ShouldBeTrue)
		})

		Convey("项目存在，包含 .git 目录，将无法通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
			dir := util.MustExtractProjectPath(ctx)

			// 新建 .git 目录
			err := os.MkdirAll(filepath.Join(dir, ".git"), constant.MkdirMode)
			So(err, ShouldBeNil)

			So(task.PreCheck(ctx), ShouldBeFalse)
		})

		Convey("项目存在，包含 githooks 目录，将无法通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
			dir := util.MustExtractProjectPath(ctx)

			// 新建 .git 目录
			err := os.MkdirAll(filepath.Join(dir, "githooks"), constant.MkdirMode)
			So(err, ShouldBeNil)

			So(task.PreCheck(ctx), ShouldBeFalse)
		})
	})
}

func TestInitTask_Exec(t *testing.T) {
	Convey("测试 GitInit 执行", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理
		projectPath := filepath.Join(dir, "tt")

		// 新生成项目
		ctx := env.ProjectCtxHelper(projectPath, env.WithProject())
		task := NewInitTask()

		Convey("正常初始化 git 及添加 hooks", func() {
			// .git 目录应该不存在
			So(util.IsExist(filepath.Join(projectPath, ".git")), ShouldBeFalse)
			// githooks 目录应该不存在
			So(util.IsExist(filepath.Join(projectPath, "githooks")), ShouldBeFalse)

			err := task.Exec(ctx)
			So(err, ShouldBeNil)

			// 执行 git init 任务后 .git 目录应该存在
			So(util.IsExist(filepath.Join(projectPath, ".git")), ShouldBeTrue)
			// 执行 git init 任务后 githooks 目录应该存在
			So(util.IsExist(filepath.Join(projectPath, "githooks")), ShouldBeTrue)
		})

		Convey("初始化 git 后，git 功能正常，hook 生效", func() {
			// 执行任务前，使用 git 命令会报错
			err := util.Exec(projectPath, "git", "add", ".")
			So(err, ShouldNotBeNil)

			err = task.Exec(ctx)
			So(err, ShouldBeNil)

			// 执行任务后，能正常使用 git 命令
			err = util.Exec(projectPath, "git", "add", ".")
			So(err, ShouldBeNil)

			// commit 格式会被检查
			err = util.Exec(projectPath, "git", "commit", "-m", "bad msg")
			So(err, ShouldNotBeNil)

			// 正常的格式可以通过检查，代表 hook 生效
			err = util.Exec(projectPath, "git", "commit", "-m", "feat: first commit")
			So(err, ShouldBeNil)
		})
	})
}
