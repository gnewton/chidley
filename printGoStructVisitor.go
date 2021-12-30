package main

import (
	"fmt"
	"io"
	"sort"
	//"strconv"
)

type PrintGoStructVisitor struct {
	alreadyVisited      map[string]bool
	alreadyVisitedNodes map[string]*Node
	globalTagAttributes map[string]([]*FQN)

	maxDepth            int
	depth               int
	nameSpaceTagMap     map[string]string
	useType             bool
	nameSpaceInJsonName bool
	writer              io.Writer
}

func (v *PrintGoStructVisitor) init(writer io.Writer, maxDepth int, globalTagAttributes map[string]([]*FQN), nameSpaceTagMap map[string]string, useType bool, nameSpaceInJsonName bool) {
	v.alreadyVisited = make(map[string]bool)
	v.alreadyVisitedNodes = make(map[string]*Node)
	v.globalTagAttributes = make(map[string]([]*FQN))
	v.globalTagAttributes = globalTagAttributes
	v.writer = writer
	v.maxDepth = maxDepth
	v.depth = 0
	v.nameSpaceTagMap = nameSpaceTagMap
	v.useType = useType
	v.nameSpaceInJsonName = nameSpaceInJsonName
}

func (v *PrintGoStructVisitor) Visit(node *Node) bool {

	v.depth += 1

	if v.AlreadyVisited(node) || node.ignoredTag {
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

func print(v *PrintGoStructVisitor, node *Node) error {
	if node.ignoredTag || node.name == "" {
		return nil
	}
	if flattenStrings && isStringOnlyField(node, len(v.globalTagAttributes[nk(node)])) {
		//v.lineChannel <- "//type " + node.makeType(namePrefix, nameSuffix)
		return nil
	}

	attributes := v.globalTagAttributes[nk(node)]
	//v.lineChannel <- "type " + node.makeType(namePrefix, nameSuffix) + " struct {"
	fmt.Fprintln(v.writer, "type "+node.makeType(namePrefix, nameSuffix)+" struct {")

	//	fmt.Fprintln(v.writer, "\tXMLName xml.Name`"+makeXmlAnnotation(node.space, false, node.name)+" "+makeJsonAnnotation(node.spaceTag, false, node.name)+"`")

	fmt.Fprintln(v.writer, "\tXMLName xml.Name `"+makeAnnotation("xml", node.space, false, false, node.name)+" "+makeJsonAnnotation(node.spaceTag, false, node.name)+"`")

	//return makeAnnotation("xml", spaceTag, true, false, name)

	makeAttributes(v.writer, attributes, v.nameSpaceTagMap)

	err := v.printInternalFields(len(attributes), node)
	if err != nil {
		return err
	}

	//v.lineChannel <- "}\n"
	fmt.Fprintln(v.writer, "}")
	fmt.Fprintln(v.writer, "")

	return nil
}

func (v *PrintGoStructVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintGoStructVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
	v.alreadyVisitedNodes[nk(n)] = n
}

func (v *PrintGoStructVisitor) printInternalFields(nattributes int, n *Node) error {
	var fields []string

	// Fields in this struct
	for i, _ := range n.children {
		child := n.children[i]
		if child.ignoredTag {
			continue
		}
		var def FieldDef
		if flattenStrings && isStringOnlyField(child, len(v.globalTagAttributes[nk(child)])) {
			//field = "\t" + child.spaceTag + child.makeType(namePrefix, nameSuffix) + " string `" + makeXmlAnnotation(child.space, false, child.name) + "`" //+ "   // ********* " + lengthTagName + ":\"" + lengthTagAttribute + lengthTagSeparator + strconv.FormatInt(child.nodeTypeInfo.maxLength+lengthTagPadding, 10) + "\""
			def.GoName = child.makeType(namePrefix, nameSuffix)
			//def.GoType = "string"
			def.GoType = findType(child.nodeTypeInfo, useType)
			def.XMLName = child.name
			def.XMLNameSpace = child.space
		} else {

			// Field name and type are the same: i.e. Person *Person or Persons []Persons
			nameAndType := child.makeType(namePrefix, nameSuffix)

			def.GoName = nameAndType
			def.GoType = nameAndType
			def.XMLName = child.name
			def.XMLNameSpace = child.space

			if child.repeats {
				def.GoTypeArrayOrPointer = "[]*"
			} else {
				def.GoTypeArrayOrPointer = "*"
			}
		}
		if flattenStrings {
			def.Length = child.nodeTypeInfo.maxLength
		}
		fieldDefString, err := render(def)
		if err != nil {
			return err
		}
		fields = append(fields, fieldDefString)
	}

	// Is this chardata Field (string)
	if n.hasCharData {
		xmlString := " `xml:\",chardata\" " + makeJsonAnnotation("", false, "") + "`"
		thisType := findType(n.nodeTypeInfo, useType)
		thisVariableName := findFieldNameFromTypeInfo(thisType)

		charField := "\t" + thisVariableName + " " + thisType + xmlString

		if flattenStrings {
			//charField += "// maxLength=" + strconv.FormatInt(n.nodeTypeInfo.maxLength, 10)
			if len(n.children) == 0 && nattributes == 0 {
				charField += "// *******************"
			}
		}
		//GOOD
		//charField += "   // maxLength=" + strconv.FormatInt(n.nodeTypeInfo.maxLength, 10)

		fields = append(fields, charField)
	}

	sort.Strings(fields)
	for i := 0; i < len(fields); i++ {
		//v.lineChannel <- fields[i]
		fmt.Fprintln(v.writer, fields[i])
	}
	return nil
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
