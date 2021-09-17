package ast

import (
	"io/fs"
	"os"

	"tianxu.xin/gox/internal/util"
)

type GoFile struct {
	path     string
	fi       fs.FileInfo
	lineNode *node
}

func ParseFile(path string) *GoFile {
	fi, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return nil
	}

	lines := util.FileEachLineWithTrim(path)
	return &GoFile{path: path, fi: fi, lineNode: newNode(lines)}
}

func (gf *GoFile) OverWrite() {
	if gf == nil {
		return
	}

	err := os.WriteFile(gf.path, []byte(gf.lineNode.String()), gf.fi.Mode())
	util.MustNotError(err)
}
