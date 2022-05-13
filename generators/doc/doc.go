package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"

	"github.com/otto-eng/go-hit/generators/helpers"

	"path/filepath"

	"fmt"

	"bufio"
	"bytes"

	"github.com/dave/jennifer/jen"

	"github.com/Eun/go-testdoc"
)

func generateDocForFile(tmpFile *jen.File, fileName string, pos int, name, doc string) {
	doc = strings.TrimSpace(doc)
	if doc == "" {
		return
	}

	d, err := testdoc.ParseDoc(doc)
	if err != nil {
		panic(err)
	}

	if s, ok := d.Fields["Usage"]; ok && s != "" {
		tmpFile.Add(generateTestBlock(fileName, pos, "Usage", false, name, s))
	}

	if s, ok := d.Fields["Example"]; ok && s != "" {
		tmpFile.Add(generateTestBlock(fileName, pos, "Example", true, name, s))
	}

	if s, ok := d.Fields["Examples"]; ok && s != "" {
		tmpFile.Add(generateTestBlock(fileName, pos, "Examples", true, name, s))
	}
}

func generateTestBlock(fileName string, pos int, funcPrefix string, expectRequest bool, name, code string) jen.Code {
	var codeLines []jen.Code

	i := 1

	scanner := bufio.NewScanner(bytes.NewReader([]byte(code)))
	var sb strings.Builder

	codeExpectRequest := jen.False()
	if expectRequest {
		codeExpectRequest = jen.True()
	}

	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			codeLines = append(codeLines, jen.Add(
				jen.Commentf("%s:%d", fileName, pos),
				jen.Line(),
				jen.Func().Id(fmt.Sprintf("Test%s%s%d", funcPrefix, strings.Title(name), i)).Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
					jen.Qual("github.com/otto-eng/go-hit/doctest", "RunTest").Call(
						codeExpectRequest,
						jen.Func().Params().Block(jen.Op(sb.String())),
					),
				),
				jen.Line(),
			))
			sb.Reset()
			i++
			continue
		}
		sb.WriteString(txt)
		sb.WriteRune('\n')
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if sb.Len() > 0 {
		codeLines = append(codeLines, jen.Add(
			jen.Commentf("%s:%d", fileName, pos),
			jen.Line(),
			jen.Func().Id(fmt.Sprintf("Test%s%s%d", funcPrefix, strings.Title(name), i)).Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Qual("github.com/otto-eng/go-hit/doctest", "RunTest").Call(
					codeExpectRequest,
					jen.Func().Params().Block(jen.Op(sb.String())),
				),
			),
			jen.Line(),
		))
	}

	return jen.Add(codeLines...)
}

func buildName(typ *ast.Ident, names []*ast.Ident) string {
	var sb strings.Builder
	fset := token.NewFileSet()
	if err := printer.Fprint(&sb, fset, typ); err != nil {
		panic(err)
	}

	for _, name := range names {
		fset := token.NewFileSet()
		if err := printer.Fprint(&sb, fset, name); err != nil {
			panic(err)
		}
	}

	return sb.String()
}

func main() {
	fmt.Println("collecting functions...")
	set := token.NewFileSet()
	pkgs, err := parser.ParseDir(set, ".", nil, parser.ParseComments)
	if err != nil {
		log.Fatal("Failed to parse package:", err)
	}

	pkg := pkgs["hit"]

	tmpFile := jen.NewFile("hit_test")
	tmpFile.HeaderComment("+build doctest")
	tmpFile.Op(`import . "github.com/otto-eng/go-hit"`)
	tmpFile.Op(`import "net/http"`)
	tmpFile.Comment("⚠️⚠️⚠️ This file was autogenerated by generators/doc ⚠️⚠️⚠️ //")

	for fileName, f := range pkg.Files {
		cleanFileName := fileName
		if ext := filepath.Ext(cleanFileName); ext != "" {
			cleanFileName = cleanFileName[:len(cleanFileName)-len(ext)]
		}

		type entry struct {
			text   string
			lineno int
		}

		docMap := make(map[string]entry)
		ast.Inspect(f, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.FuncDecl:
				if v.Doc == nil || strings.TrimSpace(v.Doc.Text()) == "" {
					return true
				}

				docMap[cleanFileName+v.Name.String()] = entry{
					text:   v.Doc.Text(),
					lineno: set.Position(v.Doc.Pos()).Line,
				}

			// find variable declarations
			case *ast.TypeSpec:
				// which are public
				if v.Name.IsExported() {
					if i, ok := v.Type.(*ast.InterfaceType); ok {
						// and are interfaces
						for _, m := range i.Methods.List {
							if m.Doc == nil || strings.TrimSpace(m.Doc.Text()) == "" {
								continue
							}
							docMap[cleanFileName+buildName(v.Name, m.Names)] = entry{
								text:   m.Doc.Text(),
								lineno: set.Position(m.Doc.Pos()).Line,
							}
						}
					}
				}
			}
			return true
		})
		fmt.Printf("generating for %s...\n", fileName)
		for name, doc := range docMap {
			generateDocForFile(tmpFile, fileName, doc.lineno, name, doc.text)
		}
	}
	if err := helpers.WriteJenFile("doc_gen_test.go", tmpFile); err != nil {
		log.Fatal(err)
	}
}
