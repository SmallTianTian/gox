package test

import (
	"io/ioutil"
	"os"
)

// 临时目录
func TempDir() string {
	dir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		panic(err)
	}
	return dir
}

// 写文件
func WriteFile(path, content string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(content); err != nil {
		panic(err)
	}
}
