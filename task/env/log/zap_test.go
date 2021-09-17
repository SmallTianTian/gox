package log

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"tianxu.xin/gox/internal/test"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task/env"
	"tianxu.xin/gox/task/env/mod"
)

func TestZapTask_Check(t *testing.T) {
	Convey("测试检查", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理

		task := NewZapTask()

		Convey("不存在的地址，将无法通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"))

			So(task.PreCheck(ctx), ShouldBeFalse)
		})

		Convey("项目存在，包含 go.mod 文件，将通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())

			So(task.PreCheck(ctx), ShouldBeTrue)
		})

		Convey("项目存在，但不包含 go.mod 文件，将无法通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"))

			So(task.PreCheck(ctx), ShouldBeFalse)
		})

		Convey("项目存在，包含 zap 文件，将无法通过检查", func() {
			ctx := env.ProjectCtxHelper(filepath.Join(dir, "tt"), env.WithProject())
			dir := util.MustExtractProjectPath(ctx)

			// 写入 zap 文件
			err := ioutil.WriteFile(filepath.Join(dir, "pkg/logger/zap.go"), []byte("not blank"), 0600)
			So(err, ShouldBeNil)

			So(task.PreCheck(ctx), ShouldBeTrue)
		})
	})
}

func TestZapTask_Exec(t *testing.T) {
	Convey("测试 zap 执行", t, func() {
		dir := test.TempDir()
		defer os.RemoveAll(dir) // 别忘记清理
		projectPath := filepath.Join(dir, "tt")

		// 新生成项目
		ctx := env.ProjectCtxHelper(projectPath, env.WithProject())
		task := NewZapTask()

		Convey("正常初始化 zap 日志及 option", func() {
			zapFileExistF := util.FileInPath("pkg/logger/zap.go")
			optFileExistF := util.FileInPath(LogOptionPath)

			// 最开始，zap 文件和 opt 文件应该不存在
			So(zapFileExistF(ctx), ShouldBeFalse)
			So(optFileExistF(ctx), ShouldBeFalse)
			// 没执行 zap task 前，main 文件中不会包含 Initzap
			content := string(util.ReadFile(filepath.Join(util.MustExtractProjectPath(ctx), "cmd/server/main.go")))
			So(strings.Contains(content, "InitZap"), ShouldBeFalse)

			// 执行 zap 任务
			err := task.Exec(ctx)
			So(err, ShouldBeNil)

			// 执行任务后，zap 文件和 opt 文件应该存在
			So(zapFileExistF(ctx), ShouldBeTrue)
			So(optFileExistF(ctx), ShouldBeTrue)
			// 执行 zap task 后，main 文件包含 Initzap
			content = string(util.ReadFile(filepath.Join(util.MustExtractProjectPath(ctx), "cmd/server/main.go")))
			So(strings.Contains(content, "InitZap"), ShouldBeTrue)

			// 执行 tidy 任务，方便后续 vet 检查
			err = mod.NewGoModTidyTask().Exec(ctx)
			So(err, ShouldBeNil)

			// vet 也不能出错
			err = util.Exec(util.MustExtractProjectPath(ctx), "go", "vet", "./...")
			So(err, ShouldBeNil)
		})

		Convey("option 存在，不更改，正常初始化 zap 日志", func() {
			const content = "not real body"
			optFilePath := filepath.Join(projectPath, LogOptionPath)
			// 写入 opt 文件
			err := ioutil.WriteFile(optFilePath, []byte(content), 0600)
			So(err, ShouldBeNil)

			zapFileExistF := util.FileInPath("pkg/logger/zap.go")
			optFileExistF := util.FileInPath(LogOptionPath)

			// 最开始，zap 文件和 opt 文件应该不存在
			So(zapFileExistF(ctx), ShouldBeFalse)
			So(optFileExistF(ctx), ShouldBeTrue)

			// 执行 zap 任务
			err = task.Exec(ctx)
			So(err, ShouldBeNil)

			// 执行任务后，zap 文件和 opt 文件应该存在
			So(zapFileExistF(ctx), ShouldBeTrue)
			So(optFileExistF(ctx), ShouldBeTrue)

			// 读出来的内容应该相同，
			// 不会改变
			current := string(util.ReadFile(optFilePath))
			So(content, ShouldEqual, current)
		})
	})
}
