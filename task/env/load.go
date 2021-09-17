package env

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/util"
	ct "tianxu.xin/gox/pkg/context"
	"tianxu.xin/gox/task"
)

// LoadProjectTask 获取项目信息的任务.
// 需要前置检查：该地址存在 go 项目
// 主要获取项目基本的架构
type LoadProjectTask struct {
	task.GoxAbstractTask
}

func NewLoadProjectTask() *LoadProjectTask {
	t := &LoadProjectTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName("load project")
	t.AppendCheck(util.ProjectPathExist)
	t.AppendCheck(util.GoModExist)
	return t
}

func (n *LoadProjectTask) Exec(ctx context.Context) error {
	const moduleStr = "module "

	config := &config.Config{}

	pwd := util.MustExtractProjectPath(ctx)

	// 1. 获取项目 module
	lines := util.FileEachLine(filepath.Join(pwd, "go.mod"))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, moduleStr) {
			config.GoEnv.Module = strings.TrimLeft(line, moduleStr)
			break
		}
	}
	if config.GoEnv.Module == "" {
		return errors.New("not go project")
	}

	// 2. 获取是否使用 vendor
	config.GoEnv.Vendor = util.FileInPath("vendor")(ctx)

	ct.MustExtract(ctx).Set(constant.Config, config)
	return nil
}
