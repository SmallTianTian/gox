package ast_util

import (
	"go/ast"
)

func searchNodeInAstFile(f *ast.File, father []string) (result ast.Node, residue []string) {
	if len(father) == 0 {
		return
	}
	if result = fileGetNode(f, father[0]); result == nil {
		return
	}

	i := 1
L:
	for result != nil && i < len(father) {
		f := father[i]
		switch n := result.(type) {
		case *ast.TypeSpec:
			result = n.Type
		case *ast.StructType:
			if tmp := structGetNode(n, f); tmp != nil {
				result = tmp
				i++
			} else {
				break L
			}
		case *ast.Field:
			if _, ok := n.Type.(*ast.Ident); ok {
				break L
			}
			result = n.Type
		// case *ast.FuncDecl:
		default:
			break L
		}
	}
	residue = father[i:]
	return
}

func fileGetNode(f *ast.File, key string) ast.Node {
	r, in := f.Scope.Objects[key]
	if in {
		if v, ok := r.Decl.(ast.Node); ok {
			return v
		}
	}
	return nil
}

func structGetNode(t *ast.StructType, key string) ast.Node {
	for _, field := range t.Fields.List {
		if key == field.Names[0].Name {
			return field
		}
	}
	return nil
}
