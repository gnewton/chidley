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

func TestTagsContainHyphens(t *testing.T) {
	err := extractor(tagsContainHyphens)
	if err != nil {
		t.Error(err)
	}
}

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

//https://github.com/gnewton/chidley/issues/14
func TestGithubIssue14(t *testing.T) {
	err := extractor(githubIssue14)
	if err != nil {
		t.Error(err)
	}
}

func extractor(xmlString string) error {
	ex := Extractor{
		namePrefix:              namePrefix,
		nameSuffix:              nameSuffix,
		useType:                 useType,
		progress:                progress,
		ignoreXmlDecodingErrors: ignoreXmlDecodingErrors,
	}

	ex.init()
	err := ex.extract(strings.NewReader(xmlString))

	if err != nil {
		return err
	}

	err = ex.extract(strings.NewReader(xmlString))
	if err != nil {
		return err
	}

	ex.done()

	buf := bytes.NewBufferString("")
	fps := make([]string, 1)
	fps[0] = "foo"
	generateGoCode(buf, fps, &ex)

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
