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

func (v *Verify) init(root *Node, namePrefix string, nameSuffix string, verifyDir string, verifyGoFile string, xmlFilename string, lineChannel chan string) (Writer, error) {
	v.root = root
	v.namePrefix = namePrefix
	v.nameSuffix = nameSuffix
	v.lineChannel = lineChannel
	v.xmlFilename = xmlFilename
	var err error
	err = os.RemoveAll(verifyDir)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(verifyDir); os.IsNotExist(err) {
		err = os.Mkdir(verifyDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	v.writer = new(fileWriter)
	v.writer.open(verifyDir+"/"+verifyGoFile, lineChannel)

	return v.writer, err

}

func (v *Verify) generateCodePre(hasStartElements bool) {
	v.lineChannel <- "package main"
	v.lineChannel <- " "

	v.lineChannel <- "import ("
	if !hasStartElements {
		v.lineChannel <- "	\"fmt\""
		v.lineChannel <- "	\"os\""
		v.lineChannel <- ")"
		return
	}
	v.lineChannel <- "	\"encoding/xml\""
	v.lineChannel <- "	\"encoding/json\""
	v.lineChannel <- "	\"os\""
	v.lineChannel <- "	\"log\""
	v.lineChannel <- "	\"fmt\""
	v.lineChannel <- "      \"bufio\""
	v.lineChannel <- "      \"compress/bzip2\""
	v.lineChannel <- "      \"compress/gzip\""
	v.lineChannel <- "      \"io\""
	v.lineChannel <- "      \"strings\""
	v.lineChannel <- ")"
}

func (v *Verify) makeName(name string) string {
	//return v.namePrefix + capitalizeFirstLetter(name) + v.nameSuffix
	return v.namePrefix + cleanName(name) + v.nameSuffix
}

func (v *Verify) generatePost(hasStartElements bool, globalTagAttributes map[string](map[string]string)) {
	v.lineChannel <- "func main(){"
	v.lineChannel <- " "
	if !hasStartElements {
		v.lineChannel <- "fmt.Println(\"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\")"
		v.lineChannel <- "fmt.Println(\"++++++ No StartElements encountered: possible error in XML or bug+++++++\")"
		v.lineChannel <- "fmt.Println(\"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\")"
		v.lineChannel <- "os.Exit(42)"
		v.lineChannel <- "}"
		return

	}
	v.lineChannel <- "reader, xmlFile, err := genericReader(\"../" + v.xmlFilename + "\")"

	v.lineChannel <- "if err != nil {"
	v.lineChannel <- "log.Fatal(err)"
	v.lineChannel <- "return"
	v.lineChannel <- "}"
	v.lineChannel <- "defer xmlFile.Close()"
	v.lineChannel <- "decoder := xml.NewDecoder(reader)"

	v.lineChannel <- "	for i:=0; i<100; i++ {"
	v.lineChannel <- "		token, _ := decoder.Token()"
	v.lineChannel <- "		if token == nil {"
	v.lineChannel <- "			break"
	v.lineChannel <- "		}"
	v.lineChannel <- "		switch se := token.(type) {"
	v.lineChannel <- "		case xml.StartElement:"

	childrenAvailable := false

	for _, node := range v.root.children {
		for _, childNode := range node.children {
			attributes := globalTagAttributes[childNode.space+childNode.name]
			if len(childNode.children) > 0 || len(attributes) > 0 {
				v.printTokenExtractor(childNode)
				childrenAvailable = true
			}
		}
	}
	if !childrenAvailable {
		for _, node := range v.root.children {
			v.printTokenExtractor(node)
		}
	}

	v.lineChannel <- "		}"

	v.lineChannel <- "	}"
	v.lineChannel <- "}"

	v.lineChannel <- "func genericReader(filename string) (io.Reader, *os.File, error) {"
	v.lineChannel <- "    file, err := os.Open(filename)"
	v.lineChannel <- "    if err != nil {"
	v.lineChannel <- "        return nil, nil, err"
	v.lineChannel <- "    }"
	v.lineChannel <- "    if strings.HasSuffix(filename, \"bz2\") {"
	v.lineChannel <- "        return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err"
	v.lineChannel <- "    }"
	v.lineChannel <- ""
	v.lineChannel <- "    if strings.HasSuffix(filename, \"gz\") {"
	v.lineChannel <- "        reader, err := gzip.NewReader(bufio.NewReader(file))"
	v.lineChannel <- "        if err != nil {"
	v.lineChannel <- "            return nil, nil, err"
	v.lineChannel <- "        }"
	v.lineChannel <- "        return bufio.NewReader(reader), file, err"
	v.lineChannel <- "    }"
	v.lineChannel <- "    return bufio.NewReader(file), file, err"
	v.lineChannel <- "}"
}

func (v *Verify) printTokenExtractor(node *Node) {
	v.lineChannel <- "			if se.Name.Local == \"" + node.name + "\" && se.Name.Space == \"" + node.space + "\" {"
	//v.lineChannel <- "				var item " + v.makeName(name)
	v.lineChannel <- "				var item " + node.makeType(v.namePrefix, v.nameSuffix)
	//v.lineChannel <- "				var item " + node.name
	v.lineChannel <- "				decoder.DecodeElement(&item, &se)"
	v.lineChannel <- "				b, err := json.Marshal(item)"
	v.lineChannel <- "				if err != nil {"
	v.lineChannel <- "				    log.Fatal(err)"
	v.lineChannel <- "				}"
	v.lineChannel <- "				fmt.Println(string(b))"
	//	v.lineChannel <- "				fmt.Printf(\"%+v\\n\", item)"
	v.lineChannel <- ""
	v.lineChannel <- "			}"
}
