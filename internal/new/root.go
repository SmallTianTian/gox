package new

import (
	"path/filepath"
	"runtime"

	"github.com/SmallTianTian/fresh-go/utils"
)

var fileAndTmpl = map[string]string{
	"main.go":              utils.ReadStatikFile("/project/main.go.tmpl"),
	"Dockerfile":           utils.ReadStatikFile("/project/Dockerfile.tmpl"),
	"Makefile":             utils.ReadStatikFile("/project/Makefile.tmpl"),
	"cmd/base.go":          utils.ReadStatikFile("/project/cmd/base.go.tmpl"),
	"cmd/root.go":          utils.ReadStatikFile("/project/cmd/root.go.tmpl"),
	"config/config.go":     utils.ReadStatikFile("/project/config/config.go.tmpl"),
	"config/config.yaml":   utils.ReadStatikFile("/project/config/config.yaml.tmpl"),
	"pkg/logger/logger.go": utils.ReadStatikFile("/project/pkg/logger/logger.go.tmpl"),
}

func NewProject(path, organization string, isVendor bool) (err error) {
	module := filepath.Join(organization, filepath.Base(path))
	var kRv = map[string]interface{}{
		"module":      module,
		"projectName": filepath.Base(path),
		"goVersion":   runtime.Version()[2:],
	}
	if kRv["modVendor"] = ""; isVendor {
		kRv["modVendor"] = "-mod=vendor"
	}
	utils.WriteByTemplate(path, fileAndTmpl, kRv)
	return nil
}
