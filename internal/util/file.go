package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// IsExist 判断目录/文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// ReadFile 读取文件，返回文件的字节数组
func ReadFile(path string) []byte {
	bs, err := ioutil.ReadFile(path)
	MustNotError(err)
	return bs
}

// FileEachLine 按行读取文本文件,
// 返回文本字符串行数组.
func FileEachLine(path string) (lines []string) {
	r := bufio.NewReader(bytes.NewBuffer(ReadFile(path)))
	for {
		line, _, err := r.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		MustNotError(err)
		lines = append(lines, string(line))
	}
	return
}

// FileEachLineWithTrim 按行读取文本文件，
// 并且将前后空格移除
func FileEachLineWithTrim(path string) (lines []string) {
	for _, line := range FileEachLine(path) {
		lines = append(lines, strings.TrimSpace(line))
	}
	return
}

// ListSubDir 列出目录下的所有文件夹/文件名称.
//
// 只是文件夹/文件名称，不包含前面的路径
func ListDir(path string, onlyDir bool) []string {
	fi, err := os.Stat(path)

	// 目录不存在
	if err != nil && !os.IsExist(err) {
		return nil
	}
	// 不是目录
	if !fi.IsDir() {
		return nil
	}

	entrys, err := os.ReadDir(path)
	MustNotError(err)

	var subs []string
	for _, e := range entrys {
		if onlyDir && !e.IsDir() {
			continue
		}
		subs = append(subs, e.Name())
	}
	return subs
}
