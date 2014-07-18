package main

import (
	"os"
)

type Verify struct {
	writer      Writer
	lineChannel chan string
	xmlFilename string
	root        *Node
	namePrefix  string
	nameSuffix  string
}

const GOFILENAME = "VerifyStructs.go"

func (v *Verify) init(root *Node, namePrefix string, nameSuffix string, dir string, xmlFilename string, lineChannel chan string) (Writer, error) {
	v.root = root
	v.namePrefix = namePrefix
	v.nameSuffix = nameSuffix
	v.lineChannel = lineChannel
	v.xmlFilename = xmlFilename
	var err error
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
	}

	v.writer = new(fileWriter)
	v.writer.open(dir+"/"+GOFILENAME, lineChannel)

	v.generateCodePre()

	return v.writer, err

}

func (v *Verify) generateCodePre() {
	v.lineChannel <- "package main"
	v.lineChannel <- " "

	v.lineChannel <- "import ("
	v.lineChannel <- "	\"encoding/xml\""
	v.lineChannel <- "	\"log\""
	v.lineChannel <- "	\"os\""
	v.lineChannel <- "	\"fmt\""
	v.lineChannel <- ")"
}

func (v *Verify) makeName(name string) string {
	return v.namePrefix + capitalizeFirstLetter(name) + v.nameSuffix
}

func (v *Verify) generatePost() {
	v.lineChannel <- "func main(){"
	v.lineChannel <- " "
	v.lineChannel <- "xmlFile, err := os.Open(\"../" + v.xmlFilename + "\")"

	v.lineChannel <- "if err != nil {"
	v.lineChannel <- "log.Fatal(err)"
	v.lineChannel <- "}"
	v.lineChannel <- "defer xmlFile.Close()"
	v.lineChannel <- "decoder := xml.NewDecoder(xmlFile)"

	v.lineChannel <- "	for {"
	v.lineChannel <- "		token, _ := decoder.Token()"
	v.lineChannel <- "		if token == nil {"
	v.lineChannel <- "			break"
	v.lineChannel <- "		}"
	v.lineChannel <- "		switch se := token.(type) {"
	v.lineChannel <- "		case xml.StartElement:"
	//for name, node := range v.root.children {
	for _, node := range v.root.children {
		for name, _ := range node.children {
			v.lineChannel <- "			if se.Name.Local == \"" + name + "\" {"
			v.lineChannel <- "				var item " + v.makeName(name)
			v.lineChannel <- "				decoder.DecodeElement(&item, &se)"
			v.lineChannel <- "				fmt.Printf(\"%+v\\n\", item)"
			v.lineChannel <- ""
			v.lineChannel <- "			}"
		}
	}
	v.lineChannel <- "		}"

	v.lineChannel <- "	}"
	v.lineChannel <- "}"

}
