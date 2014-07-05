package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const STYPE = "_Type"

var globalTagAttributes map[string](map[string]bool)

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
	root.initialize("root", "", nil)

	globalTagAttributes = make(map[string](map[string]bool))

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
			//fmt.Println(space, name)
			//for _, att := range startElement.Attr {
			//fmt.Println(att.Name.Space, "|", att.Name.Local, "|", att.Value)
			//}

			if name == "" {
				continue
			}
			var child *Node
			var attributes map[string]bool
			child, ok := this.children[name]
			if ok {
				this.childCount[name] += 1
				attributes = globalTagAttributes[name]
				//fmt.Println("Exists", name)
			} else {
				//fmt.Println("!Exists", name)
				newNode := new(Node)
				newNode.initialize(name, space, this)
				child = newNode
				this.children[name] = child
				this.childCount[name] = 1

				attributes = make(map[string]bool)
				globalTagAttributes[name] = attributes
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

	alreadyPrinted[n.name] = true

	if n.name == startName {
		foundStartString = true
	}

	attributes := globalTagAttributes[n.name]
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
	//fmt.Printf("11 %+v\n", n)
	//fmt.Printf("22 %+v\n", n.nodeTypeInfo)
	fmt.Println("type " + capitalizeFirstLetter(n.name) + STYPE + " struct {")

	//fields := makeAttributes(n.attributes)
	attributes := globalTagAttributes[n.name]
	fields := makeAttributes(attributes)

	xmlName := "\tXMLName  xml.Name `xml:\"" + n.name + ",omitempty\" json:\",omitempty\"`"
	fields = append(fields, xmlName)

	//fmt.Println(fields)

	var field string
	for _, v := range n.children {
		name := capitalizeFirstLetter(v.name)
		field = "\t" + name + " "
		//fmt.Print()
		childAttributes := globalTagAttributes[v.name]
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
			field += name + STYPE
			//fmt.Print()
		}
		xmlString := " `xml:\"" + v.name + ",omitempty\" json:\",omitempty\"`"

		field += xmlString

		fields = append(fields, field)

		//fmt.Println(xmlString)
	}

	//if len(n.children) == 0 && len(n.attributes) > 0 {
	if len(n.children) == 0 && len(attributes) > 0 {
		xmlString := " `xml:\"chardata,omitempty\" json:\",omitempty\"`"
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
