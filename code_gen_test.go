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
	ex := Extractor{
		namePrefix:              namePrefix,
		nameSuffix:              nameSuffix,
		reader:                  strings.NewReader(modx),
		useType:                 useType,
		progress:                progress,
		ignoreXmlDecodingErrors: ignoreXmlDecodingErrors,
	}

	ex.init()
	err := ex.extract()

	if err != nil {
		t.Error(err)
	}
	ex.done()

	buf := bytes.NewBufferString("")
	generateGoCode(buf, "foo", &ex)

	//log.Println(buf.String())

	err = parseAndType(buf.String())
	if err != nil {
		t.Error(err)
	}
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
