package main

import (
	"sort"
	"strconv"
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
	if flattenStrings && isStringOnlyField(node, len(v.globalTagAttributes[nk(node)])) {
		//v.lineChannel <- "//type " + node.makeType(namePrefix, nameSuffix)
		return
	}

	attributes := v.globalTagAttributes[nk(node)]
	v.lineChannel <- "type " + node.makeType(namePrefix, nameSuffix) + " struct {"
	makeAttributes(v.lineChannel, attributes, v.nameSpaceTagMap)
	v.printInternalFields(len(attributes), node)
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

func (v *PrintGoStructVisitor) printInternalFields(nattributes int, n *Node) {
	var fields []string

	var field string

	for i, _ := range n.children {
		child := n.children[i]
		var def OutVariableDef
		if flattenStrings && isStringOnlyField(child, len(v.globalTagAttributes[nk(child)])) {
			field = "\t" + child.makeType(namePrefix, nameSuffix) + " string `" + makeXmlAnnotation(child.space, false, child.name) + "`" + "   // ********* " + lengthTagName + ":\"" + lengthTagAttribute + lengthTagSeparator + strconv.FormatInt(child.nodeTypeInfo.maxLength+lengthTagPadding, 10) + "\""
			def.GoName = child.makeType(namePrefix, nameSuffix)
			def.GoType = "string"
			def.XMLName = child.name
			def.XMLNameSpace = child.space
		} else {

			// Field name and type are the same: i.e. Person *Person or Persons []Persons
			nameAndType := child.makeType(namePrefix, nameSuffix)

			def.GoName = nameAndType
			def.GoType = nameAndType
			def.XMLName = child.name
			def.XMLNameSpace = child.space

			field = "\t" + nameAndType + " "
			if child.repeats {
				field += "[]*"
			} else {
				field += "*"
			}
			field += nameAndType

			jsonAnnotation := makeJsonAnnotation(child.spaceTag, v.nameSpaceInJsonName, child.name)
			xmlAnnotation := makeXmlAnnotation(child.space, false, child.name)
			dbAnnotation := ""
			if addDbMetadata {
				dbAnnotation = " " + makeDbAnnotation(child.space, false, child.name)
			}

			annotation := " `" + xmlAnnotation + " " + jsonAnnotation + dbAnnotation + "`"

			field += annotation
			if flattenStrings {
				field += "   // maxLength=" + strconv.FormatInt(child.nodeTypeInfo.maxLength, 10)
			}
		}
		fields = append(fields, field)
	}

	if n.hasCharData {
		xmlString := " `xml:\",chardata\" " + makeJsonAnnotation("", false, "") + "`"
		charField := "\t" + "Text" + " " + findType(n.nodeTypeInfo, useType) + xmlString

		if flattenStrings {
			charField += "// maxLength=" + strconv.FormatInt(n.nodeTypeInfo.maxLength, 10)
			if len(n.children) == 0 && nattributes == 0 {
				charField += "// *******************"
			}
		}
		fields = append(fields, charField)
	}
	sort.Strings(fields)
	for i := 0; i < len(fields); i++ {
		v.lineChannel <- fields[i]
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
