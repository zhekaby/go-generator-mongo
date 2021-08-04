package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
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

type Parser struct {
	Cs, CsVar, DbVar, In, Dir, PkgPath, PkgName string
	isDir                                       bool
	collections                                 []*collection
	//Decls []ast.Decl
}

type collection struct {
	Typ, Name string
	Fields    []*field
}

type field struct {
	Prop, Type, BsonProp, BsonPath, GoPath string
}

type visitor struct {
	*Parser

	name string
}

var (
	structComment = "mongowarapper:collection"
)

func NewParser(cs, csVar, dbVar, in string) *Parser {
	root, _ := os.Getwd()
	fin := path.Join(root, in)
	fInfo, err := os.Stat(fin)
	if err != nil {
		fmt.Printf("Error parsing %v: %v\n", in, err)
		os.Exit(1)
	}

	p := &Parser{
		Cs: cs, CsVar: csVar, DbVar: dbVar, In: fin, isDir: fInfo.IsDir(),
	}

	if fInfo.IsDir() {
		p.Dir = fin
	} else {
		p.Dir = filepath.Dir(fin)
	}
	return p
}

func (p *Parser) Parse() error {
	var err error
	if p.PkgPath, err = getPkgPath(p.In, p.isDir); err != nil {
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
				c := &collection{Typ: v.name, Name: collectionName, Fields: make([]*field, 0, len(s.Fields.List))}
				v.Parser.collections = append(v.Parser.collections, c)

				c.Fields = expandStruct(s)
			}

			//keys := strings.Split(reflect.StructTag(f.Tag.Value[1:len(f.Tag.Value)-1]).Get("json"), ",")
			//res := &structProps{Name: v.name, QueryFields: make([]*QueryFieldProps, 0, 20)}
			//if s, ok := n.Type.(*ast.StructType); ok {
			//	for _, f := range s.Fields.List {
			//		if f.Names[0].Name == "Body" {
			//			star, ok := f.Type.(*ast.StarExpr)
			//			if !ok {
			//				panic("only star expression supported for Body. Use: Body *type")
			//			}
			//			res.BodyType = star.X.(*ast.Ident).Name
			//		} else if f.Names[0].Name == "Query" {
			//			star, ok := f.Type.(*ast.StarExpr)
			//			if !ok {
			//				panic("only star expression supported for Query. Use: Query *type")
			//			}
			//			res.QueryType = star.X.(*ast.Ident).Name
			//			ss, _ := v.findStruct(res.QueryType)
			//			for _, f := range ss.Fields.List {
			//				p := &QueryFieldProps{
			//					Type:  f.Type.(*ast.Ident).String(),
			//					Field: f.Names[0].Name,
			//				}
			//				if f.Tag.Value == "" {
			//					p.Name = f.Names[0].Name
			//				} else {
			//					keys := strings.Split(reflect.StructTag(f.Tag.Value[1:len(f.Tag.Value)-1]).Get("json"), ",")
			//					p.Name = keys[0]
			//				}
			//				res.QueryFields = append(res.QueryFields, p)
			//			}
			//		}
			//	}

			//}
			//v.StructNames = append(v.StructNames, res)
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

func excludeTestFiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

func expandStruct(s *ast.StructType) []*field {
	return fnWalk("", "", s.Fields.List)
}

func fnWalk(prefix, goPrefix string, astFields []*ast.Field) []*field {
	fields := make([]*field, 0, 50)
	if len(prefix) > 0 {
		prefix += "."
	}
	for _, f := range astFields {
		tag := getTag(f.Tag, "bson", f.Names[0].Name, 0)
		bsonPath := prefix + tag
		goPath := goPrefix + f.Names[0].Name

		switch n := f.Type.(type) {
		case *ast.Ident:
			fields = append(fields, &field{
				Prop:     f.Names[0].Name,
				GoPath:   goPath,
				BsonPath: bsonPath,
				BsonProp: tag,
				Type:     n.Name,
			})
			if n.Obj == nil {
				continue
			}
			if t, ok := n.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType); ok {
				fields = append(fields, fnWalk(bsonPath, goPath, t.Fields.List)...)
			}
		case *ast.StructType:
			fields = append(fields, fnWalk(bsonPath, goPath, n.Fields.List)...)
		case *ast.StarExpr:
			if t, ok := n.X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType); ok {
				fields = append(fields, fnWalk(bsonPath, goPath, t.Fields.List)...)
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
