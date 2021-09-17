package util

import (
	"context"
	"path/filepath"

	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/internal/constant"
	ct "tianxu.xin/gox/pkg/context"
	"tianxu.xin/gox/pkg/logger"
)

func MustExtractLog(ctx context.Context) logger.Logger {
	logInter := ct.MustExtract(ctx).Get(constant.Logger)

	log, ok := logInter.(logger.Logger)
	if !ok {
		panic("no log instance")
	}
	return log
}

func MustExtractConf(ctx context.Context) *config.Config {
	confInter := ct.MustExtract(ctx).Get(constant.Config)

	conf, ok := confInter.(*config.Config)
	if !ok {
		panic("no config instance")
	}
	return conf
}

func MustExtractProjectPath(ctx context.Context) string {
	ppInter := ct.MustExtract(ctx).Get(constant.ProjectPath)

	pp, ok := ppInter.(string)
	if !ok {
		panic("no project path")
	}
	return pp
}

func MustGetAbsolutePath(ctx context.Context, fullFileName string) string {
	return filepath.Join(MustExtractProjectPath(ctx), fullFileName)
}

func ExtractModRequire(ctx context.Context) []string {
	mInter := ct.MustExtract(ctx).Get(constant.ModRequire)
	if mInter == nil {
		return nil
	}

	pp, ok := mInter.([]string)
	if !ok {
		panic("no mod require")
	}
	return pp
}

func AddModRequire(ctx context.Context, require string) {
	m := ExtractModRequire(ctx)
	ct.MustExtract(ctx).Set(constant.ModRequire, append(m, require))
}
