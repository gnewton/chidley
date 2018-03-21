package main

import (
	"bytes"
	//"fmt"
	"log"
	//"strconv"
	"text/template"
)

type FieldDef struct {
	XMLName              string
	XMLNameSpace         string
	GoName               string
	GoType               string
	GoTypeArrayOrPointer string
	Length               int64
	//Foo          bool
}

var fieldTemplate *template.Template

var fieldTemplateString = `{{.GoName}} {{.GoTypeArrayOrPointer}}{{.GoType}} ` + "`" + `xml:"{{if notEmpty .XMLNameSpace}}{{.XMLNameSpace}} {{end}}{{.XMLName}},omitempty" json:"{{.XMLName}},omitempty"` + "`" + ``

func render(otd FieldDef) (string, error) {
	var err error
	//	fieldTemplate = template.Must(template.New("fieldTemplate").Funcs(template.FuncMap{
	fieldTemplate, err = template.New("fieldTemplate").Funcs(template.FuncMap{
		"notEmpty": func(feature string) bool {
			return len(feature) > 0
		},
	}).Parse(fieldTemplateString)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	//err := t.Execute(os.Stdout, ot)
	err = fieldTemplate.Execute(&buf, otd)
	if err != nil {
		return "", err
	}
	//fmt.Println(buf.String())
	//return "\t" + buf.String() + "   // ZZmaxLength=" + strconv.FormatInt(otd.Length, 10), nil
	return "\t" + buf.String(), nil
}

func runValidateFieldTemplate(printToLog bool) error {
	fieldDef := FieldDef{
		XMLName:              "xmlName",
		XMLNameSpace:         "xmlNameSpace",
		GoName:               "Foobar",
		GoType:               "string",
		GoTypeArrayOrPointer: "[]",
		Length:               32,
	}

	string, err := render(fieldDef)
	if printToLog {
		log.Println("validateFieldTemplate")
		log.Println("Using template:", fieldTemplateString)
		log.Println(string)

	}
	if err != nil {
		log.Println("Error with field template:", err)
	}
	return err
}
