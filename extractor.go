package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"sort"
	"strings"
)

var nameMapper = map[string]string{
	"-": "_",
	".": "_dot_",
}

type Extractor struct {
	globalTagAttributes map[string](map[string]string)
	globalNodeMap       map[string]*Node
	namePrefix          string
	nameSpaceTagMap     map[string]string
	nameSuffix          string
	verify              bool
	xmlFilename         string
}

func (ex *Extractor) extract() error {
	reader, xmlFile, err := genericReader(ex.xmlFilename)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer xmlFile.Close()

	ex.globalTagAttributes = make(map[string](map[string]string))
	ex.nameSpaceTagMap = make(map[string]string)
	ex.globalNodeMap = make(map[string]*Node)

	decoder := xml.NewDecoder(reader)
	depth := 0

	root := new(Node)
	root.initialize("root", "", "", nil)

	this := root
	//var last *Node
	//last = nil

	hasStartElements := false

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		if token == nil {
			break
		}
		//fmt.Printf("token: %+v\n", token)
		switch element := token.(type) {
		case xml.Comment:
			if DEBUG {
				fmt.Printf("Comment: %+v\n", string(element))
			}
		case xml.ProcInst:
			if DEBUG {
				fmt.Printf("ProcInst: %+v\n", element)
			}
		case xml.Directive:
			if DEBUG {
				fmt.Printf("Directive: %+v\n", string(element))
			}
		case xml.CharData:
			if DEBUG {
				fmt.Printf("CharData: %+v\n", string(element))
			}
			this.nodeTypeInfo.checkFieldType(string(element))
		case xml.StartElement:
			if DEBUG {
				fmt.Printf("StartElement: %+v\n", element)
			}
			hasStartElements = true

			if element.Name.Local == "" {
				continue
			}
			this = ex.handleStartElement(element, this)
			depth += 1

		case xml.EndElement:
			if DEBUG {
				fmt.Printf("endElement: %+v\n", element)
			}
			for key, c := range this.childCount {
				if c > 1 {
					this.children[key].repeats = true
				}
				this.childCount[key] = 0
			}
			if this.parent != nil {
				this = this.parent
			}
		}
	}

	alreadyPrinted := make(map[string]bool)
	var writer Writer

	lineChannel := make(chan string)
	var verify *Verify
	if ex.verify {
		verify = new(Verify)
		writer, err = verify.init(root, namePrefix, nameSuffix, "chidleyVerity", "ChidleyVerify.go", ex.xmlFilename, lineChannel)
		verify.generateCodePre(hasStartElements)
		if err != nil {
			fmt.Println(err)
			return err
		}

	} else {
		writer = new(stdoutWriter)
		writer.open("", lineChannel)
	}

	ex.printStruct(root, lineChannel, "Document", true, alreadyPrinted)
	if ex.verify {
		verify.generatePost(hasStartElements, ex.globalTagAttributes)
	}

	//ex.printTree(root, lineChannel, 0, "", true)

	close(lineChannel)
	writer.close()

	return nil
}

func (ex *Extractor) printTree(n *Node, lineChannel chan string, d int, startName string, foundStartString bool) {
	if n.name == startName {
		foundStartString = true
	}
	repeats := ""
	if n.repeats {
		repeats = "*"
	}
	if foundStartString {
		fmt.Println(indent(d) + n.name + repeats)
		lineChannel <- indent(d) + n.name + repeats
		d += 1
	}

	for _, v := range n.children {
		ex.printTree(v, lineChannel, d, startName, foundStartString)
	}
}

func (ex *Extractor) printStruct(n *Node, lineChannel chan string, startName string, foundStartString bool, alreadyPrinted map[string]bool) {
	_, ok := alreadyPrinted[n.space+n.name]
	if ok {
		//fmt.Println("printStruct: already printed: " + n.space + ":" + n.name)
		return
	}
	//fmt.Println(">> printStruct: printing: " + n.space + ":" + n.name)

	alreadyPrinted[n.space+n.name] = true

	if n.space+n.name == startName {
		foundStartString = true
	}
	attributes := ex.globalTagAttributes[n.space+n.name]
	//if n.parent != nil && foundStartString {
	if foundStartString {
		if len(n.children) > 0 || len(attributes) > 0 {
			lineChannel <- "type " + n.makeType(namePrefix, nameSuffix) + " struct {"
			ex.printInternalFields(n, lineChannel)
			lineChannel <- "}\n"
		}
	}

	for _, v := range n.children {
		ex.printStruct(v, lineChannel, startName, foundStartString, alreadyPrinted)
	}
}

func makeAttributes(attributes map[string]string) []string {
	all := make([]string, 0)
	for att, space := range attributes {
		//name := capitalizeFirstLetter(att)
		name := att
		if space != "" {
			space = space + " "
		}
		attStr := "\t" + attributePrefix + cleanName(name) + " string `xml:\"" + space + att + ",attr\"`"
		all = append(all, attStr)
		//fmt.Println("APPEND: ", all)
		//fmt.Println(attStr)
	}
	return all
}

func findType(nti *NodeTypeInfo) string {
	return "string"

	if nti.alwaysBool {
		return "bool"
	}

	if nti.alwaysInt08 {
		return "int8"
	}
	if nti.alwaysInt16 {
		return "int16"
	}
	if nti.alwaysInt32 {
		return "int32"
	}
	if nti.alwaysInt64 {
		return "int64"
	}

	if nti.alwaysInt0 {
		return "int"
	}

	if nti.alwaysFloat32 {
		return "float32"
	}
	if nti.alwaysFloat64 {
		return "float64"
	}

	return "string"
}

func (ex *Extractor) printInternalFields(n *Node, lineChannel chan string) {
	//fmt.Println("type " + n.makeType(namePrefix) + " struct {")

	//fields := makeAttributes(n.attributes)
	attributes := ex.globalTagAttributes[n.space+n.name]
	fields := makeAttributes(attributes)
	var xmlName string
	if n.space != "" {
		xmlName = "\tXMLName  xml.Name `xml:\"" + n.space + " " + n.name + ",omitempty\" json:\",omitempty\"`"
	} else {
		//xmlName = "\tXMLName  xml.Name `xml:\"" + n.name + ",omitempty\" json:\",omitempty\"`"
	}
	fields = append(fields, xmlName)

	var field string

	for _, v := range n.children {
		field = "\t " + v.makeType(namePrefix, namePrefix) + " "
		childAttributes := ex.globalTagAttributes[v.space+v.name]
		if len(v.children) == 0 && len(childAttributes) == 0 {
			if v.repeats {
				field += "[]"
			}
			field += findType(v.nodeTypeInfo)

		} else {
			if v.repeats {
				field += "[]"
			} else {
				field += "*"
			}
			field += v.makeType(namePrefix, namePrefix)
		}
		var xmlString string
		if v.space != "" {
			xmlString = " `xml:\"" + v.space + " " + v.name + ",omitempty\" json:\",omitempty\"`"
		} else {
			xmlString = " `xml:\"" + v.name + ",omitempty\" json:\",omitempty\"`"
		}
		field += xmlString
		fields = append(fields, field)
	}

	if len(n.children) == 0 && len(attributes) > 0 {
		xmlString := " `xml:\",chardata\" json:\",omitempty\"`"
		//charField := "\t" + capitalizeFirstLetter(n.name) + " string" + xmlString
		charField := "\t" + "Text" + " string" + xmlString
		fields = append(fields, charField)
	}
	sort.Strings(fields)

	for i := 0; i < len(fields); i++ {
		lineChannel <- fields[i]
		//fmt.Println(fields[i])
	}
}

func (ex *Extractor) findNewNameSpaces(attrs []xml.Attr) {
	for _, attr := range attrs {
		if attr.Name.Space == "xmlns" {
			ex.nameSpaceTagMap[attr.Value] = attr.Name.Local
		}
	}

}

func cleanName(name string) string {
	for old, new := range nameMapper {
		name = strings.Replace(name, old, new, -1)
	}
	return name
}

func (ex *Extractor) handleStartElement(startElement xml.StartElement, this *Node) *Node {
	name := startElement.Name.Local
	space := startElement.Name.Space

	ex.findNewNameSpaces(startElement.Attr)

	var child *Node
	var attributes map[string]string
	child, ok := this.children[space+name]
	// Does this node already exist as child
	//fmt.Println(space, name)
	if ok {
		this.childCount[space+name] += 1
		attributes = ex.globalTagAttributes[space+name]
		//fmt.Println("Exists", name)
	} else {
		// if this node does not exist as child, it may still exist as child on other node:
		child, ok = ex.globalNodeMap[space+name]
		if !ok {
			child = new(Node)
			ex.globalNodeMap[space+name] = child
			spaceTag, _ := ex.nameSpaceTagMap[space]
			child.initialize(name, space, spaceTag, this)
			this.childCount[space+name] = 1

			attributes = make(map[string]string)
			ex.globalTagAttributes[space+name] = attributes
		} else {
			attributes = ex.globalTagAttributes[space+name]
			//fmt.Println("Exists global:", space, name)
		}
		this.children[space+name] = child

	}

	for _, attr := range startElement.Attr {
		attributes[attr.Name.Local] = attr.Name.Space
		//fmt.Println("Attr: name=" + attr.Name.Local + " Space=" + attr.Name.Space)
		//fmt.Println(child.name + ":: " + attr.Name.Local)
	}
	return child
}
