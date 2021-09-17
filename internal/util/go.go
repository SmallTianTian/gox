package util

import "tianxu.xin/gox/pkg/logger"

// GoFmtCode 格式化 golang 代码.
func GoFmtCode(path string) bool {
	if CheckCommandExists("goimports") {
		if err := Exec(path, "goimports", "-w", "."); err == nil {
			return true
		}
	}

	if err := Exec(path, "gofmt", "-s", "-w", "."); err != nil {
		logger.Errorf("Go fmt code failed.\n Please exec `gofmt -s -w .` in `%s`.\n Error: %v", path, err)
		return false
	}
	return true
}
