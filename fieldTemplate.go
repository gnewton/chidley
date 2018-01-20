package main

import (
	"bytes"
	//"fmt"
	"strconv"
	"text/template"
)

type OutVariableDef struct {
	XMLName              string
	XMLNameSpace         string
	GoName               string
	GoType               string
	GoTypeArrayOrPointer string
	Length               int64
	//Foo          bool
}

var fieldTemplate *template.Template

// const t = `
// type Article struct{
//        {{.GoName}} {{.GoType}} ` + "`" + `xml:"{{if notEmpty .XMLNameSpace}}{{.XMLNameSpace}} {{end}}{{.XMLName}},omitempty" json:"{{.XMLName}},omitempty"{{if .Foo}} db:"size={{.Length}}" {{- else}} gorm:"name:name,size={{.Length}}"{{- end}}` + "`" + `
// }
// `

var fieldTemplateString = `{{.GoName}} {{.GoTypeArrayOrPointer}}{{.GoType}} ` + "`" + `xml:"{{if notEmpty .XMLNameSpace}}{{.XMLNameSpace}} {{end}}{{.XMLName}},omitempty" json:"{{.XMLName}},omitempty"` + "`" + ``

func render(otd OutVariableDef) string {
	fieldTemplate = template.Must(template.New("fieldTemplate").Funcs(template.FuncMap{
		"notEmpty": func(feature string) bool {
			return len(feature) > 0
		},
	}).Parse(fieldTemplateString))

	var buf bytes.Buffer

	//err := t.Execute(os.Stdout, ot)
	err := fieldTemplate.Execute(&buf, otd)
	if err != nil {
		panic(err)
	}
	//fmt.Println(buf.String())
	return "\t" + buf.String() + "   // ZZmaxLength=" + strconv.FormatInt(otd.Length, 10)

}
