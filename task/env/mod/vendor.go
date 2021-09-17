package mod

import (
	"context"

	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task"
)

type GoModVendorTask struct {
	task.GoxAbstractTask
}

func NewGoModVendorTask() *GoModVendorTask {
	t := &GoModVendorTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName("go mod vendor")
	t.AppendCheck(util.FileInPath("go.mod"))

	return t
}

func (n *GoModVendorTask) Exec(ctx context.Context) error {
	return util.Exec(util.MustExtractProjectPath(ctx), "go", "mod", "vendor")
}
