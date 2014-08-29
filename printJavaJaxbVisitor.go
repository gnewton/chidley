package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
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
	fullPath := v.javaDir + "/" + strings.Replace(javaPackage, ".", "/", -1)
	os.RemoveAll(fullPath)
	os.MkdirAll(fullPath, 0755)
}

func (v *PrintJavaJaxbVisitor) Visit(node *Node) bool {
	if v.AlreadyVisited(node) {
		return false
	}
	v.SetAlreadyVisited(node)

	attributes := v.globalTagAttributes[nk(node)]

	jb := new(JaxbInfo)
	jb.init()
	jb.PackageName = v.javaPackage
	jb.ClassName = cleanName(capitalizeFirstLetter(node.name))
	jb.HasValue = node.hasCharData
	jb.Name = node.name

	for _, fqn := range attributes {
		jat := new(JaxbAttribute)
		cleanName := cleanName(fqn.name)
		jat.Name = fqn.name
		jat.NameUpper = capitalizeFirstLetter(cleanName)
		jat.NameLower = lowerFirstLetter(cleanName)
		jat.NameSpace = fqn.space
		jb.Attributes = append(jb.Attributes, jat)
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
		jb.Fields = append(jb.Fields, jaf)

	}

	printJaxbClass(jb, v.javaDir)

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

func printJaxbClass(jb *JaxbInfo, dir string) {
	t := template.Must(template.New("chidleyJaxbGen").Parse(jaxbTemplate))
	//err := t.Execute(os.Stdout, jb)
	writer, f, err := javaClassWriter(dir, jb.PackageName, jb.ClassName)
	defer f.Close()
	err = t.Execute(writer, jb)
	if err != nil {
		log.Println("executing template:", err)
	}
	bufio.NewWriter(writer).Flush()
}

func javaClassWriter(dir string, packageName string, className string) (io.Writer, *os.File, error) {
	fullPath := dir + "/" + strings.Replace(packageName, ".", "/", -1) + "/" + className + ".java"
	fi, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	return bufio.NewWriter(fi), fi, nil
}
