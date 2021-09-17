package mod

import (
	"context"

	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task"
)

type GoModTidyTask struct {
	task.GoxAbstractTask
}

func NewGoModTidyTask() *GoModTidyTask {
	t := &GoModTidyTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName("go mod tidy")
	t.AppendCheck(util.FileInPath("go.mod"))

	return t
}

func (n *GoModTidyTask) Exec(ctx context.Context) error {
	if n == nil {
		return nil
	}

	path := util.MustExtractProjectPath(ctx)
	if err := util.Exec(path, "go", "mod", "tidy"); err != nil {
		return err
	}
	return n.AfterExec(ctx)
}
