package ast

import (
	"go/ast"
	"go/token"
)

func newDeclStruct(name string) *ast.GenDecl {
	return &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{newSpecStruct(name)},
	}
}

func newSpecStruct(name string) *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{Name: name},
		Type: newTypeStruct(),
	}
}

func newField(name, _type string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{{Name: name}},
		Type:  ast.NewIdent(_type),
	}
}

func newFieldStruct(name string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{{Name: name}},
		Type:  newTypeStruct(),
	}
}

func newTypeStruct() *ast.StructType {
	return &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{},
		},
	}
}
