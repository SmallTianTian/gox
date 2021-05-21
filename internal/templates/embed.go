package templates

import (
	"embed"
	"io/ioutil"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/utils"
)

var ReadTemplateFile = readEmbedFile()

//go:embed project/**
var ht embed.FS

func readEmbedFile() func(name string) string {
	return func(name string) string {
		name = filepath.Join("project", name)

		f, err := ht.Open(name)
		utils.MustNotError(err)
		bs, err := ioutil.ReadAll(f)
		utils.MustNotError(err)
		return string(bs)
	}
}
