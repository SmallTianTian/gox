package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SmallTianTian/fresh-go/utils"
)

// 初始化 gomod
func InitGoMod(remote, owner, name, dir string) {
	if !utils.IsExist(dir) {
		utils.MustNotError(os.MkdirAll(dir, os.ModePerm))
	}

	var mods []string
	for _, v := range []string{remote, owner, name} {
		if strings.TrimSpace(v) != "" {
			mods = append(mods, v)
		}
	}

	path := filepath.Join(dir, "go.mod")
	content := fmt.Sprintf(`
module %s

go 1.15`, strings.Join(mods, "/"))
	WriteFile(path, content)
}

// 初始化 vendor，
// isDir 为 true，创建 vendor 目录，为 false，创建 vendor 文件
func InitVendor(dir string, isDir bool) {
	vendor := filepath.Join(dir, "vendor")
	if isDir {
		if err := os.Mkdir(vendor, os.ModePerm); err != nil {
			panic(err)
		}
		return
	}
	WriteFile(vendor, "This is a file.")
}
