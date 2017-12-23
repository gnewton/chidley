package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type OutVariableDef struct {
	XMLName      string
	XMLNameSpace string
	GoName       string
	GoType       string
	Length       int
	//Foo          bool
}

var fieldTemplate *template.Template

// const t = `
// type Article struct{
//        {{.GoName}} {{.GoType}} ` + "`" + `xml:"{{if notEmpty .XMLNameSpace}}{{.XMLNameSpace}} {{end}}{{.XMLName}},omitempty" json:"{{.XMLName}},omitempty"{{if .Foo}} db:"size={{.Length}}" {{- else}} gorm:"name:name,size={{.Length}}"{{- end}}` + "`" + `
// }
// `

const fieldTemplateString = `
       {{.GoName}} {{.GoType}} ` + "`" + `xml:"{{if notEmpty .XMLNameSpace}}{{.XMLNameSpace}} {{end}}{{.XMLName}},omitempty" json:"{{.XMLName}},omitempty"` + "`" + `
`

func render() {
	fieldTemplate = template.Must(template.New("fieldTemplate").Funcs(template.FuncMap{
		"notEmpty": func(feature string) bool {
			return len(feature) > 0
		},
	}).Parse(fieldTemplateString))

	ot := OutVariableDef{
		"author",
		"", //"http://w3/org/mmm",
		"Author",
		"string",
		16,
		//	false
	}

	var buf bytes.Buffer

	//err := t.Execute(os.Stdout, ot)
	err := fieldTemplate.Execute(&buf, ot)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

}
