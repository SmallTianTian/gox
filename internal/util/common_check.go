package util

import (
	"context"
	"path/filepath"
)

// Not 将原有条件取反，避免再写一次
func Not(f func(context.Context) bool) func(context.Context) bool {
	return func(c context.Context) bool { return !f(c) }
}

// FileInPath 上下文中包含的项目地址中是否存在该文件
func FileInPath(file string) func(context.Context) bool {
	return func(c context.Context) bool {
		dir := MustExtractProjectPath(c)
		fileRealPath := filepath.Join(dir, file)
		return IsExist(fileRealPath)
	}
}

// ProjectPathExist 上下文中包含的项目地址是否存在
func ProjectPathExist(ctx context.Context) bool {
	return FileInPath(".")(ctx)
}

// GoModExist 上下文中包含的项目地址中是否存在 go.mod 文件
func GoModExist(ctx context.Context) bool {
	return FileInPath("go.mod")(ctx)
}
