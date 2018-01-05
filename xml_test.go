package main

import (
	//	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	//"log"
	"bytes"
	"strings"
	"testing"
)

func TestSameNameDifferentNameSpaceXML(t *testing.T) {
	err := extractor(sameNameDifferentNameSpaceXML)
	if err != nil {
		t.Error(err)
	}
}

func TestMixedCaseSameNameXML(t *testing.T) {
	err := extractor(mixedCaseSameNameXML)
	if err != nil {
		t.Error(err)
	}
}

func extractor(xml string) error {
	ex := Extractor{
		namePrefix:              namePrefix,
		nameSuffix:              nameSuffix,
		reader:                  strings.NewReader(xml),
		useType:                 useType,
		progress:                progress,
		ignoreXmlDecodingErrors: ignoreXmlDecodingErrors,
	}

	ex.init()
	err := ex.extract()

	if err != nil {
		return err
	}
	ex.done()

	buf := bytes.NewBufferString("")
	generateGoCode(buf, "foo", &ex)

	//log.Println(buf.String())

	err = parseAndType(buf.String())
	return err
}

// Derived from example at https://github.com/golang/example/tree/master/gotypes#an-example
func parseAndType(s string) error {
	fset := token.NewFileSet()
	//f, err := parser.ParseFile(fset, "hello.go", hello, 0)
	f, err := parser.ParseFile(fset, "hello.go", s, 0)
	if err != nil {
		//t.Error("Expected 1.5, got ", v)
		//log.Fatal(err) // parse error
		return err
	}

	conf := types.Config{Importer: importer.Default()}

	// Type-check the package containing only file f.
	// Check returns a *types.Package.
	_, err = conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
	if err != nil {
		//log.Fatal(err) // type error
		return err
	}
	return nil
}
