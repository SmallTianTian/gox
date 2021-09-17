package git

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/templates"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task"
)

// InitTask 初始化 git 任务.
// 需要前置检查：
// 1. 该地址没有该项目
// 2. 该地址是否已经初始化 git
// 3. 该项目中是否已经存在 githooks
// 主要做以下几件事情：
// 1. 执行 git init
// 2. 增加部分 git hook
// 3. 修改 Makefile 文件
type InitTask struct {
	task.GoxAbstractTask
	tmpFiles []*templates.FileTemp
}

func NewInitTask() *InitTask {
	t := &InitTask{GoxAbstractTask: task.GoxAbstractTask{}}
	t.WithName("git init")
	t.AppendCheck(util.ProjectPathExist)
	t.AppendCheck(util.Not(util.FileInPath(".git")))
	t.AppendCheck(util.Not(util.FileInPath("githooks")))

	t.tmpFiles = []*templates.FileTemp{
		templates.NewEmbedTemp("githooks/commit-msg"),
		templates.NewEmbedTemp("githooks/go_pre_commit.sh"),
		templates.NewEmbedTemp("githooks/pre-commit"),
		templates.NewEmbedTemp(".gitignore"),
	}

	return t
}

func (n *InitTask) Exec(ctx context.Context) error {
	dir := util.MustExtractProjectPath(ctx)
	// 执行 git init
	err := util.Exec(dir, "git", "init")
	util.MustNotError(err)

	// 先创建 githooks 目录
	err = os.MkdirAll(filepath.Join(dir, "githooks"), constant.MkdirMode)
	util.MustNotError(err)

	// 将 git hook 文件写入项目中
	for _, tmp := range n.tmpFiles {
		tmp.Write(ctx, map[string]interface{}{})
	}

	// 在 git 中加入 hook
	util.MustNotError(util.Exec(dir, "chmod", "+x", "githooks/go_pre_commit.sh"))
	util.MustNotError(util.Exec(dir, "cp", "githooks/pre-commit", ".git/hooks"))
	util.MustNotError(util.Exec(dir, "chmod", "+x", ".git/hooks/pre-commit"))
	util.MustNotError(util.Exec(dir, "cp", "githooks/commit-msg", ".git/hooks"))
	util.MustNotError(util.Exec(dir, "chmod", "+x", ".git/hooks/commit-msg"))

	// 修改 Makefile 文件
	lines := util.FileEachLine(filepath.Join(dir, "Makefile"))
	var bb bytes.Buffer
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "VERSION="):
			bb.WriteString("VERSION=$(shell git describe --tags --always)")
		case strings.HasPrefix(line, "GO_FLAGS=") && strings.Contains(line, "no git"):
			pl := strings.Split(line, "no git")
			bb.WriteString(pl[0])
			bb.WriteString("`git rev-parse HEAD`")
			bb.WriteString(pl[1])
		default:
			bb.WriteString(line)
		}
		bb.WriteString("\n")
	}

	err = os.WriteFile(filepath.Join(dir, "Makefile"), bb.Bytes(), constant.WriteFileMode)
	util.MustNotError(err)

	return n.AfterExec(ctx)
}
