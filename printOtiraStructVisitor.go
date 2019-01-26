package main

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"text/template"
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

func makeOtiraAttributes(writer io.Writer, tableName string, attributes []*FQN, nameSpaceTagMap map[string]string, fieldCounter *int, template *template.Template, namePrefix string) {
	sort.Sort(fqnSorter(attributes))

	for _, fqn := range attributes {
		fmt.Fprintln(writer, "\t// Attribute")
		printStringField(writer, tableName, fqn.space, fqn.name, 33, fieldCounter, template, namePrefix, "Attribute")
	}

}

type TStringField struct {
	FieldVariableName, FieldName, TableVariableName string
	FieldLength                                     int
	Comment                                         string
}

func printStringField(writer io.Writer, tableName, nameSpace, localName string, length int, fieldCounter *int, template *template.Template, namePrefix string, comment string) {
	var name string
	if nameSpace == "" {
		name = localName
	} else {
		name = nameSpace + "__" + localName
	}
	//name =
	if namePrefix != "" {
		namePrefix = namePrefix + "_"
	}
	name = namePrefix + name
	sqlName := sqlizeString(name)

	name = "v" + name + strconv.Itoa(*fieldCounter)
	*fieldCounter = *fieldCounter + 1

	data := TStringField{name, sqlName, tableName, length, comment}

	err := template.Execute(writer, data)
	if err != nil {
		log.Fatal("executing template:", err)
	}
}

func printOtiraTables(v *PrintOtiraVisitor, node *Node, fieldCounter *int, template *template.Template, collapsedXmlTagsList []string) error {
	if node == nil {
		return nil
	}
	if node.ignoredTag || node.name == "" {
		return nil
	}

	if flattenStrings && isStringOnlyField(node, len(v.globalTagAttributes[nk(node)])) {
		return nil
	}

	if contains(collapsedXmlTagsList, node.name) {
		return nil
	}

	zz := lowerFirstLetter(node.makeType("", nameSuffix))
	if zz == "" {
		return nil
	}

	tableName := zz + "Table"

	fmt.Fprintln(v.writer, "\n\t// Table "+tableName)
	fmt.Fprintln(v.writer, "\t//")
	fmt.Fprintln(v.writer, "\t"+tableName+", err := otira.NewTableDef(\""+sqlizeString(zz)+"\")")
	fmt.Fprintln(v.writer, "\tif err != nil{")
	fmt.Fprintln(v.writer, "\t\treturn")
	fmt.Fprintln(v.writer, "\t}")

	makePrimaryKey(v.writer, zz)

	// Start anoymous function
	attributes := v.globalTagAttributes[nk(node)]
	makeOtiraAttributes(v.writer, zz, attributes, v.nameSpaceTagMap, fieldCounter, template, "")

	if node.hasCharData {
		printStringField(v.writer, tableName, node.space, zz, int(node.nodeTypeInfo.maxLength), fieldCounter, template, "", "CharContent")
	}
	err := v.printInternalFields(zz, len(attributes), node, fieldCounter, template)
	if err != nil {
		return err
	}

	return nil
}

func makePrimaryKey(w io.Writer, tablename string) {
	printUint64Field(w, tablename)
}

func printUint64Field(w io.Writer, tablename string) {
	fmt.Fprintln(w, "pk = new(otira.FieldDefUint64)")
	fmt.Fprintln(w, "pk.SetName(\"id\")")
	fmt.Fprintln(w, tablename+"Table.Add(pk)")
}

func printOtiraRelations(v *PrintOtiraVisitor, node *Node, one2mT, m2mT *template.Template, collapsedXmlTagsList []string) error {
	zz := lowerFirstLetter(node.makeType("", nameSuffix))

	err := v.printInternalRelationFields(zz, node, one2mT, m2mT, collapsedXmlTagsList)
	if err != nil {
		return err
	}
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

func (v *PrintOtiraVisitor) printInternalRelationFields(tableName string, n *Node, one2mT, m2mT *template.Template, collapsedXmlTagsList []string) error {
	if n.name == "" {
		return nil
	}
	m2mCounter, one2mCounter := 0, 0
	// Fields in this struct
	for i, _ := range n.children {

		child := n.children[i]
		if child.ignoredTag || contains(collapsedXmlTagsList, child.name) {
			continue
		}

		if flattenStrings && isStringOnlyField(child, len(v.globalTagAttributes[nk(child)])) {

		} else {
			fmt.Fprintln(v.writer, "\t//"+child.name+" complexField")
			if child.repeats {
				fmt.Fprintln(v.writer, "\t//"+child.name+":repeats "+strconv.Itoa(child.maxNumInstances))
				printManyToMany(v.writer, n.name, child.name, m2mCounter, m2mT)
				m2mCounter += 1
			} else {
				printOneToMany(v.writer, n.name, child.name, one2mCounter, one2mT)
				fmt.Fprintln(v.writer, "\t //"+child.name+":single "+strconv.Itoa(child.maxNumInstances))
				one2mCounter += 1
			}
		}
	}
	return nil
}

func (v *PrintOtiraVisitor) printInternalFields(tableName string, nattributes int, n *Node, fieldCounter *int, template *template.Template) error {

	// Fields in this struct
	for i, _ := range n.children {

		child := n.children[i]
		if child.ignoredTag {
			continue
		}

		// Collapsed tags
		if contains(collapsedXmlTagsList, child.name) {
			fmt.Fprintln(v.writer, "\t //Collapsed Field")

			printStringField(v.writer, tableName, child.space, child.name, int(child.nodeTypeInfo.maxLength), fieldCounter, template, "", "Collapsed field")
			v.printInternalFields(tableName, nattributes, child, fieldCounter, template)
			attributes := v.globalTagAttributes[nk(child)]
			zz := lowerFirstLetter(child.makeType("", nameSuffix))
			if zz == "" {
				return nil
			}

			makeOtiraAttributes(v.writer, tableName, attributes, v.nameSpaceTagMap, fieldCounter, template, child.name)
		} else {
			if flattenStrings && isStringOnlyField(child, len(v.globalTagAttributes[nk(child)])) {
				printStringField(v.writer, tableName, child.space, child.name, int(child.nodeTypeInfo.maxLength), fieldCounter, template, "", "DOES THIS HAPPEN")
			}
		}
	}
	return nil
}

func printManyToMany(w io.Writer, left, right string, counter int, m2mT *template.Template) {
	l := lowerFirstLetter(left)
	r := lowerFirstLetter(right)

	//t := template.Must(template.New(many2manyTemplateName).Parse(many2manyTemplate))
	data := TRelations{r, l, "", counter}
	err := m2mT.Execute(w, data)
	if err != nil {
		log.Fatal("executing template:", err)
	}

}

type TRelations struct {
	Right, Left, RightSql string
	Counter               int
}

func printOneToMany(w io.Writer, left, right string, counter int, one2mT *template.Template) {
	l := lowerFirstLetter(left)
	r := lowerFirstLetter(right)

	//t := template.Must(template.New(one2manyTemplateName).Parse(one2manyTemplate))
	data := TRelations{r, l, sqlizeString(right), counter}
	err := one2mT.Execute(w, data)
	if err != nil {
		log.Fatal("executing template:", err)
	}
}

func printRelation(w io.Writer, left, right string, counter int, template *template.Template) {
	l := lowerFirstLetter(left)
	r := lowerFirstLetter(right)

	data := TRelations{r, l, "", counter}
	err := template.Execute(w, data)
	if err != nil {
		log.Fatal("executing template:", err)
	}

}
