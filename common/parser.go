package common

import (
	"fmt"
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
	structCommentMongo   = "mongowrapper:collection"
	structCommentSwagger = "swagger:parameters"
	bodyComment          = "in: body"
)

type Parser struct {
	In, Dir, PkgPath, PkgName string
	isDir                     bool
	Collections               []*Collection
	Structs                   []*StructInfo
	Decls                     map[string]*ast.StructType
}

type StructInfo struct {
	Name         string
	Body         *ast.StructType
	BodyTypeName string
	Fields       []Field
}
type Field struct {
	Prop, Type, JsonProp, JsonPath, BsonProp, BsonPath, GoPath, Ns, NsShort, NsCompact string
	Validations                                                                        map[string]string
}
type Collection struct {
	Typ, Name string
	Fields    []Field
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
		In: fin, isDir: fInfo.IsDir(),
		Structs:     make([]*StructInfo, 0, 50),
		Collections: make([]*Collection, 0, 50),
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
	if p.PkgPath, err = GetPkgPath(p.In, p.isDir); err != nil {
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
				FillDecls(f, p.Decls)
			}
			ast.Walk(&visitor{Parser: p}, pckg)
		}
	} else {
		log.Printf("parse file: %s", p.In)
		f, err := parser.ParseFile(fset, p.In, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		FillDecls(f, p.Decls)
		ast.Walk(&visitor{Parser: p}, f)
	}

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
		args := v.needType(n.Doc, structCommentMongo)
		if len(args) == 0 {
			args = v.needType(n.Doc, structCommentSwagger)
		}

		if len(args) > 0 {
			for _, nc := range n.Specs {
				switch nct := nc.(type) {
				case *ast.TypeSpec:
					nct.Doc = n.Doc
				}
			}
		}

		return v
	case *ast.TypeSpec:
		args := v.needType(n.Doc, structCommentMongo)
		if len(args) > 0 {

			v.name = n.Name.String()
			fmt.Printf("parsing %s\n", v.name)

			fields := make([]Field, 0, 100)
			deep(n.Type, "", "", "", "", "", "", "", "", &fields)
			v.Parser.Collections = append(v.Parser.Collections, &Collection{
				Typ:    v.name,
				Name:   args[1],
				Fields: fields,
			})
		}

		args = v.needType(n.Doc, structCommentSwagger)
		if len(args) > 0 {
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
						fields := make([]Field, 0, 100)
						deep(st, "", "", "", "", "", "", name, "", &fields)
						v.Structs = append(v.Structs, &StructInfo{
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

func (p *Parser) needType(comments *ast.CommentGroup, reqComment string) (arguments []string) {
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
				return data
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

func deep(n ast.Node, fieldName, jsonTag, jsonPrefix, bsonTag, bsonPrefix, goPrefix, ns, tag string, fields *[]Field) {
	switch n := n.(type) {
	case *ast.TypeSpec:
		switch ts := n.Type.(type) {
		case *ast.StarExpr:
			deep(ts.X, "", "", "", "", bsonPrefix, goPrefix, ns, "", fields)
		case *ast.StructType:
			deep(ts, "", "", "", "", bsonPrefix, goPrefix, ns, "", fields)
		default:
			return
		}
	case *ast.GenDecl:
		for _, nc := range n.Specs {
			switch nct := nc.(type) {
			case *ast.TypeSpec:
				deep(nc, nct.Name.Name, jsonTag, jsonPrefix, bsonTag, bsonPrefix, goPrefix, ns, "", fields)

			}
		}
	case *ast.StructType:
		if len(bsonPrefix) > 0 {
			bsonPrefix += "."
		}
		if len(jsonPrefix) > 0 {
			jsonPrefix += "."
		}
		for _, f := range n.Fields.List {
			bsonTag := GetTag(f.Tag, "bson", f.Names[0].Name, 0)
			jsonTag := GetTag(f.Tag, "json", f.Names[0].Name, 0)
			var tag = ""
			if f.Tag != nil {
				tag = f.Tag.Value
			}
			switch ss := f.Type.(type) {
			case *ast.StructType:
				deep(ss, f.Names[0].Name, jsonTag, jsonPrefix+jsonTag, bsonTag, bsonPrefix+bsonTag, goPrefix+f.Names[0].Name, ns, tag, fields)
			case *ast.StarExpr:
				if ident, ok := ss.X.(*ast.Ident); ok {
					if ident.Obj != nil {
						if ts, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
							deep(ts.Type, f.Names[0].Name, jsonTag, jsonPrefix+jsonTag, bsonTag, bsonPrefix+bsonTag, goPrefix+f.Names[0].Name, ns+"."+f.Names[0].Name, tag, fields)
						}
					} else {
						deep(f.Type, f.Names[0].Name, jsonTag, jsonPrefix, bsonTag, bsonPrefix, goPrefix, ns, tag, fields)
					}
				}
			default:
				deep(f.Type, f.Names[0].Name, jsonTag, jsonPrefix, bsonTag, bsonPrefix, goPrefix, ns, tag, fields)
			}

		}
	case *ast.Ident:
		var typ string
		if n.Obj != nil {
			typ = n.Obj.Name
		} else {
			typ = n.Name
		}
		ns += "." + fieldName
		idx := strings.IndexByte(ns, byte('.'))
		f := &Field{
			Prop:        fieldName,
			GoPath:      goPrefix + fieldName,
			JsonProp:    jsonTag,
			JsonPath:    jsonPrefix + jsonTag,
			BsonProp:    bsonTag,
			BsonPath:    bsonPrefix + bsonTag,
			Type:        typ,
			Ns:          ns,
			NsShort:     ns[idx+1:],
			NsCompact:   strings.Replace(ns, ".", "", -1),
			Validations: getValidateRules(tag),
		}
		*fields = append(*fields, *f)
	case *ast.SelectorExpr:
		var typ string
		if e, ok := n.X.(*ast.Ident); ok {
			if e.Name != "" {
				typ = e.Name + "."
			}
		}
		ns := goPrefix + "." + fieldName
		idx := strings.IndexByte(ns, byte('.'))
		f := &Field{
			Prop:        fieldName,
			GoPath:      goPrefix + fieldName,
			JsonProp:    jsonTag,
			JsonPath:    jsonPrefix + jsonTag,
			BsonProp:    bsonTag,
			BsonPath:    bsonPrefix + bsonTag,
			Type:        typ + n.Sel.Name,
			Ns:          ns,
			NsShort:     ns[idx+1:],
			NsCompact:   strings.Replace(ns, ".", "", -1),
			Validations: getValidateRules(tag),
		}
		*fields = append(*fields, *f)
		break
	case *ast.StarExpr:
		deep(n.X, fieldName, jsonTag, jsonPrefix, bsonTag, bsonPrefix, goPrefix, ns, tag, fields)
		break
	default:
	}
}

func getValidateRules(tag string) map[string]string {
	if tag == "" {
		return nil
	}
	m := make(map[string]string, 5)
	keys := strings.Split(reflect.StructTag(tag[1:len(tag)-1]).Get("validate"), ",")
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