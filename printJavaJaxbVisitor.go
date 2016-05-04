package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"text/template"
	"time"
)

type PrintJavaJaxbVisitor struct {
	alreadyVisited      map[string]bool
	globalTagAttributes map[string]([]*FQN)
	nameSpaceTagMap     map[string]string
	useType             bool
	javaDir             string
	javaPackage         string
	namePrefix          string
	Date                time.Time
}

func (v *PrintJavaJaxbVisitor) Visit(node *Node) bool {
	if v.AlreadyVisited(node) {
		return false
	}
	v.SetAlreadyVisited(node)

	attributes := v.globalTagAttributes[nk(node)]

	class := new(JaxbClassInfo)
	class.init()
	class.Date = v.Date
	class.PackageName = v.javaPackage
	class.ClassName = v.namePrefix + cleanName(capitalizeFirstLetter(node.name))
	class.HasValue = node.hasCharData
	class.ValueType = findJavaType(node.nodeTypeInfo, v.useType)
	class.Name = node.name

	for _, fqn := range attributes {
		jat := new(JaxbAttribute)
		cleanName := cleanName(fqn.name)
		jat.Name = fqn.name
		jat.NameUpper = capitalizeFirstLetter(cleanName)
		if v.namePrefix != "" {
			jat.NameLower = lowerFirstLetter(v.namePrefix) + capitalizeFirstLetter(cleanName)
		} else {
			jat.NameLower = lowerFirstLetter(cleanName)
		}
		jat.NameSpace = fqn.nameSpace
		class.Attributes = append(class.Attributes, jat)
	}

	for _, child := range node.children {
		jaf := new(JaxbField)
		jaf.Name = child.name
		cleanName := cleanName(child.name)
		jaf.NameUpper = capitalizeFirstLetter(cleanName)
		if v.namePrefix != "" {
			jaf.NameLower = lowerFirstLetter(v.namePrefix) + capitalizeFirstLetter(cleanName)
		} else {
			jaf.NameLower = lowerFirstLetter(cleanName)
		}
		jaf.NameSpace = child.nameSpace
		jaf.Repeats = child.repeats
		jaf.TypeName = child.makeJavaType(v.namePrefix, "")
		class.Fields = append(class.Fields, jaf)

	}

	printJaxbClass(class, v.javaDir+"/xml")

	for _, child := range node.children {
		v.Visit(child)
	}

	return true
}

func (v *PrintJavaJaxbVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *PrintJavaJaxbVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
}

func printJaxbClass(class *JaxbClassInfo, dir string) {
	t := template.Must(template.New("chidleyJaxbGen").Parse(jaxbClassTemplate))
	//err := t.Execute(os.Stdout, jb)
	writer, f, err := javaClassWriter(dir, class.PackageName+".xml", class.ClassName)
	defer f.Close()
	err = t.Execute(writer, class)
	if err != nil {
		log.Println("executing template:", err)
	}
	bufio.NewWriter(writer).Flush()
}

func javaClassWriter(dir string, packageName string, className string) (io.Writer, *os.File, error) {
	fullPath := dir + "/" + className + ".java"
	log.Print("Writing java Class file: " + fullPath)
	fi, err := os.Create(fullPath)
	if err != nil {
		log.Print("Problem creating file: " + fullPath)
		panic(err)
	}
	return bufio.NewWriter(fi), fi, nil
}
