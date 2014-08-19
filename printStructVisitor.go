package main

import (
	"sort"
)

type PrintStructVisitor struct {
	alreadyVisited      map[string]bool
	globalTagAttributes map[string]([]*FQN)
	lineChannel         chan string
	maxDepth            int
	depth               int
	nameSpaceTagMap     map[string]string
	useType             bool
}

func (v *PrintStructVisitor) init(lineChannel chan string, maxDepth int, globalTagAttributes map[string]([]*FQN), nameSpaceTagMap map[string]string, useType bool) {
	v.alreadyVisited = make(map[string]bool)
	v.globalTagAttributes = make(map[string]([]*FQN))
	v.globalTagAttributes = globalTagAttributes
	v.lineChannel = lineChannel
	v.maxDepth = maxDepth
	v.depth = 0
	v.nameSpaceTagMap = nameSpaceTagMap
	v.useType = useType
}

func (v *PrintStructVisitor) Visit(node *Node) bool {
	v.depth += 1
	//	if v.depth >= v.maxDepth {
	//		return false
	//	}

	if v.AlreadyVisited(node) {
		v.depth += 1
		return false
	}
	v.SetAlreadyVisited(node)

	attributes := v.globalTagAttributes[nk(node)]

	//if len(node.children) > 0 || len(attributes) > 0 {
	if true {
		//if true {
		//if len(node.children) > 0 {

		v.lineChannel <- "type " + node.makeType(namePrefix, nameSuffix) + " struct {"
		makeAttributes(v.lineChannel, attributes, v.nameSpaceTagMap)
		//if v.depth+1 < v.maxDepth {

		v.printInternalFields(node)
		//}
		if node.space != "" {
			v.lineChannel <- "\tXMLName  xml.Name `xml:\"" + node.space + " " + node.name + ",omitempty\" json:\",omitempty\"`"
		} else {
			//xmlName = "\tXMLName  xml.Name `xml:\"" + n.name + ",omitempty\" json:\",omitempty\"`"
		}
		v.lineChannel <- "}\n"
		v.lineChannel <- "\n"
	}
	for _, child := range node.children {
		v.Visit(child)
	}
	v.depth += 1
	return true
}

func (v *PrintStructVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintStructVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
}

func (pn *PrintStructVisitor) printInternalFields(n *Node) {
	//attributes := pn.globalTagAttributes[nk(n)]
	//fields := makeAttributes(attributes)
	var fields []string

	var field string

	for _, v := range n.children {
		field = "\t" + v.makeType(namePrefix, nameSuffix) + " "

		//childAttributes := pn.globalTagAttributes[nk(v)]
		//if len(v.children) == 0 && len(childAttributes) == 0 {
		if false {
			//if false {
			//if len(v.children) == 0 {
			if v.repeats {
				field += "[]"
			}
			field += findType(v.nodeTypeInfo, useType)
		} else {
			if v.repeats {
				field += "[]"
			} else {
				field += "*"
			}
			field += v.makeType(namePrefix, nameSuffix)
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
	//if (len(n.children) == 0 && len(attributes) > 0) || n.hasCharData {
	if n.hasCharData {
		xmlString := " `xml:\",chardata\" json:\",omitempty\"`"
		//charField := "\t" + capitalizeFirstLetter(n.name) + " string" + xmlString

		charField := "\t" + "Text" + " " + findType(n.nodeTypeInfo, useType) + xmlString
		//charField := "\t" + "Text" + " string" + xmlString
		fields = append(fields, charField)
	}
	sort.Strings(fields)
	for i := 0; i < len(fields); i++ {
		pn.lineChannel <- fields[i]
	}
}
