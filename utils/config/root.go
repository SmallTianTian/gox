package config_util

import (
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/utils"
)

var configs = map[string]func(content []byte, k string, v interface{}, father []string) []byte{
	"config/config.yaml": yaml,
}

func WriteConfig(path, k string, v interface{}, father []string) {
	for cf, f := range configs {
		file := filepath.Join(path, cf)
		if utils.IsExist(file) {
			bs := f(utils.ReadFile(file), k, v, father)
			utils.OverwritingFile(file, "", string(bs))
		}
	}
}
