package main

import (
	"os"
)

type CodeGenerator struct {
	writer          Writer
	lineChannel     chan string
	source          Source
	root            *Node
	namePrefix      string
	nameSuffix      string
	nameSpaceTagMap map[string]string
	globalNodeMap   map[string]*Node
}

func (v *CodeGenerator) init(root *Node, globalNodeMap map[string]*Node, namePrefix string, nameSpaceTagMap map[string]string, nameSuffix string, codegenDir string, codegenGoFile string, source Source, lineChannel chan string) (Writer, error) {
	v.root = root
	v.namePrefix = namePrefix
	v.nameSuffix = nameSuffix
	v.lineChannel = lineChannel
	v.source = source
	v.nameSpaceTagMap = nameSpaceTagMap
	v.globalNodeMap = globalNodeMap

	var err error
	err = os.RemoveAll(codegenDir)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(codegenDir); os.IsNotExist(err) {
		err = os.Mkdir(codegenDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	v.writer = new(fileWriter)
	v.writer.open(codegenDir+"/"+codegenGoFile, lineChannel)

	return v.writer, err

}

func (v *CodeGenerator) generateCodePre(hasStartElements bool, url bool, writeJson bool, writeXml bool) {
	v.lineChannel <- "package main"
	v.lineChannel <- " "

	v.lineChannel <- "import ("
	if !hasStartElements {
		v.lineChannel <- "	\"os\""
		v.lineChannel <- "	\"log\""
		v.lineChannel <- ")"
		return
	}
	if url {
		v.lineChannel <- "	\"net/http\""
	}
	v.lineChannel <- "	\"encoding/xml\""
	if writeJson {
		v.lineChannel <- "\"encoding/json\""
	}

	if writeJson || writeXml {
		v.lineChannel <- "	\"fmt\""
	}
	v.lineChannel <- "	\"os\""
	v.lineChannel <- "      \"bufio\""
	v.lineChannel <- "      \"log\""
	v.lineChannel <- "      \"compress/bzip2\""
	v.lineChannel <- "      \"compress/gzip\""
	v.lineChannel <- "      \"io\""
	v.lineChannel <- "      \"strings\""
	v.lineChannel <- ")"

	v.lineChannel <- "\n\n"

	makeGenericReaderCode(v.lineChannel)
}

func (v *CodeGenerator) makeName(name string) string {
	//return v.namePrefix + capitalizeFirstLetter(name) + v.nameSuffix
	return v.namePrefix + cleanName(name) + v.nameSuffix
}

func (v *CodeGenerator) generateConvertToCode(firstNode *Node, writeJson bool, writeXml bool, prettyPrint bool) {
	v.lineChannel <- "func main(){"
	v.lineChannel <- " "
	v.makeDecoder(firstNode, writeJson, writeXml, prettyPrint)
	v.lineChannel <- "}"
}

func (v *CodeGenerator) generateVerifyCode(hasStartElements bool, globalTagAttributes map[string]([]*FQN), url bool) {
	v.lineChannel <- "func main(){"
	v.lineChannel <- " "
	if !hasStartElements {
		v.lineChannel <- "log.Print(\"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\")"
		v.lineChannel <- "log.Print(\"++++++ No StartElements encountered: possible error in XML or bug+++++++\")"
		v.lineChannel <- "log.Print(\"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\")"
		v.lineChannel <- "os.Exit(42)"
		v.lineChannel <- "}"
		return
	}

	for _, node := range v.globalNodeMap {
		v.makeDecoder(node, false, false, false)
	}
	v.lineChannel <- "}"

}

func (v *CodeGenerator) printTokenExtractor(name string, space string, spaceTag string, writeJson bool, writeXml bool, prettyPrint bool) {
	v.lineChannel <- "			if se.Name.Local == \"" + name + "\" && se.Name.Space == \"" + space + "\" {"
	if !writeJson && !writeXml {
		v.lineChannel <- "				count += 1"
	}
	v.lineChannel <- "				var item " + makeTypeGeneric(name, spaceTag, v.namePrefix, v.nameSuffix)
	v.lineChannel <- "				decoder.DecodeElement(&item, &se)"
	if writeJson {
		if prettyPrint {
			v.lineChannel <- "				b, err := json.MarshalIndent(item, \"\", \" \")"
		} else {
			v.lineChannel <- "				b, err := json.Marshal(item)"
		}
	} else {
		if writeXml {
			v.lineChannel <- "				tmp := struct {"
			v.lineChannel <- "				" + "C_PubmedArticle"
			v.lineChannel <- "				XMLName struct{} `xml:\"" + "PubmedArticle" + "\"`"
			v.lineChannel <- "				}{" + "C_PubmedArticle" + ": item}"

			v.lineChannel <- "				if err := enc.Encode(tmp); err != nil {"
			v.lineChannel <- "				    fmt.Printf(\"error: %v\\n\", err)"
			v.lineChannel <- "				}"
		}
	}

	if writeJson || writeXml {
		v.lineChannel <- "				if err != nil {"
		v.lineChannel <- "				    log.Fatal(err)"
		v.lineChannel <- "				}"
	}
	if writeJson {
		v.lineChannel <- "				fmt.Println(string(b))"
	}
	v.lineChannel <- ""
	v.lineChannel <- "			}"
}

func (v *CodeGenerator) makeReader(url bool) {
	if !url {
		v.makeFileReader()
	} else {
		v.makeUrlReader()
	}

}

func (v *CodeGenerator) makeUrlReader() {
	v.lineChannel <- "        res, err := http.Get(\"" + v.source.getName() + "\")"
	v.lineChannel <- "        if err != nil {"
	v.lineChannel <- "        		log.Fatal(err)"
	v.lineChannel <- "        		return"
	v.lineChannel <- "        	}"
	v.lineChannel <- "        if res.StatusCode != 200 {"
	v.lineChannel <- "        		log.Fatal(\"ERROR: bad http status code != 200: \", res.StatusCode)"
	v.lineChannel <- "        		return"
	v.lineChannel <- "        	}"
	v.lineChannel <- "        reader := res.Body"
}

func (v *CodeGenerator) makeFileReader() {

	v.lineChannel <- "reader, xmlFile, err := genericReader(\"../" + v.source.getName() + "\")"

	v.lineChannel <- "if err != nil {"
	v.lineChannel <- "log.Fatal(err)"
	v.lineChannel <- "return"
	v.lineChannel <- "}"
	v.lineChannel <- "defer xmlFile.Close()"
}

func (v *CodeGenerator) makeCountVar(node *Node) string {
	return "count_" + node.name + "_" + v.nameSpaceTagMap[node.space]
}

func (v *CodeGenerator) makeTagSpaceMap(varName string) string {
	s := "\n\t" + varName + ":= map[string]string{"
	for _, v := range v.globalNodeMap {
		s += "\n\t\t\"" + v.name + "\":\"" + v.space + "\","
	}
	s += "\n\t}"
	return s
}

func (v *CodeGenerator) makeTagSpaceTagMap(nameSpaceTagMap map[string]string, varName string) string {
	s := "\n\t" + varName + ":= map[string]string{"
	for k, v := range nameSpaceTagMap {
		s += "\n\t\t\"" + k + "\":\"" + v + "\","
	}
	s += "\n\t}"
	return s
}

func (v *CodeGenerator) makeDecoder(node *Node, writeJson bool, writeXml bool, prettyPrint bool) {
	v.lineChannel <- "     {"
	if !writeJson && !writeXml {
		v.lineChannel <- "       count:= 0"
	}

	v.makeReader(url)
	if writeXml {
		v.lineChannel <- "enc := xml.NewEncoder(os.Stdout)"
		if prettyPrint {
			v.lineChannel <- "enc.Indent(\"  \", \"    \")"
		}
	}

	v.lineChannel <- "decoder := xml.NewDecoder(reader)"

	v.lineChannel <- "	for{"
	v.lineChannel <- "		token, _ := decoder.Token()"
	v.lineChannel <- "		if token == nil {"
	v.lineChannel <- "		    break"
	v.lineChannel <- "		}"
	v.lineChannel <- "		switch se := token.(type) {"
	v.lineChannel <- "		case xml.StartElement:"

	v.printTokenExtractor(node.name, node.space, node.spaceTag, writeJson, writeXml, prettyPrint)
	v.lineChannel <- "	   }"
	v.lineChannel <- "        }"
	spaceTag := ""
	if node.spaceTag != "" {
		spaceTag = node.spaceTag + ":"
	}

	if !writeJson && !writeXml {
		v.lineChannel <- "     log.Print(\"Number of " + spaceTag + node.name + "= \", count)"
	}

	v.lineChannel <- "    }"
}

func makeGenericReaderCode(lineChannel chan string) {
	lineChannel <- "func genericReader(filename string) (io.Reader, *os.File, error) {"
	lineChannel <- "    file, err := os.Open(filename)"
	lineChannel <- "    if err != nil {"
	lineChannel <- "        return nil, nil, err"
	lineChannel <- "    }"
	lineChannel <- "    if strings.HasSuffix(filename, \"bz2\") {"
	lineChannel <- "        return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err"
	lineChannel <- "    }"
	lineChannel <- ""
	lineChannel <- "    if strings.HasSuffix(filename, \"gz\") {"
	lineChannel <- "        reader, err := gzip.NewReader(bufio.NewReader(file))"
	lineChannel <- "        if err != nil {"
	lineChannel <- "            return nil, nil, err"
	lineChannel <- "        }"
	lineChannel <- "        return bufio.NewReader(reader), file, err"
	lineChannel <- "    }"
	lineChannel <- "    return bufio.NewReader(file), file, err"
	lineChannel <- "}"
}
