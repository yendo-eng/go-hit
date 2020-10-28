package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Eun/go-hit/generators/helpers"

	"github.com/dave/jennifer/jen"

	"github.com/Eun/yaegi-template/codebuffer"
)

func main() {
	if err := doit(); err != nil {
		panic(err)
	}
}

func doit() error {
	fl, err := os.Open("README.md")
	if err != nil {
		return err
	}
	defer fl.Close()
	cb := codebuffer.New(fl, []rune("```go"), []rune("```"))
	it, err := cb.Iterator()
	if err != nil {
		return err
	}

	f := jen.NewFile("hit_test")
	f.HeaderComment("+build doctest")
	f.PackageComment("⚠️⚠️⚠️ This file was autogenerated by generators/readme/readme ⚠️⚠️⚠️ //")
	f.Op(`import . "github.com/Eun/go-hit"`)
	f.Op(`import "net/http"`)

	i := 0
	for it.Next() {
		v := it.Value()
		if v.Type == codebuffer.CodePartType {
			if bytes.HasPrefix(v.Content, []byte(" //ignore")) {
				continue
			}
			f.Func().Id(fmt.Sprintf("TestReadmeCodePart%d", i)).Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(
				jen.Op(strings.TrimSpace(string(v.Content))),
			)
			i++
		}
	}
	if err := it.Error(); err != nil {
		return err
	}

	if err := helpers.WriteJenFile("readme_gen_test.go", f); err != nil {
		return err
	}

	return nil
}
