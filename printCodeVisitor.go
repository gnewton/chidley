package main

import (
	"sort"
)

type PrintCodeVisitor struct {
	alreadyVisited      map[string]bool
	globalTagAttributes map[string]([]*FQN)
	lineChannel         chan string
}

func (v *PrintCodeVisitor) init(lineChannel chan string) {
	v.alreadyVisited = make(map[string]bool)
	v.globalTagAttributes = make(map[string]([]*FQN))
	v.lineChannel = lineChannel

}

func (v *PrintCodeVisitor) Visit(node *Node) bool {
	if v.AlreadyVisited(node) {
		return false
	}

	v.SetAlreadyVisited(node)

	attributes := v.globalTagAttributes[nk(node)]
	if len(node.children) > 0 || len(attributes) > 0 {
		v.lineChannel <- "type " + node.makeType(namePrefix, nameSuffix) + " struct {"
		v.printInternalFields(node)
		v.lineChannel <- "}\n"
	}
	for _, child := range node.children {
		v.Visit(child)
	}
	return true
}

func (v *PrintCodeVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintCodeVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
}

func (pn *PrintCodeVisitor) printInternalFields(n *Node) {
	attributes := pn.globalTagAttributes[nk(n)]
	//fields := makeAttributes(attributes)
	var fields []string
	var xmlName string
	if n.space != "" {
		xmlName = "\tXMLName  xml.Name `xml:\"" + n.space + " " + n.name + ",omitempty\" json:\",omitempty\"`"
	} else {
		//xmlName = "\tXMLName  xml.Name `xml:\"" + n.name + ",omitempty\" json:\",omitempty\"`"
	}
	fields = append(fields, xmlName)
	var field string

	for _, v := range n.children {
		field = "\t" + v.makeType(namePrefix, namePrefix) + " "
		childAttributes := pn.globalTagAttributes[v.space+v.name]
		if len(v.children) == 0 && len(childAttributes) == 0 {
			if v.repeats {
				field += "[]"
			}
			field += findType(v.nodeTypeInfo, false)

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
		pn.lineChannel <- fields[i]
	}
}
