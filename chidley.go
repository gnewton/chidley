package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"sort"
	//	"strconv"
	"strings"
)

const STYPE = "_Type"

var globalTagAttributes map[string](map[string]bool)
var nameSpaceIds map[string]int
var nameSpaceTagMap map[string]string = make(map[string]string)

func main() {
	//xmlFile, err := os.Open("pubmed_xml_12550251")
	//xmlFile, err := os.Open("nrn_rrn_ab_kml_en.kml")
	//xmlFile, err := os.Open("xml/b.xml")
	//xmlFile, err := os.Open("/home/newtong/2014/Carleton/thesis/data/pubmed/nlmmedlinecitationset_140101.dtd")
	//xmlFile, err := os.Open("/home/newtong/work/kml/a.xml")
	//xmlFile, err := os.Open("xml/pubmed_xml_12750255")
	//xmlFile, err := os.Open("xml/bookCatalog.xml")
	xmlFile, err := os.Open("xml/ns1.xml")

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

	globalTagAttributes = make(map[string](map[string]bool))
	nameSpaceIds = make(map[string]int)

	startingNameSpaceId := 0

	this := root
	var last *Node
	last = nil
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.Comment:
			//fmt.Println("Comment: ", this.name)
		case xml.ProcInst:
			//fmt.Println("ProcInst: ", this.name)
		case xml.Directive:
			//fmt.Println("Directive: ", this.name)
		case xml.CharData:
			this.nodeTypeInfo.checkFieldType(string(startElement))
			//fmt.Println(this.name + ":[" + strings.TrimSpace(string(startElement)) + "]")
		case xml.StartElement:
			name := startElement.Name.Local
			space := startElement.Name.Space
			fmt.Println("===========================================")
			fmt.Printf("%+v\n", token)
			//fmt.Println(token, "  |||  ", name, "  |||  ", space)

			findNewNameSpaces(startElement.Attr)

			for _, attr := range startElement.Attr {
				fmt.Println(attr.Name.Local, attr.Name.Space, attr.Value)
			}

			_, ok := nameSpaceIds[space]
			if !ok {
				nameSpaceIds[space] = startingNameSpaceId
				startingNameSpaceId += 1
			}

			//fmt.Println(space, name)
			//for _, att := range startElement.Attr {
			//fmt.Println(att.Name.Space, "|", att.Name.Local, "|", att.Value)
			//}

			if name == "" {
				continue
			}
			var child *Node
			var attributes map[string]bool
			child, ok = this.children[space+name]
			if ok {
				this.childCount[space+name] += 1
				attributes = globalTagAttributes[space+name]
				//fmt.Println("Exists", name)
			} else {
				//fmt.Println("!Exists", name)
				newNode := new(Node)
				spaceTag, _ := nameSpaceTagMap[space]
				// if !ok {
				// 	spaceTag = ""
				// }
				newNode.initialize(name, space, spaceTag, this)
				child = newNode
				this.children[space+name] = child
				this.childCount[space+name] = 1

				attributes = make(map[string]bool)
				globalTagAttributes[space+name] = attributes
			}

			for _, attr := range startElement.Attr {
				attributes[attr.Name.Local] = true
				//fmt.Println(child.name + ":: " + attr.Name.Local)
			}
			this = child
			if last == child {
				//fmt.Println("repreat: " + last.name)
			}
			depth += 1
			//fmt.Println(indent(depth) + startElement.Name.Local)
			//if startElement.Name.Local == "entry" {
			// do what you need to do for each entry
			//}
		case xml.EndElement:
			//fmt.Println(indent(depth) + "-" + startElement.Name.Local)
			for key, c := range this.childCount {
				if c > 1 {
					//fmt.Println("**** repreat: " + key)
					this.children[key].repeats = true
				}
				this.childCount[key] = 0
			}

			last = this
			if this.parent != nil {
				this = this.parent
			}

			//if startElement.Name.Local == "entry" {
			// do what you need to do for each entry
			//}
		}
	}

	//printTree(root, 0, "Document", false)
	//printTree(root, 0, "Document", true)

	//printStruct(root, "Document", false)
	alreadyPrinted := make(map[string]bool)
	printStruct(root, "Document", true, alreadyPrinted)

	fmt.Printf("%+v\n", nameSpaceTagMap)

}

func printTree(n *Node, d int, startName string, foundStartString bool) {
	if n.name == startName {
		foundStartString = true
	}
	repeats := ""
	if n.repeats {
		repeats = "*"
	}
	if foundStartString {
		fmt.Println(indent(d) + n.name + repeats)
		d += 1
	}

	for _, v := range n.children {
		printTree(v, d, startName, foundStartString)
	}
}

func printStruct(n *Node, startName string, foundStartString bool, alreadyPrinted map[string]bool) {
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
		//if len(n.children) > 0 || len(n.attributes) > 0 {
		if len(n.children) > 0 || len(attributes) > 0 {
			printStruct0(n)
		}
	}

	for _, v := range n.children {
		printStruct(v, startName, foundStartString, alreadyPrinted)
	}
}

func makeAttributes(attributes map[string]bool) []string {
	all := make([]string, 0)
	for att, _ := range attributes {
		name := capitalizeFirstLetter(att)
		attStr := "\t" + name + " string `xml:\"" + att + ",attr\"`"
		all = append(all, attStr)
		//fmt.Println("APPEND: ", all)
		//fmt.Println(attStr)
	}
	return all
}

func makeType(nti *NodeTypeInfo) string {
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

func printStruct0(n *Node) {

	//nameSpaceId, ok := nameSpaceIds[n.space]
	// var nameSpaceIdString string
	// if ok {
	// 	nameSpaceIdString = strconv.Itoa(nameSpaceId)
	// } else {
	// 	nameSpaceIdString = ""
	// }
	//fmt.Printf("11 %+v\n", n)
	//fmt.Printf("22 %+v\n", n.nodeTypeInfo)
	fmt.Println("type " + n.makeType() + " struct {")

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
			field += makeType(v.nodeTypeInfo)

		} else {
			if v.repeats {
				field += "[]"
				//fmt.Print("[]")
			}
			//field += name + STYPE
			field += v.makeType()
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
		fmt.Println(fields[i])
	}

	//fmt.Println("}")
	//field += "}"
	//fmt.Println(field)
	fmt.Println("}")
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
