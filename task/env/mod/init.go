package mod

import (
	"context"

	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task"
)

var (
	moduleCheck = func(ctx context.Context) bool {
		config := util.MustExtractConf(ctx)
		return config.GoEnv.Module != ""
	}
)

type GoModInitTask struct {
	task.GoxAbstractTask
}

func NewGoModInitTask() *GoModInitTask {
	t := &GoModInitTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName("go mod init")
	t.AppendCheck(util.ProjectPathExist)
	t.AppendCheck(moduleCheck)

	t.AppendChildTask(NewGoModTidyTask())
	return t
}

func (n *GoModInitTask) Exec(ctx context.Context) error {
	if n == nil {
		return nil
	}

	path := util.MustExtractProjectPath(ctx)
	config := util.MustExtractConf(ctx)

	if err := util.Exec(path, "go", "mod", "init", config.GoEnv.Module); err != nil {
		return err
	}

	return n.AfterExec(ctx)
}
