package main

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
)

const OFloat = "FieldDefFloat"
const OString = "FieldDefString"
const OUint64 = "FieldDefUint64"
const OByte = "FieldDefByte"
const OBool = "FieldDefBool"

func findOtiraType(nti *NodeTypeInfo, useType bool) string {
	if !useType {
		return OString
	}

	if nti.alwaysBool {
		return OBool
	}

	if nti.alwaysInt08 || nti.alwaysInt16 || nti.alwaysInt32 || nti.alwaysInt64 || nti.alwaysInt0 {
		return OUint64 // TODO does not handle negative numbers//need to check ange
	}

	if nti.alwaysFloat32 || nti.alwaysFloat64 {
		return OFloat
	}
	return OString
}

type PrintOtiraVisitor struct {
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

func (v *PrintOtiraVisitor) init(writer io.Writer, maxDepth int, globalTagAttributes map[string]([]*FQN), nameSpaceTagMap map[string]string, useType bool, nameSpaceInJsonName bool) {
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

func (v *PrintOtiraVisitor) Visit(node *Node) bool {
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

func makeOtiraAttributes(writer io.Writer, attributes []*FQN, nameSpaceTagMap map[string]string) {
	sort.Sort(fqnSorter(attributes))

	for _, fqn := range attributes {

		name := fqn.name
		nameSpace := fqn.space

		fmt.Fprintln(writer, "\t ATTRIBUTES")
		fmt.Fprintln(writer, "\t tmp:= new(FieldDefString)")
		fmt.Fprintln(writer, "\t tmp.SetName(\""+name+"\")")
		//fmt.Fprintln(writer, "\n"+variableName+" "+variableType+"`xml:\""+nameSpace+name+",attr\"  json:\",omitempty\"`")

		nameSpaceTag, ok := nameSpaceTagMap[nameSpace]
		if ok && nameSpaceTag != "" {
			nameSpaceTag = nameSpaceTag + "Space"
		} else {
			nameSpaceTag = nameSpace
		}

		nameSpaceTag = goVariableNameSanitize(nameSpaceTag)
		if len(nameSpace) > 0 {
			nameSpace = nameSpace + " "
		}
	}
}

func printOtira(v *PrintOtiraVisitor, node *Node) error {
	if node == nil {
		return nil
	}
	if node.ignoredTag || node.name == "" {
		return nil
	}
	if flattenStrings && isStringOnlyField(node, len(v.globalTagAttributes[nk(node)])) {
		//v.lineChannel <- "//type " + node.makeType(namePrefix, nameSuffix)
		return nil
	}

	attributes := v.globalTagAttributes[nk(node)]
	//v.lineChannel <- "type " + node.makeType(namePrefix, nameSuffix) + " struct {"
	//fmt.Fprintln(v.writer, "type "+node.makeType(namePrefix, nameSuffix)+" struct {")

	zz := lowerFirstLetter(node.makeType("", nameSuffix))
	fmt.Fprintln(v.writer, zz+"Table, err := NewTable(\""+zz+"\")"+"//  minDepth="+strconv.Itoa(node.minDepth))
	fmt.Fprintln(v.writer, "if err != nil{")
	fmt.Fprintln(v.writer, "\treturn err")
	fmt.Fprintln(v.writer, "}")

	fmt.Fprintln(v.writer, "START ATTRIBUTES")
	makeOtiraAttributes(v.writer, attributes, v.nameSpaceTagMap)

	fmt.Fprintln(v.writer, "START INTERNAL")
	err := v.printInternalFields(len(attributes), node)
	if err != nil {
		return err
	}

	//v.lineChannel <- "}\n"
	fmt.Fprint(v.writer, "}\n\n")

	fmt.Fprint(v.writer, "**************************\n")
	return nil
}

func (v *PrintOtiraVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintOtiraVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
	v.alreadyVisitedNodes[nk(n)] = n
}

func (v *PrintOtiraVisitor) printInternalFields(nattributes int, n *Node) error {
	var fields []string

	// Fields in this struct
	for i, _ := range n.children {

		child := n.children[i]
		if child.ignoredTag {
			continue
		}
		fields = append(fields, "\n INTERNAL----- "+child.name+"   maxNumInstances: "+strconv.Itoa(child.maxNumInstances)+" :"+findOtiraType(child.nodeTypeInfo, useType))

		var def FieldDef
		if flattenStrings && isStringOnlyField(child, len(v.globalTagAttributes[nk(child)])) {
			def.GoName = child.makeType(namePrefix, nameSuffix)

			def.GoType = findType(child.nodeTypeInfo, useType)
			def.XMLName = child.name
			def.XMLNameSpace = child.space
			fields = append(fields, child.name+" simpleField "+strconv.FormatInt(child.nodeTypeInfo.maxLength, 10))
		} else {
			fields = append(fields, child.name+" complexField")
			log.Println(child.name + " complexField: ")
			log.Println(v.globalTagAttributes[nk(child)])
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
		thisType := findType(n.nodeTypeInfo, useType)
		thisVariableName := findFieldNameFromTypeInfo(thisType)

		if flattenStrings {
			//charField += "// maxLength=" + strconv.FormatInt(n.nodeTypeInfo.maxLength, 10)
			if len(n.children) == 0 && nattributes == 0 {
				//charField += "// *******************"
			}
		}
		fields = append(fields, "===== "+thisVariableName)

	}

	sort.Strings(fields)
	for i := 0; i < len(fields); i++ {
		//v.lineChannel <- fields[i]
		fmt.Fprintln(v.writer, fields[i])
	}
	return nil
}
