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
	"strings"
)

type Parser struct {
	In, Dir, PkgPath, PkgName string
	isDir                     bool
	Collections               []*collection
	Structs                   []*structInfo
	Decls                     map[string]*ast.StructType
}

type structInfo struct {
	Name         string
	Body         *ast.StructType
	BodyTypeName string
	Fields       []field
}

type collection struct {
	Typ, Name string
	Fields    []field
}

type field struct {
	Prop, Type, BsonProp, BsonPath, GoPath string
}

type visitor struct {
	*Parser

	name string
}

var (
	structCommentMongo   = "mongowrapper:collection"
	structCommentSwagger = "swagger:parameters"
	bodyComment          = "in: body"
)

func NewParser(in string) *Parser {
	root, _ := os.Getwd()
	fin := path.Join(root, in)
	fInfo, err := os.Stat(fin)
	if err != nil {
		fmt.Printf("Error parsing %v: %v\n", in, err)
		os.Exit(1)
	}

	p := &Parser{
		In: fin, isDir: fInfo.IsDir(),
		Structs:     make([]*structInfo, 0, 50),
		Collections: make([]*collection, 0, 50),
		Decls:       make(map[string]*ast.StructType, 200),
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
			for _, f := range pckg.Files {
				common.FillDecls(f, p.Decls)
			}
			ast.Walk(&visitor{Parser: p}, pckg)
		}
	} else {
		log.Printf("parse file: %s", p.In)
		f, err := parser.ParseFile(fset, p.In, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		common.FillDecls(f, p.Decls)
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
		collectionName := v.needType(n.Doc, structCommentMongo)
		if collectionName == "" {
			collectionName = v.needType(n.Doc, structCommentSwagger)
		}

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
		collectionName := v.needType(n.Doc, structCommentMongo)
		if collectionName != "" {

			v.name = n.Name.String()
			fmt.Printf("parsing %s\n", v.name)

			fields := make([]field, 0, 100)
			deep(n.Type, "", "", "", "", &fields)
			v.Parser.Collections = append(v.Parser.Collections, &collection{
				Typ:    v.name,
				Name:   collectionName,
				Fields: fields,
			})
		}

		collectionName = v.needType(n.Doc, structCommentSwagger)
		if collectionName != "" {
			if s, ok := n.Type.(*ast.StructType); ok {
				for _, f := range s.Fields.List {
					if needField(f, bodyComment) {
						ident, ok := f.Type.(*ast.StarExpr).X.(*ast.Ident)
						if !ok {
							return nil
						}
						var st *ast.StructType
						if ident.Obj == nil {
							st = v.Decls[ident.Name]
						} else {
							st, ok = ident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
							if !ok {
								st, ok = f.Type.(*ast.StructType)
								if !ok {
									return nil
								}
							}
						}
						name := ident.Name
						fields := make([]field, 0, 100)
						deep(st, "", "", "", "", &fields)
						v.Structs = append(v.Structs, &structInfo{
							Name:         n.Name.Name,
							Body:         st,
							BodyTypeName: name,
							Fields:       fields,
						})
					}
				}
			}
		}

		return nil
	case *ast.StructType:
		//v.StructNames = append(v.StructNames, &structProps{Name: v.name})
		return nil
	}
	return nil
}

func (p *Parser) needType(comments *ast.CommentGroup, reqComment string) (collection string) {
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

			if strings.HasPrefix(comment, reqComment) {
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

func deep(n ast.Node, fieldName, tag, bsonPrefix, goPrefix string, fields *[]field) {
	switch n := n.(type) {
	case *ast.TypeSpec:
		switch ns := n.Type.(type) {
		case *ast.StarExpr:
			deep(ns.X, "", "", bsonPrefix, goPrefix, fields)
		case *ast.StructType:
			deep(ns, "", "", bsonPrefix, goPrefix, fields)
		default:
			return
		}
	case *ast.GenDecl:
		for _, nc := range n.Specs {
			switch nct := nc.(type) {
			case *ast.TypeSpec:
				deep(nc, nct.Name.Name, tag, bsonPrefix, goPrefix, fields)

			}
		}
	case *ast.StructType:
		if len(bsonPrefix) > 0 {
			bsonPrefix += "."
		}
		for _, f := range n.Fields.List {
			tag := common.GetTag(f.Tag, "bson", f.Names[0].Name, 0)
			switch ss := f.Type.(type) {
			case *ast.StructType:
				deep(ss, f.Names[0].Name, tag, bsonPrefix+tag, goPrefix+f.Names[0].Name, fields)
			case *ast.StarExpr:
				if ident, ok := ss.X.(*ast.Ident); ok {
					if ident.Obj != nil {
						if ts, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
							deep(ts.Type, f.Names[0].Name, tag, bsonPrefix+tag, goPrefix+f.Names[0].Name, fields)
						}
					} else {
						deep(f.Type, f.Names[0].Name, tag, bsonPrefix, goPrefix, fields)
					}
				}
			default:
				deep(f.Type, f.Names[0].Name, tag, bsonPrefix, goPrefix, fields)
			}

		}
	case *ast.Ident:
		var typ string
		if n.Obj != nil {
			typ = n.Obj.Name
		} else {
			typ = n.Name
		}
		f := &field{
			Prop:     fieldName,
			GoPath:   goPrefix + fieldName,
			BsonPath: bsonPrefix + tag,
			BsonProp: tag,
			Type:     typ,
		}
		*fields = append(*fields, *f)
	case *ast.SelectorExpr:
		var typ string
		if e, ok := n.X.(*ast.Ident); ok {
			if e.Name != "" {
				typ = e.Name + "."
			}
		}
		f := &field{
			Prop:     fieldName,
			GoPath:   goPrefix + fieldName,
			BsonPath: bsonPrefix + tag,
			BsonProp: tag,
			Type:     typ + n.Sel.Name,
		}
		*fields = append(*fields, *f)
		break
	case *ast.StarExpr:
		deep(n.X, fieldName, tag, bsonPrefix, goPrefix, fields)
		break
	default:
	}
}
