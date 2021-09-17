package task

import (
	"context"
	"errors"

	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/internal/constant"
	ct "tianxu.xin/gox/pkg/context"
	"tianxu.xin/gox/pkg/logger"
)

type GoXRootTask struct {
	GoxAbstractTask

	ctx context.Context
	log logger.Logger
}

func New(conf *config.Config, projectPath string) *GoXRootTask {
	if conf == nil {
		panic("nil config")
	}

	log := logger.InitLogger(conf.Logger.Level)

	tagCtx := ct.NewTagContext(context.Background())
	ct.MustExtract(tagCtx).Set(constant.Logger, log)
	ct.MustExtract(tagCtx).Set(constant.Config, conf)
	ct.MustExtract(tagCtx).Set(constant.ProjectPath, projectPath)

	return &GoXRootTask{
		GoxAbstractTask: GoxAbstractTask{Name: "root"},
		ctx:             tagCtx,
		log:             log,
	}
}

func (root *GoXRootTask) Exec() error {
	if root == nil {
		return nil
	}

	if !root.PreCheck(root.ctx) {
		return errors.New("root check failed")
	}

	return root.AfterExec(root.ctx)
}
