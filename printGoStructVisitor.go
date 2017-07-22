package main

import (
	"sort"
)

type PrintGoStructVisitor struct {
	alreadyVisited      map[string]bool
	alreadyVisitedNodes map[string]*Node
	globalTagAttributes map[string]([]*FQN)
	lineChannel         chan string
	maxDepth            int
	depth               int
	nameSpaceTagMap     map[string]string
	useType             bool
	nameSpaceInJsonName bool
	charDataFieldName   string
}

func (v *PrintGoStructVisitor) init(lineChannel chan string, maxDepth int, globalTagAttributes map[string]([]*FQN), nameSpaceTagMap map[string]string, useType bool, nameSpaceInJsonName bool) {
	v.alreadyVisited = make(map[string]bool)
	v.alreadyVisitedNodes = make(map[string]*Node)
	v.globalTagAttributes = make(map[string]([]*FQN))
	v.globalTagAttributes = globalTagAttributes
	v.lineChannel = lineChannel
	v.maxDepth = maxDepth
	v.depth = 0
	v.nameSpaceTagMap = nameSpaceTagMap
	v.useType = useType
	v.nameSpaceInJsonName = nameSpaceInJsonName
	v.charDataFieldName = "Text"
}

func (v *PrintGoStructVisitor) Visit(node *Node) bool {
	v.depth += 1

	if v.AlreadyVisited(node) {
		v.depth += 1
		return false
	}
	v.SetAlreadyVisited(node)

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
	if node.space != "" {
		v.lineChannel <- "\tXMLName  xml.Name `" + makeXmlAnnotation(node.space, false, node.name) + " " + makeJsonAnnotation(node.spaceTag, false, node.name) + "`"
	}
	v.lineChannel <- "}\n"

}

func (v *PrintGoStructVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintGoStructVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
	v.alreadyVisitedNodes[nk(n)] = n
}

func (pn *PrintGoStructVisitor) printInternalFields(n *Node) {
	var fields []string

	var field string

	for i, _ := range n.children {
		v := n.children[i]
		field = "\t" + v.makeType(namePrefix, nameSuffix) + " "
		if v.repeats {
			field += "[]*"
		} else {
			field += "*"
		}
		field += v.makeType(namePrefix, nameSuffix)

		jsonAnnotation := makeJsonAnnotation(v.spaceTag, pn.nameSpaceInJsonName, v.name)
		xmlAnnotation := makeXmlAnnotation(v.space, false, v.name)
		dbAnnotation := ""
		if addDbMetadata {
			dbAnnotation = " " + makeDbAnnotation(v.space, false, v.name)
		}

		annotation := " `" + xmlAnnotation + " " + jsonAnnotation + dbAnnotation + "`"

		field += annotation
		fields = append(fields, field)
	}

	if n.hasCharData {
		xmlString := " `xml:\",chardata\" " + makeJsonAnnotation("", false, "") + "`"
		charField := "\t" + pn.charDataFieldName + " " + findType(n.nodeTypeInfo, useType) + xmlString
		fields = append(fields, charField)
	}
	sort.Strings(fields)
	for i := 0; i < len(fields); i++ {
		pn.lineChannel <- fields[i]
	}
}

func makeJsonAnnotation(spaceTag string, useSpaceTagInName bool, name string) string {
	return makeAnnotation("json", spaceTag, false, useSpaceTagInName, name)
}

func makeXmlAnnotation(spaceTag string, useSpaceTag bool, name string) string {
	return makeAnnotation("xml", spaceTag, true, false, name)
}

func makeDbAnnotation(spaceTag string, useSpaceTag bool, name string) string {
	return makeAnnotation("db", spaceTag, true, false, name)
}

func makeAnnotation(annotationId string, spaceTag string, useSpaceTag bool, useSpaceTagInName bool, name string) (annotation string) {
	annotation = annotationId + ":\""

	if useSpaceTag {
		annotation = annotation + spaceTag
		annotation = annotation + " "
	}

	if useSpaceTagInName {
		if spaceTag != "" {
			annotation = annotation + spaceTag
		}
	}

	annotation = annotation + name + ",omitempty\""

	return annotation
}
