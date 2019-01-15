package main

import (
	"fmt"
	"io"
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

func makeOtiraAttributes(writer io.Writer, tableName string, attributes []*FQN, nameSpaceTagMap map[string]string) {
	sort.Sort(fqnSorter(attributes))

	for _, fqn := range attributes {
		// Because golang/xml does not properly handle name spaces....
		//if fqn.space == "xmlns" || fqn.name == "xmlns" {
		//continue
		//}
		fmt.Fprintln(writer, "\t ATTRIBUTES")
		// Attributes have no namespaces
		printStringField(writer, tableName, fqn.space, fqn.name, 33)
	}
}

func printStringField(writer io.Writer, tableName, nameSpace, localName string, length int) {
	var name string
	if nameSpace == "" {
		name = localName
	} else {
		name = nameSpace + "__" + localName
	}
	name = lowerFirstLetter(name)
	fmt.Fprintln(writer, "\t "+name+" := new(FieldDefString)"+"  // *****************")
	fmt.Fprintln(writer, "\t "+name+".SetName(\""+name+"\")")
	if length > 0 {
		fmt.Fprintln(writer, "\t "+name+".SetLength(\""+strconv.Itoa(length)+"\")")
	}
	fmt.Fprintln(writer, "\t "+tableName+".Add(\""+name+"\")")
}

func printOtiraNode(v *PrintOtiraVisitor, node *Node) error {
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

	// Start anoymous function
	fmt.Fprintln(v.writer, "func(){")
	fmt.Fprintln(v.writer, "START ATTRIBUTES")
	makeOtiraAttributes(v.writer, zz, attributes, v.nameSpaceTagMap)

	fmt.Fprintln(v.writer, "START INTERNAL")
	err := v.printInternalFields(zz, len(attributes), node)
	if err != nil {
		return err
	}

	// END Anonymous function
	fmt.Fprintln(v.writer, "}()")

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

func (v *PrintOtiraVisitor) printInternalFields(tableName string, nattributes int, n *Node) error {
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
			printStringField(v.writer, tableName, child.space, child.name, int(child.nodeTypeInfo.maxLength))
		} else {
			fields = append(fields, child.name+" complexField")
			// Field name and type are the same: i.e. Person *Person or Persons []Persons
			nameAndType := child.makeType(namePrefix, nameSuffix)

			def.GoName = nameAndType
			def.GoType = nameAndType
			def.XMLName = child.name
			def.XMLNameSpace = child.space

			if child.repeats {
				fmt.Fprintln(v.writer, "\t "+child.name+":repeats "+strconv.Itoa(child.maxNumInstances))
				printManyToMany(v.writer, n.name, child.name)
			} else {
				printOneToMany(v.writer, n.name, child.name)
				fmt.Fprintln(v.writer, "\t "+child.name+":single "+strconv.Itoa(child.maxNumInstances))
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

func printManyToMany(w io.Writer, left, right string) {
	l := lowerFirstLetter(left)
	r := lowerFirstLetter(right)
	fmt.Fprintln(w, "m2m := NewManyToMany()")
	fmt.Fprintln(w, "m2m.leftTable = "+l+"Table")
	fmt.Fprintln(w, "m2m.rightTable = "+r+"Table")
	fmt.Fprintln(w, l+"Table.AddManyToMany(m2m)")
}

func printOneToMany(w io.Writer, left, right string) {
	l := lowerFirstLetter(left)
	r := lowerFirstLetter(right)
	fmt.Fprintln(w, "one2m := NewOneToMany()")
	fmt.Fprintln(w, "one2m.leftTable = "+l+"Table")
	fmt.Fprintln(w, "one2m.rightTable = "+r+"Table")
	fmt.Fprintln(w, "leftField := new(FieldDefUint64)")
	fmt.Fprintln(w, "leftField.SetName(\""+r+"\")")
	fmt.Fprintln(w, l+"Table.Add(leftField)")

	fmt.Fprintln(w, "one2m.leftKeyField = leftField")
	fmt.Fprintln(w, "one2m.rightKeyField = "+r+"Table.PrimaryKey()")
	fmt.Fprintln(w, l+"Table.AddOneToMany(one2m)")
}
