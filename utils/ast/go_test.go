package ast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"
)

func readTempAst() (f *ast.File) {
	fset := token.NewFileSet()
	var err error
	if f, err = parser.ParseFile(fset, "./temp_test.go", nil, parser.AllErrors); err != nil {
		panic(err)
	}
	return
}

func writeAst(name string, f *ast.File) error {
	fset := token.NewFileSet()
	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := format.Node(buffer, fset, f); err != nil {
		return err
	}
	return ioutil.WriteFile(name, buffer.Bytes(), 0644)
}

func Test_Struct(t *testing.T) {
	att := readTempAst()
	AddField2AstFile(att, "B", "bool", []string{"_struct1"})
	AddStruct2AstFile(att, "C", []string{""})
	writeAst("Test_Struct.go", att)
}

func Test_Field(t *testing.T) {
	att := readTempAst()
	AddField2AstFile(att, "B", "bool", []string{"_struct1"})
	AddField2AstFile(att, "C", "int", []string{"_struct1", "E"})
	writeAst("Test_Struct.go", att)
}

func Test_Func(t *testing.T) {
	att := readTempAst()
	AppendFuncCall2AstFile(att, "F1", []string{"name"}, []string{"F1"})
	writeAst("Test_Func.go", att)
}

func Test_Import(t *testing.T) {
	att := readTempAst()
	SetImport2AstFile(att, "math")
	writeAst("Test_Import.go", att)
}
