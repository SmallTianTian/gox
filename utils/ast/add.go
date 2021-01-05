package ast_util

import (
	"fmt"
	"go/ast"
	"go/token"
)

func AddStruct2AstFile(astF *ast.File, structName string, father []string) {
	node, residue := searchNodeInAstFile(astF, father)
	for i, r := range residue {
		AddStruct2AstFile(astF, r, father[:len(father)-len(residue)+i])
	}
	if node == nil {
		astF.Decls = append(astF.Decls, newDeclStruct(structName))
		return
	}

	var _struct *ast.StructType
Match_Type:
	switch s := node.(type) {
	case *ast.StructType:
		_struct = s
	case *ast.TypeSpec:
		node = s.Type
		goto Match_Type
	}
	_struct.Fields.List = append(_struct.Fields.List, newFieldStruct(structName))
}

func AddField2AstFile(astF *ast.File, fn, _type string, father []string) {
	node, residue := searchNodeInAstFile(astF, father)
	for i, r := range residue {
		AddStruct2AstFile(astF, r, father[:len(father)-len(residue)+i])
	}
	if len(residue) != 0 {
		// reserch
		node, _ = searchNodeInAstFile(astF, father)
	}

	if node == nil {
		astF.Decls = append(astF.Decls)
		return
	}

	var _struct *ast.StructType
Match_Type:
	switch s := node.(type) {
	case *ast.StructType:
		_struct = s
	case *ast.TypeSpec:
		node = s.Type
		goto Match_Type
	case *ast.Field:
		node = s.Type
		goto Match_Type
	}
	_struct.Fields.List = append(_struct.Fields.List, newField(fn, _type))
}

func AppendFuncCall2AstFile(astF *ast.File, fn string, pNames, father []string) {
	node, residue := searchNodeInAstFile(astF, father)
	if len(residue) != 0 || node == nil {
		return
	}

	switch f := node.(type) {
	case *ast.FuncDecl:
		x := &ast.CallExpr{Fun: ast.NewIdent(fn)}
		stmt := &ast.ExprStmt{X: x}
		for _, pn := range pNames {
			x.Args = append(x.Args, ast.NewIdent(pn))
		}
		f.Body.List = append(f.Body.List, stmt)
	}
}

func SetImport2AstFile(astF *ast.File, imp string) {
	var impdec *ast.GenDecl
	for _, decl := range astF.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			impdec = gd
			break
		}
	}

	impor := &ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf(`"%s"`, imp)}}
	if impdec != nil {
		for _, sp := range impdec.Specs {
			if is, ok := sp.(*ast.ImportSpec); ok {
				if is.Path.Value == imp {
					return
				}
			}
		}
		impdec.Specs = append(impdec.Specs, impor)
		return
	}

	gd := &ast.GenDecl{Tok: token.IMPORT, Specs: []ast.Spec{impor}}
	astF.Decls = append([]ast.Decl{gd}, astF.Decls...)
}
