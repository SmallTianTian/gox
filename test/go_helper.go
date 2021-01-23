package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/utils"
)

// 初始化 gomod
func InitGoMod(o, p, dir string) {
	if !utils.IsExist(dir) {
		utils.MustNotError(os.MkdirAll(dir, os.ModePerm))
	}

	path := filepath.Join(dir, "go.mod")
	content := fmt.Sprintf(`
module %s/%s

go 1.15`, o, p)
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
