package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var nameSuffix = "_Type"
var namePrefix = "Chidley_"
var VERIFY = false

var globalTagAttributes map[string](map[string]string)
var nameSpaceTagMap map[string]string = make(map[string]string)

type Writer interface {
	open(s string, lineChannel chan string) error
	close()
}

func init() {
	flag.StringVar(&nameSuffix, "s", nameSuffix, "Suffix to element names")
	flag.StringVar(&namePrefix, "p", namePrefix, "Prefix to element names")
	flag.BoolVar(&VERIFY, "V", VERIFY, "Do full code generation & see if it can decode the original source file")
}

func handleParameters() bool {
	flag.Parse()
	//fmt.Println("nameSuffix=" + nameSuffix)
	//fmt.Println("namePrefix=" + namePrefix)
	return true
}

func main() {
	handleParameters()

	//xmlFilename := "xml/bookCatalog.xml"
	xmlFilename := "xml/pubmed_xml_12750255"
	//xmlFile, err := os.Open("pubmed_xml_12550251")
	//xmlFile, err := os.Open("nrn_rrn_ab_kml_en.kml")
	//xmlFile, err := os.Open("xml/b.xml")
	//xmlFile, err := os.Open("/home/newtong/2014/Carleton/thesis/data/pubmed/nlmmedlinecitationset_140101.dtd")
	//xmlFile, err := os.Open("/home/newtong/work/kml/a.xml")
	//xmlFile, err := os.Open("xml/pubmed_xml_12750255")
	xmlFile, err := os.Open(xmlFilename)
	//xmlFile, err := os.Open("xml/ns1.xml")
	//xmlFile, err := os.Open("xml/MODIS-Imagery-OilSpill.kml")

	//xmlFile, err := os.Open("catalog.xml")
	//xmlFile, err := os.Open("simple1.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	depth := 0

	root := new(Node)
	root.initialize("root", "", "", nil)

	globalTagAttributes = make(map[string](map[string]string))

	this := root
	//var last *Node
	//last = nil

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.Comment:
			//fmt.Printf("Comment: %+v\n", string(startElement))
		case xml.ProcInst:
			//fmt.Printf("ProcInst: %+v\n", startElement)
		case xml.Directive:
			//fmt.Printf("Directive: %+v\n", string(startElement))
		case xml.CharData:
			//fmt.Printf("CharData: %+v\n", string(startElement))
			this.nodeTypeInfo.checkFieldType(string(startElement))
		case xml.StartElement:
			//fmt.Printf("StartElement: %+v\n", startElement)
			name := startElement.Name.Local
			if name == "" {
				continue
			}
			space := startElement.Name.Space

			findNewNameSpaces(startElement.Attr)

			//for _, attr := range startElement.Attr {
			//fmt.Println(attr.Name.Local, attr.Name.Space, attr.Value)
			//}

			var child *Node
			var attributes map[string]string
			//fmt.Println("Space=" + space + " Name=" + name)
			child, ok := this.children[space+name]
			if ok {
				this.childCount[space+name] += 1
				attributes = globalTagAttributes[space+name]
				//fmt.Println("Exists", name)
			} else {
				//fmt.Println("!Exists", name)
				newNode := new(Node)
				spaceTag, _ := nameSpaceTagMap[space]

				newNode.initialize(name, space, spaceTag, this)
				child = newNode
				this.children[space+name] = child
				this.childCount[space+name] = 1

				attributes = make(map[string]string)
				globalTagAttributes[space+name] = attributes
			}

			for _, attr := range startElement.Attr {
				attributes[attr.Name.Local] = attr.Name.Space
				//fmt.Println("Attr: name=" + attr.Name.Local + " Space=" + attr.Name.Space)
				//fmt.Println(child.name + ":: " + attr.Name.Local)
			}
			this = child
			depth += 1

		case xml.EndElement:
			for key, c := range this.childCount {
				if c > 1 {
					this.children[key].repeats = true
				}
				this.childCount[key] = 0
			}

			//last = this
			if this.parent != nil {
				this = this.parent
			}
		}
	}

	alreadyPrinted := make(map[string]bool)
	var writer Writer

	lineChannel := make(chan string)
	var verify *Verify
	if VERIFY {
		verify = new(Verify)
		writer, err = verify.init(root, namePrefix, nameSuffix, "newDir", xmlFilename, lineChannel)

		if err != nil {
			fmt.Println(err)
		}

	} else {
		writer = new(stdoutWriter)
		writer.open("", lineChannel)
	}

	printStruct(root, lineChannel, "Document", true, alreadyPrinted)
	if VERIFY {

	}
	verify.generatePost()
	close(lineChannel)
	writer.close()

	//fmt.Printf("%+v\n", nameSpaceTagMap)

}

func printTree(n *Node, lineChannel chan string, d int, startName string, foundStartString bool) {
	if n.name == startName {
		foundStartString = true
	}
	repeats := ""
	if n.repeats {
		repeats = "*"
	}
	if foundStartString {
		//fmt.Println(indent(d) + n.name + repeats)
		lineChannel <- indent(d) + n.name + repeats
		d += 1
	}

	for _, v := range n.children {
		printTree(v, lineChannel, d, startName, foundStartString)
	}
}

func printStruct(n *Node, lineChannel chan string, startName string, foundStartString bool, alreadyPrinted map[string]bool) {
	_, ok := alreadyPrinted[n.name]
	if ok {
		return
	}

	alreadyPrinted[n.space+n.name] = true

	if n.space+n.name == startName {
		foundStartString = true
	}

	attributes := globalTagAttributes[n.space+n.name]
	if n.parent != nil && foundStartString {
		if len(n.children) > 0 || len(attributes) > 0 {
			printStruct0(n, lineChannel)
		}
	}

	for _, v := range n.children {
		printStruct(v, lineChannel, startName, foundStartString, alreadyPrinted)
	}
}

func makeAttributes(attributes map[string]string) []string {
	all := make([]string, 0)
	for att, space := range attributes {
		name := capitalizeFirstLetter(att)
		if space != "" {
			space = space + " "
		}
		attStr := "\t" + name + " string `xml:\"" + space + att + ",attr\"`"
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

func printStruct0(n *Node, lineChannel chan string) {
	//fmt.Println("type " + n.makeType(namePrefix) + " struct {")
	lineChannel <- "type " + n.makeType(namePrefix, nameSuffix) + " struct {"

	//fields := makeAttributes(n.attributes)
	attributes := globalTagAttributes[n.space+n.name]
	fields := makeAttributes(attributes)
	var xmlName string
	if n.space != "" {
		xmlName = "\tXMLName  xml.Name `xml:\"" + n.space + " " + n.name + ",omitempty\" json:\",omitempty\"`"
	} else {
		xmlName = "\tXMLName  xml.Name `xml:\"" + n.name + ",omitempty\" json:\",omitempty\"`"
	}
	fields = append(fields, xmlName)

	//fmt.Println(fields)

	var field string
	for _, v := range n.children {
		//nameSpaceId, ok := nameSpaceIds[v.space]
		// if ok {
		// 	nameSpaceIdString = strconv.Itoa(nameSpaceId)
		// } else {
		// 	nameSpaceIdString = ""
		// }
		//name := capitalizeFirstLetter(v.name)
		//field = "\t" + name + nameSpaceIdString + " "
		field = "\t" + v.makeName() + " "
		//fmt.Print()
		childAttributes := globalTagAttributes[v.space+v.name]
		//if len(v.children) == 0 && len(v.attributes) == 0 {
		if len(v.children) == 0 && len(childAttributes) == 0 {
			if v.repeats {
				//fmt.Print("[]")
				field += "[]"
			}
			//fmt.Print(makeType(v.nodeTypeInfo))
			field += findType(v.nodeTypeInfo)

		} else {
			if v.repeats {
				field += "[]"
				//fmt.Print("[]")
			}
			//field += name + nameSuffix
			field += v.makeType(namePrefix, namePrefix)
			//fmt.Print()
		}
		var xmlString string
		if v.space != "" {
			xmlString = " `xml:\"" + v.space + " " + v.name + ",omitempty\" json:\",omitempty\"`"
		} else {
			xmlString = " `xml:\"" + v.name + ",omitempty\" json:\",omitempty\"`"
		}

		field += xmlString

		fields = append(fields, field)

		//fmt.Println(xmlString)
	}

	//if len(n.children) == 0 && len(n.attributes) > 0 {
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

	//fmt.Println("}")
	//field += "}"
	//fmt.Println(field)

	//fmt.Println("}\n")
	lineChannel <- "}\n"
}

func attributes(atts map[string]bool) string {
	ret := ": "

	for k, _ := range atts {
		ret = ret + k + ", "
	}
	return ret
}

func indent(d int) string {
	indent := ""
	for i := 0; i < d; i++ {
		indent = indent + "\t"
	}
	return indent
}

func capitalizeFirstLetter(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func findNewNameSpaces(attrs []xml.Attr) {
	for _, attr := range attrs {
		if attr.Name.Space == "xmlns" {
			nameSpaceTagMap[attr.Value] = attr.Name.Local
		}
	}

}
