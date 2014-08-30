package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"text/template"
)

type PrintJavaJaxbVisitor struct {
	alreadyVisited      map[string]bool
	globalTagAttributes map[string]([]*FQN)
	nameSpaceTagMap     map[string]string
	useType             bool
	javaDir             string
	javaPackage         string
}

func (v *PrintJavaJaxbVisitor) init() {
	fullPath := v.javaDir + "/xml"
	os.RemoveAll(fullPath)
	os.MkdirAll(fullPath, 0755)
}

func (v *PrintJavaJaxbVisitor) Visit(node *Node) bool {
	if v.AlreadyVisited(node) {
		return false
	}
	v.SetAlreadyVisited(node)

	attributes := v.globalTagAttributes[nk(node)]

	class := new(JaxbClassInfo)
	class.init()
	class.PackageName = v.javaPackage
	class.ClassName = cleanName(capitalizeFirstLetter(node.name))
	class.HasValue = node.hasCharData
	class.Name = node.name

	for _, fqn := range attributes {
		jat := new(JaxbAttribute)
		cleanName := cleanName(fqn.name)
		jat.Name = fqn.name
		jat.NameUpper = capitalizeFirstLetter(cleanName)
		jat.NameLower = lowerFirstLetter(cleanName)
		jat.NameSpace = fqn.space
		class.Attributes = append(class.Attributes, jat)
	}

	for _, child := range node.children {
		jaf := new(JaxbField)
		jaf.Name = child.name
		cleanName := cleanName(child.name)
		jaf.NameUpper = capitalizeFirstLetter(cleanName)
		jaf.NameLower = lowerFirstLetter(cleanName)
		jaf.NameSpace = child.space
		jaf.Repeats = child.repeats
		jaf.TypeName = capitalizeFirstLetter(child.makeType("", ""))
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
	fi, err := os.Create(fullPath)
	if err != nil {
		log.Print("Problem creating file: " + fullPath)
		panic(err)
	}
	return bufio.NewWriter(fi), fi, nil
}
