package project

import (
	"path/filepath"
	"runtime"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/templates"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
)

var fileAndTmpl = map[string]string{
	"main.go":              templates.ReadTemplateFile("project/main.go.tmpl"),
	"Dockerfile":           templates.ReadTemplateFile("project/Dockerfile.tmpl"),
	"Makefile":             templates.ReadTemplateFile("project/Makefile.tmpl"),
	"cmd/base.go":          templates.ReadTemplateFile("project/cmd/base.go.tmpl"),
	"cmd/root.go":          templates.ReadTemplateFile("project/cmd/root.go.tmpl"),
	"config/config.go":     templates.ReadTemplateFile("project/config/config.go.tmpl"),
	"config/config.yaml":   templates.ReadTemplateFile("project/config/config.yaml.tmpl"),
	"pkg/logger/logger.go": templates.ReadTemplateFile("project/pkg/logger/logger.go.tmpl"),
}

func NewProject() {
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
		"module":      module,
		"projectName": pro,
		"goVersion":   runtime.Version()[2:],
	}
	if kRv["modVendor"] = ""; isVendor {
		kRv["modVendor"] = "-mod=vendor"
	}
	utils.WriteByTemplate(dir, fileAndTmpl, kRv)
}
