// 由于日志有太大的共性
// 抽出来统一写
package log

import (
	"context"
	"fmt"
	"path/filepath"

	"tianxu.xin/gox/internal/templates"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/internal/util/ast"
	"tianxu.xin/gox/task"
	"tianxu.xin/gox/task/env/conf"
)

// commonTask 新增日志.
// 需要前置检查：
// 1. 该地址是否是 go 项目
// 2. 该项目存在 pkg/logger 目录
// 3. 该项目中不存在对应的日志文件
//
// 主要做以下几件事情：
// 1. 新建日志文件
// 2. 如果不存在 log_option 文件，则新建。
// 3. 在 cmd 各目录 main.go 文件中，写入该日志初始化
type commonTask struct {
	task.GoxAbstractTask
	tmpFiles []*templates.FileTemp
}

func newCommonTask(name string) *commonTask {
	t := &commonTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName(name)
	t.AppendCheck(util.ProjectPathExist)
	t.AppendCheck(util.FileInPath("pkg/logger"))

	t.tmpFiles = []*templates.FileTemp{
		templates.NewEmbedTemp(fmt.Sprintf("pkg/logger/%s.go", name)),
	}
	return t
}

func (n *commonTask) Exec(ctx context.Context) error {
	// 将必要的文件先写
	for _, tmp := range n.tmpFiles {
		tmp.Write(ctx, map[string]interface{}{})
	}

	// 如果没有 logoption 文件，则写入
	if !util.FileInPath(LogOptionPath)(ctx) {
		templates.NewEmbedTemp(LogOptionPath).Write(ctx, map[string]interface{}{})
	}

	dir := util.MustExtractProjectPath(ctx)

	subDirs := util.ListDir(filepath.Join(dir, "cmd"), true)
	for _, sd := range subDirs {
		// 只看各个目录下的 main.go 文件，
		// 如果 main 函数不在 main.go 里面，
		// 则不予理会
		subMain := filepath.Join(dir, "cmd", sd, "main.go")
		if !util.IsExist(subMain) {
			continue
		}

		// 由于大多数
		gf := ast.ParseFile(subMain)
		gf.FuncAddCall("main", "conf.InitConfig()", fmt.Sprintf("初始化 %s 日志", n.Name),
			fmt.Sprintf("logger.Init%s(config.Logger.App, config.Logger.Level)",
				util.FirstUp(n.Name)), false)
		gf.OverWrite()
	}

	config := util.MustExtractConf(ctx)

	n.ChildsTask = append(n.ChildsTask, conf.NewAddTask(
		conf.NewAddConfig([]string{"Logger"}, "Level", "string", "debug"),
		conf.NewAddConfig([]string{"Logger"}, "App", "string", filepath.Base(config.GoEnv.Module)),
	))
	return n.AfterExec(ctx)
}
