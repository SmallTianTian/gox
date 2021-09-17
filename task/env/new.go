package env

import (
	"context"
	"os"
	"path/filepath"

	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/templates"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task"
)

// FreshProjectTask 新项目任务.
// 需要前置检查：该地址没有该项目
// 主要做以下几件事情：
// 1. 创建项目目录
// 2. 创建一个可运行项目基本的架构
type FreshProjectTask struct {
	task.GoxAbstractTask
	tmpFiles []*templates.FileTemp
}

func NewFreshProjectTask() *FreshProjectTask {
	t := &FreshProjectTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName("new project")
	t.AppendCheck(util.Not(util.ProjectPathExist))

	t.tmpFiles = []*templates.FileTemp{
		templates.NewEmbedTemp("cmd/server/main.go"),
		templates.NewEmbedTemp("configs/config.yaml"),
		templates.NewEmbedTemp("configs/config.yaml.example"),
		templates.NewEmbedTemp("internal/conf/config.go"),
		templates.NewEmbedTemp("pkg/application/app_server.go"),
		templates.NewEmbedTemp("pkg/application/app.go"),
		templates.NewEmbedTemp("pkg/logger/logger.go"),
		templates.NewEmbedTemp(".golangci.yml"),
		templates.NewEmbedTemp("Makefile"),
	}

	return t
}

func (n *FreshProjectTask) Exec(ctx context.Context) error {
	pp := util.MustExtractProjectPath(ctx)
	// 先创建项目目录
	err := os.MkdirAll(pp, constant.MkdirMode)
	util.MustNotError(err)

	// 创建一个项目的基本架构
	for _, tmp := range n.tmpFiles {
		tmp.Write(ctx, map[string]interface{}{
			constant.TempModule:  util.MustExtractConf(ctx).GoEnv.Module,
			constant.ProjectName: filepath.Base(util.MustExtractConf(ctx).GoEnv.Module), // 项目名称一般为 module 最后一段
		})
	}
	return n.AfterExec(ctx)
}
