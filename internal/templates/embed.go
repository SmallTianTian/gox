package templates

import (
	"embed"
	"io/ioutil"
	"path/filepath"

	"tianxu.xin/gox/internal/util"
)

var ReadTemplateFile = readEmbedFile()
var File = ReadTemplateFile

//go:embed project/**
var ht embed.FS

func readEmbedFile() func(name string) string {
	return func(name string) string {
		name = filepath.Join("project", name)

		f, err := ht.Open(name)
		util.MustNotError(err)
		bs, err := ioutil.ReadAll(f)
		util.MustNotError(err)
		return string(bs)
	}
}
