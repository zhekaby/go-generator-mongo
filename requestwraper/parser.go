package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/zhekaby/go-generator-mongo/common"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	structComment = "swagger:parameters"
	bodyComment   = "in: body"
)

type Parser struct {
	In, Dir, PkgPath, PkgName string
	isDir                     bool
	Structs                   []*structInfo
}

type structInfo struct {
	Name         string
	Body         *ast.StructType
	BodyTypeName string
	Fields       []*field
}
type visitor struct {
	*Parser

	name string
}

func NewParser(in string) *Parser {
	root, _ := os.Getwd()
	fin := path.Join(root, in)
	fInfo, err := os.Stat(fin)
	if err != nil {
		fmt.Printf("Error parsing %v: %v\n", in, err)
		os.Exit(1)
	}

	p := &Parser{
		In: fin, isDir: fInfo.IsDir(), Structs: make([]*structInfo, 0, 20),
	}

	if fInfo.IsDir() {
		p.Dir = fin
	} else {
		p.Dir = filepath.Dir(fin)
	}
	return p
}

func (p *Parser) Parse() error {
	os.Remove(fmt.Sprintf("%s/requestwrapper_validator_middleware.go", p.Dir))
	var err error
	if p.PkgPath, err = common.GetPkgPath(p.In, p.isDir); err != nil {
		return err
	}

	fset := token.NewFileSet()
	if p.isDir {
		log.Printf("parse dir: %s", p.Dir)
		packages, err := parser.ParseDir(fset, p.Dir, excludeTestFiles, parser.ParseComments)
		if err != nil {
			return err
		}

		if len(packages) != 1 {
			return fmt.Errorf("only one package in directory supported\n")
		}

		for _, pckg := range packages {
			//for _, f := range pckg.Files {
			//	//p.Decls = append(p.Decls, f.Decls...)
			//}
			ast.Walk(&visitor{Parser: p}, pckg)
		}
	} else {
		log.Printf("parse file: %s", p.In)
		f, err := parser.ParseFile(fset, p.In, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		//p.Decls = append(p.Decls, f.Decls...)
		ast.Walk(&visitor{Parser: p}, f)
	}

	spew.Dump(p)

	return nil
}

func (v *visitor) Visit(n ast.Node) (w ast.Visitor) {
	switch n := n.(type) {
	case *ast.Package:
		return v
	case *ast.File:
		v.PkgName = n.Name.String()
		return v

	case *ast.GenDecl:
		collectionName := v.needType(n.Doc)

		if collectionName != "" {
			for _, nc := range n.Specs {
				switch nct := nc.(type) {
				case *ast.TypeSpec:
					nct.Doc = n.Doc
				}
			}
		}

		return v
	case *ast.TypeSpec:
		collectionName := v.needType(n.Doc)
		if collectionName == "" {
			return nil
		}

		v.name = n.Name.String()

		// Allow to specify non-structs explicitly independent of '-all' flag.

		{
			if s, ok := n.Type.(*ast.StructType); ok {
				for _, f := range s.Fields.List {
					if needField(f, bodyComment) {
						st, ok := f.Type.(*ast.StarExpr).X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
						if !ok {
							st, ok = f.Type.(*ast.StructType)
							if !ok {
								return nil
							}
						}
						name := f.Type.(*ast.StarExpr).X.(*ast.Ident).Obj.Name
						v.Structs = append(v.Structs, &structInfo{
							Name:         n.Name.Name,
							Body:         st,
							BodyTypeName: name,
							Fields:       expandStruct(st, name),
						})
					}
				}
			}
			return nil
		}

		return v
	case *ast.StructType:
		//v.StructNames = append(v.StructNames, &structProps{Name: v.name})
		return nil
	}
	return nil
}

func (p *Parser) needType(comments *ast.CommentGroup) (collection string) {
	if comments == nil {
		return
	}

	for _, v := range comments.List {
		comment := v.Text

		if len(comment) > 2 {
			switch comment[1] {
			case '/':
				// -style comment (no newline at the end)
				comment = comment[2:]
			case '*':
				/*-style comment */
				comment = comment[2 : len(comment)-2]
			}
		}

		for _, comment := range strings.Split(comment, "\n") {

			comment = strings.TrimSpace(comment)

			if strings.HasPrefix(comment, structComment) {
				data := strings.FieldsFunc(comment, func(r rune) bool {
					return r == ' '
				})
				return data[1]
			}
		}
	}

	return
}

func needField(f *ast.Field, reqComment string) bool {
	if f.Doc == nil {
		return false
	}
	for _, c := range f.Doc.List {
		if strings.Contains(c.Text, reqComment) {
			return true
		}
	}
	return false
}

func excludeTestFiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

type field struct {
	Prop, Type, JsonProp, JsonPath, Ns, NsShort, NsCompact string
	Validations                                            map[string]string
}

func expandStruct(s *ast.StructType, ns string) []*field {
	return fnWalk("", ns, s.Fields.List)
}

func fnWalk(prefix, goPrefix string, astFields []*ast.Field) []*field {
	fields := make([]*field, 0, 50)
	if len(prefix) > 0 {
		prefix += "."
	}
	for _, f := range astFields {
		tag := getTag(f.Tag, "json", f.Names[0].Name, 0)
		bsonPath := prefix + tag
		ns := goPrefix + "." + f.Names[0].Name
		idx := strings.IndexByte(ns, byte('.'))
		switch n := f.Type.(type) {
		case *ast.Ident:

			fields = append(fields, &field{
				Prop:        f.Names[0].Name,
				Ns:          ns,
				NsShort:     ns[idx+1:],
				NsCompact:   strings.Replace(ns, ".", "", -1),
				JsonPath:    bsonPath,
				JsonProp:    tag,
				Type:        n.Name,
				Validations: getValidateRules(f.Tag),
			})
			if n.Obj == nil {
				continue
			}
			if t, ok := n.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType); ok {
				fields = append(fields, fnWalk(bsonPath, ns, t.Fields.List)...)
			}
		case *ast.StructType:
			fields = append(fields, fnWalk(bsonPath, ns, n.Fields.List)...)
		case *ast.StarExpr:
			if t, ok := n.X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType); ok {
				fields = append(fields, fnWalk(bsonPath, ns, t.Fields.List)...)
			}

		default:
			spew.Dump(1)
		}
	}
	return fields
}

func getTag(tag *ast.BasicLit, name, defaultValue string, pos int) string {
	if tag == nil {
		return defaultValue
	}
	keys := strings.Split(reflect.StructTag(tag.Value[1:len(tag.Value)-1]).Get(name), ",")
	if len(keys)-1 < pos {
		return defaultValue
	}
	return keys[pos]
}

func getValidateRules(tag *ast.BasicLit) map[string]string {
	if tag == nil {
		return nil
	}
	m := make(map[string]string, 5)
	keys := strings.Split(reflect.StructTag(tag.Value[1:len(tag.Value)-1]).Get("validate"), ",")
	for _, k := range keys {
		if k == "" || k == "required" {
			continue
		}
		idx := strings.IndexByte(k, '=')
		if idx > 0 {
			m[k[:idx]] = k[idx+1:]
		} else {
			m[k] = ""
		}
	}
	return m
}
