package project

import (
	"path/filepath"
	"strings"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/templates"
	"github.com/SmallTianTian/fresh-go/model"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
)

var base = []*model.FileTemp{
	{Name: "cmd/server/main.go", Content: templates.ReadTemplateFile("project/cmd/server/main.go.tmpl")},
	{Name: "cmd/server/wire.go", Content: templates.ReadTemplateFile("project/cmd/server/wire.go.tmpl")},
	{Name: "configs/config.yaml", Content: templates.ReadTemplateFile("project/configs/config.yaml.tmpl")},
	{Name: "githooks/commit-msg", Content: templates.ReadTemplateFile("project/githooks/commit-msg")},
	{Name: "githooks/go_pre_commit.sh", Content: templates.ReadTemplateFile("project/githooks/go_pre_commit.sh")},
	{Name: "githooks/pre-commit", Content: templates.ReadTemplateFile("project/githooks/pre-commit")},
	{Name: "internal/conf/config.go", Content: templates.ReadTemplateFile("project/internal/conf/config.go.tmpl")},
	{Name: "pkg/application/app_server.go", Content: templates.ReadTemplateFile("project/pkg/application/app_server.go.tmpl")}, // nolint
	{Name: "pkg/application/app.go", Content: templates.ReadTemplateFile("project/pkg/application/app.go.tmpl")},
	{Name: "pkg/logger/logger.go", Content: templates.ReadTemplateFile("project/pkg/logger/logger.go.tmpl")},
	{Name: ".gitignore", Content: templates.ReadTemplateFile("project/.gitignore.tmpl")},
	{Name: "Makefile", Content: templates.ReadTemplateFile("project/Makefile.tmpl")},
}

var log = map[string]*model.FileTemp{
	"zap":    {Name: "pkg/logger/zap.go", Content: templates.ReadTemplateFile("project/pkg/logger/zap.go.tmpl")},
	"logrus": {Name: "pkg/logger/logrus.go", Content: templates.ReadTemplateFile("project/pkg/logger/logrus.go.tmpl")},
}

func NewProject() {
	logger.Debug("Will begin create new project.")
	pro := config.DefaultConfig.Project.Name
	// 项目目录需要追加项目名
	dir := filepath.Join(config.DefaultConfig.Project.Path, pro)
	config.DefaultConfig.Project.Path = dir
	if utils.CheckGoProject(dir) {
		panic("Cloudn't init go project again in" + dir)
	}

	org := config.DefaultConfig.Project.Org
	isVendor := config.DefaultConfig.Project.Vendor
	module := filepath.Join(org, pro)
	logger.Debugf("Project name: %s\nOrganization: %s\nPath: %s\nUse model: %v.", pro, org, dir, isVendor)

	var kRv = map[string]interface{}{
		"module": module,
		"vendor": isVendor,
		"name":   pro,
	}

	// 0. 初始化项目基本结构
	utils.WriteByTemplate(dir, kRv, base...)

	// 1. 初始化日志
	initLog(dir)

	// 做一些项目初始化操作
	doSomeInit(dir)
}

// 初始化日志
func initLog(dir string) {
	switch strings.ToLower(config.DefaultConfig.Adapter.Logger) {
	case "logrus":
		utils.WriteByTemplate(dir, nil, log["logrus"])
	case "zap":
		fallthrough
	default:
		utils.WriteByTemplate(dir, nil, log["zap"])
	}
}

func doSomeInit(dir string) {
	// 初始化 mod
	utils.FirstMod(config.DefaultConfig.Project.Path,
		filepath.Join(config.DefaultConfig.Project.Org, config.DefaultConfig.Project.Name),
		config.DefaultConfig.Project.Vendor)

	utils.MustNotError(utils.Exec(dir, "wire", "./..."))
	utils.MustNotError(utils.Exec(dir, "git", "init"))
	utils.MustNotError(utils.Exec(dir, "chmod", "+x", "githooks/go_pre_commit.sh"))
	utils.MustNotError(utils.Exec(dir, "cp", "githooks/pre-commit", ".git/hooks"))
	utils.MustNotError(utils.Exec(dir, "chmod", "+x", "./.git/hooks/pre-commit"))
	utils.MustNotError(utils.Exec(dir, "cp", "githooks/commit-msg", ".git/hooks"))
	utils.MustNotError(utils.Exec(dir, "chmod", "+x", "./.git/hooks/commit-msg"))
}
