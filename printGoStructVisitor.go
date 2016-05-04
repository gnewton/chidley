package main

import (
	"sort"
)

type PrintGoStructVisitor struct {
	alreadyVisited      map[string]bool
	globalTagAttributes map[string]([]*FQN)
	lineChannel         chan string
	maxDepth            int
	depth               int
	nameSpaceTagMap     map[string]string
	useType             bool
	nameSpaceInJsonName bool
}

func (v *PrintGoStructVisitor) init(lineChannel chan string, maxDepth int, globalTagAttributes map[string]([]*FQN), nameSpaceTagMap map[string]string, useType bool, nameSpaceInJsonName bool) {
	v.alreadyVisited = make(map[string]bool)
	v.globalTagAttributes = make(map[string]([]*FQN))
	v.globalTagAttributes = globalTagAttributes
	v.lineChannel = lineChannel
	v.maxDepth = maxDepth
	v.depth = 0
	v.nameSpaceTagMap = nameSpaceTagMap
	v.useType = useType
	v.nameSpaceInJsonName = nameSpaceInJsonName
}

func (v *PrintGoStructVisitor) Visit(node *Node) bool {
	v.depth += 1

	if v.AlreadyVisited(node) {
		v.depth += 1
		return false
	}
	v.SetAlreadyVisited(node)

	print(v, node)

	for _, child := range node.children {
		v.Visit(child)
	}
	v.depth += 1
	return true
}

func print(v *PrintGoStructVisitor, node *Node) {
	attributes := v.globalTagAttributes[nk(node)]
	v.lineChannel <- "type " + node.makeType(namePrefix, nameSuffix) + " struct {"
	makeAttributes(v.lineChannel, attributes, v.nameSpaceTagMap)
	v.printInternalFields(node)
	if node.nameSpace != "" {
		v.lineChannel <- "\tXMLName  xml.Name `" + makeXmlAnnotation(node.nameSpace, false, node.name) + " " + makeJsonAnnotation(node.nameSpaceTag, false, node.name) + "`"
	}
	v.lineChannel <- "}\n"

}

func (v *PrintGoStructVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintGoStructVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
}

func (pn *PrintGoStructVisitor) printInternalFields(n *Node) {
	var fields []string

	var field string

	for _, v := range n.children {
		field = "\t" + v.makeType(namePrefix, nameSuffix) + " "
		if v.repeats {
			field += "[]*"
		} else {
			field += "*"
		}
		field += v.makeType(namePrefix, nameSuffix)

		jsonAnnotation := makeJsonAnnotation(v.nameSpaceTag, pn.nameSpaceInJsonName, v.name)
		xmlAnnotation := makeXmlAnnotation(v.nameSpace, false, v.name)
		dbAnnotation := ""
		if addDbMetadata {
			dbAnnotation = " " + makeDbAnnotation(v.nameSpace, false, v.name)
		}

		annotation := " `" + xmlAnnotation + " " + jsonAnnotation + dbAnnotation + "`"

		field += annotation
		fields = append(fields, field)
	}

	if n.hasCharData {
		xmlString := " `xml:\",chardata\" " + makeJsonAnnotation("", false, "") + "`"
		charField := "\t" + "Text" + " " + findType(n.nodeTypeInfo, useType) + xmlString
		fields = append(fields, charField)
	}
	sort.Strings(fields)
	for i := 0; i < len(fields); i++ {
		pn.lineChannel <- fields[i]
	}
}

func makeJsonAnnotation(nameSpaceTag string, useNameSpaceTagInName bool, name string) string {
	return makeAnnotation("json", nameSpaceTag, false, useNameSpaceTagInName, name)
}

func makeXmlAnnotation(nameSpaceTag string, useNameSpaceTag bool, name string) string {
	return makeAnnotation("xml", nameSpaceTag, true, false, name)
}

func makeDbAnnotation(nameSpaceTag string, useNameSpaceTag bool, name string) string {
	return makeAnnotation("db", nameSpaceTag, true, false, name)
}

func makeAnnotation(annotationId string, nameSpaceTag string, useNameSpaceTag bool, useNameSpaceTagInName bool, name string) (annotation string) {
	annotation = annotationId + ":\""

	if useNameSpaceTag {
		annotation = annotation + nameSpaceTag
		if nameSpaceTag != ""{
		    annotation = annotation + " "
		    }
	}

	if useNameSpaceTagInName {
		if nameSpaceTag != "" {
			annotation = annotation + nameSpaceTag + "__"
		}
	}

	annotation = annotation + name + ",omitempty\""

	return annotation
}
