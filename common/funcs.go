package common

import (
	"go/ast"
	"reflect"
	"strings"
)

func HasComment(d ast.Decl, comment string) bool {
	if d == nil {
		return false
	}
	if g, ok := d.(*ast.GenDecl); !ok {
		return false
	} else {
		if g.Doc == nil {
			return false
		}
		for _, c := range g.Doc.List {
			if strings.Contains(c.Text, comment) {
				return true
			}
		}
	}
	return false
}

func GetTag(tag *ast.BasicLit, name, defaultValue string, pos int) string {
	if tag == nil {
		return defaultValue
	}
	keys := strings.FieldsFunc(reflect.StructTag(tag.Value[1:len(tag.Value)-1]).Get(name), func(r rune) bool {
		return r == ','
	})
	if len(keys)-1 < pos {
		return defaultValue
	}
	return keys[pos]
}

func FillDecls(f *ast.File, m map[string]*ast.StructType) {
	for _, node := range f.Decls {
		switch node.(type) {
		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)
					if s, ok := typeSpec.Type.(*ast.StructType); ok {
						m[typeSpec.Name.Name] = s
					}
				}
			}
		}
	}
}
