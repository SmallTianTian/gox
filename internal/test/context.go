package test

import (
	"context"

	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/internal/constant"
	ct "tianxu.xin/gox/pkg/context"
	"tianxu.xin/gox/pkg/logger"
)

func GetTestCtx(projectPath string) context.Context {
	ctx := context.Background()
	ctx = ct.NewTagContext(ctx)
	tag := ct.MustExtract(ctx)

	tag.Set(constant.ProjectPath, projectPath)
	tag.Set(constant.Logger, logger.InitLogger("debug"))
	tag.Set(constant.Config, config.GetDefaultConfig())
	return ctx
}
