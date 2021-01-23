package utils

import (
	"io/ioutil"

	// Register some standard stuff
	_ "github.com/SmallTianTian/fresh-go/statik"
	"github.com/rakyll/statik/fs"
)

var ReadStatikFile = readStatikFile()

func readStatikFile() func(name string) string {
	ht, err := fs.New()
	MustNotError(err)
	return func(name string) string {
		f, err := ht.Open(name)
		MustNotError(err)
		bs, err := ioutil.ReadAll(f)
		MustNotError(err)
		return string(bs)
	}
}
