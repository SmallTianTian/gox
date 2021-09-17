package conf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/internal/util/ast"
	"tianxu.xin/gox/task"
)

type AddConfig struct {
	fathers          []string
	param, _type, dv string
}

func NewAddConfig(fathers []string, param, _type, defaultValue string) *AddConfig {
	return &AddConfig{
		fathers: fathers,
		param:   param,
		_type:   _type,
		dv:      defaultValue,
	}
}

type AddTask struct {
	task.GoxAbstractTask
	configs []*AddConfig
}

func NewAddTask(configs ...*AddConfig) *AddTask {
	t := &AddTask{
		GoxAbstractTask: task.GoxAbstractTask{},
		configs:         configs,
	}
	t.WithName("add config")
	t.AppendCheck(util.ProjectPathExist)
	t.AppendCheck(util.FileInPath("configs/config.yaml"))
	t.AppendCheck(util.FileInPath("configs/config.yaml.example"))
	t.AppendCheck(util.FileInPath("internal/conf/config.go"))
	t.AppendCheck(func(c context.Context) bool {
		for _, c := range t.configs {
			if !util.LegalVarName(c.param) {
				return false
			}
		}
		return true
	})
	return t
}

func (n *AddTask) Exec(ctx context.Context) error {
	dir := util.MustExtractProjectPath(ctx)
	if len(n.configs) == 0 {
		return nil
	}

	// 写入 config.go
	gf := ast.ParseFile(filepath.Join(dir, "internal/conf/config.go"))
	for _, c := range n.configs {
		// go file 默认填充 Config
		goFathers := append([]string{"Config"}, c.fathers...)
		gf.StructAddField(goFathers, "", fmt.Sprintf("%s %s", util.FirstUp(c.param), c._type))
	}
	gf.OverWrite()

	// 写入 yaml 文件
	yamlConfigs := []string{"configs/config.yaml.example", "configs/config.yaml"}
	for _, yc := range yamlConfigs {
		ycd := filepath.Join(dir, yc)
		bs := util.ReadFile(ycd)
		for _, c := range n.configs {
			bs = yaml(bs, c.param, c.dv, c.fathers)
		}
		err := os.WriteFile(ycd, bs, constant.WriteFileMode)
		util.MustNotError(err)
	}
	return n.AfterExec(ctx)
}
