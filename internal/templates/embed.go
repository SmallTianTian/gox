package templates

import (
	"embed"
	"io/ioutil"

	"github.com/SmallTianTian/fresh-go/utils"
)

var ReadTemplateFile = readEmbedFile()

//go:embed */**
var ht embed.FS

func readEmbedFile() func(name string) string {
	return func(name string) string {
		f, err := ht.Open(name)
		utils.MustNotError(err)
		bs, err := ioutil.ReadAll(f)
		utils.MustNotError(err)
		return string(bs)
	}
}
