package utils

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func OverwritingFile(path, fpName, content string) {
	fullPath := filepath.Join(path, fpName)
	if IsExist(fullPath) {
		err := ioutil.WriteFile(fullPath, []byte(content), 0644)
		MustNotError(err)
	}
	err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	MustNotError(err)
	err = ioutil.WriteFile(fullPath, []byte(content), 0644)
	MustNotError(err)
}

func ReadFile(path string) []byte {
	bs, err := ioutil.ReadFile(path)
	MustNotError(err)
	return bs
}

func ReadTxtFileEachLine(path string) (lines []string) {
	r := bufio.NewReader(bytes.NewBuffer(ReadFile(path)))
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		MustNotError(err)
		lines = append(lines, string(line))
	}
	return
}

func File2GoAST(path string) *ast.File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	MustNotError(err)
	return f
}

func WriteAstFile(path, fpName string, f *ast.File) {
	fset := token.NewFileSet()
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := format.Node(buffer, fset, f)
	MustNotError(err)
	OverwritingFile(path, fpName, buffer.String())
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func WriteByTemplate(path string, fileAndTmpl map[string]string, keyAndRealValue map[string]interface{}) {
	for file, tmpl := range fileAndTmpl {
		realContent, err := StringFormat(tmpl, keyAndRealValue)
		MustNotError(err)
		OverwritingFile(path, file, realContent)
	}
}
