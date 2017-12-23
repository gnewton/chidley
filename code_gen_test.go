package main

import (
	//	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	//"log"
	"testing"
)

const hello = `package main

import (
      "encoding/xml"
)

type ChiAuthor struct {
	ChiName *ChiAuthor ` + "`" + `xml:\"http://www.w3.org/2005/Atom name,omitempty\" json:\"name,omitempty\""` + "`" + `
	XMLName  xml.Name ` + "`" + `xml:\"http://www.w3.org/2005/Atom author,omitempty\" json:\"author,omitempty\"` + "`" + `
}
`

func TestAverage(t *testing.T) {
	err := parseAndType(hello)
	if err != nil {
		t.Error(err)
	}
}

// Derived from example at https://github.com/golang/example/tree/master/gotypes#an-example
func parseAndType(s string) error {
	fset := token.NewFileSet()
	//f, err := parser.ParseFile(fset, "hello.go", hello, 0)
	f, err := parser.ParseFile(fset, "hello.go", hello, 0)
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
