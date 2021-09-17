package templates

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"text/template"

	"tianxu.xin/gox/internal/constant"
	"tianxu.xin/gox/internal/util"
)

type FileTemp struct {
	path string
	name string
	temp *template.Template
}

// NewEmbedTemp 从 embed 文件中生成实例
// embed 文件名为 name 追加 .tmpl 作为文件地址
func NewEmbedTemp(name string) *FileTemp {
	return NewTemp(name, File(name+".tmpl"))
}

func NewTemp(name, content string) *FileTemp {
	t, err := template.New("").Parse(content)
	util.MustNotError(err)
	return &FileTemp{
		path: filepath.Dir(name),
		name: filepath.Base(name),
		temp: t,
	}
}

func (ft *FileTemp) Write(ctx context.Context, kv map[string]interface{}) {
	if ft == nil {
		return
	}

	// 先看看目录是否存在，不存在先创建目录
	fp := util.MustGetAbsolutePath(ctx, ft.path)
	if !util.IsExist(fp) {
		err := os.MkdirAll(fp, constant.MkdirMode)
		util.MustNotError(err)
	}

	var bs bytes.Buffer
	err := ft.temp.Execute(&bs, kv)
	util.MustNotError(err)

	// 后写文件
	filePath := util.MustGetAbsolutePath(ctx, filepath.Join(ft.path, ft.name))
	err = os.WriteFile(filePath, bs.Bytes(), constant.WriteFileMode)
	util.MustNotError(err)
}
