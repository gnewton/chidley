package main

// Copyright 2014,2015,2016 Glen Newton
// glen.newton@gmail.com

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"time"
)

var DEBUG = false
var attributePrefix = "Attr"
var codeGenConvert = false
var classicStructNamesWithUnderscores = false
var nameSpaceInJsonName = false
var prettyPrint = false
var progress = false
var readFromStandardIn = false
var sortByXmlOrder = false
var structsToStdout = true

var codeGenDir = "codegen"
var codeGenFilename = "CodeGenStructs.go"

// Java out
const javaBasePackage = "ca.gnewton.chidley"
const mavenJavaBase = "src/main/java"

var javaBasePackagePath = strings.Replace(javaBasePackage, ".", "/", -1)
var javaAppName = "jaxb"
var writeJava = false
var baseJavaDir = "java"
var userJavaPackageName = ""
var charDataFieldName = ""

var namePrefix = "Chi"
var nameSuffix = ""
var xmlName = false
var url = false
var useType = false
var addDbMetadata = false

type structSortFunc func(v *PrintGoStructVisitor)

var structSort = printStructsAlphabetical

type Writer interface {
	open(s string, lineChannel chan string) error
	close()
}

var outputs = []*bool{
	&codeGenConvert,
	&structsToStdout,
	&writeJava,
}

func init() {

	flag.BoolVar(&classicStructNamesWithUnderscores, "C", classicStructNamesWithUnderscores, "Structs have underscores instead of CamelCase; how chidley used to produce output; includes name spaces (see -n)")

	flag.BoolVar(&DEBUG, "d", DEBUG, "Debug; prints out much information")
	flag.BoolVar(&addDbMetadata, "B", addDbMetadata, "Add database metadata to created Go structs")
	flag.BoolVar(&sortByXmlOrder, "X", sortByXmlOrder, "Sort output of structs in Go code by order encounered in source XML (default is alphabetical order)")
	flag.BoolVar(&codeGenConvert, "W", codeGenConvert, "Generate Go code to convert XML to JSON or XML (latter useful for validation) and write it to stdout")
	flag.BoolVar(&nameSpaceInJsonName, "n", nameSpaceInJsonName, "Use the XML namespace prefix as prefix to JSON name")
	flag.BoolVar(&prettyPrint, "p", prettyPrint, "Pretty-print json in generated code (if applicable)")
	flag.BoolVar(&progress, "r", progress, "Progress: every 50000 input tags (elements)")
	flag.BoolVar(&readFromStandardIn, "c", readFromStandardIn, "Read XML from standard input")
	flag.BoolVar(&structsToStdout, "G", structsToStdout, "Only write generated Go structs to stdout")
	flag.BoolVar(&url, "u", url, "Filename interpreted as an URL")
	flag.BoolVar(&useType, "t", useType, "Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete")
	flag.BoolVar(&writeJava, "J", writeJava, "Generated Java code for Java/JAXB")
	flag.BoolVar(&xmlName, "x", xmlName, "Add XMLName (Space, Local) for each XML element, to JSON")
	flag.StringVar(&attributePrefix, "a", attributePrefix, "Prefix to attribute names")
	flag.StringVar(&baseJavaDir, "D", baseJavaDir, "Base directory for generated Java code (root of maven project)")
	flag.StringVar(&javaAppName, "k", javaAppName, "App name for Java code (appended to ca.gnewton.chidley Java package name))")
	flag.StringVar(&namePrefix, "e", namePrefix, "Prefix to struct (element) names; must start with a capital")
	flag.StringVar(&userJavaPackageName, "P", userJavaPackageName, "Java package name (rightmost in full package name")
	flag.StringVar(&charDataFieldName, "cdata", "Text", "Field name for any chardata field")
}

func handleParameters() error {
	flag.Parse()

	if codeGenConvert || writeJava {
		structsToStdout = false
	}

	numBoolsSet := countNumberOfBoolsSet(outputs)
	if numBoolsSet > 1 {
		log.Print("  ERROR: Only one of -W -J -X -V -c can be set")
	} else if numBoolsSet == 0 {
		log.Print("  ERROR: At least one of -W -J -X -V -c must be set")
	}
	if sortByXmlOrder {
		structSort = printStructsByXml
	}
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := handleParameters()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err != nil {
		flag.Usage()
		return
	}

	if len(flag.Args()) != 1 && !readFromStandardIn {
		fmt.Println("chidley <flags> xmlFileName|url")
		fmt.Println("xmlFileName can be .gz or .bz2: uncompressed transparently")
		flag.Usage()
		return
	}

	var sourceName string

	if !readFromStandardIn {
		sourceName = flag.Args()[0]
	}
	if !url && !readFromStandardIn {
		sourceName, err = filepath.Abs(sourceName)
		if err != nil {
			log.Fatal("FATAL ERROR: " + err.Error())
		}
	}

	source, err := makeSourceReader(sourceName, url, readFromStandardIn)
	if err != nil {
		log.Fatal("FATAL ERROR: " + err.Error())
	}

	ex := Extractor{
		namePrefix: namePrefix,
		nameSuffix: nameSuffix,
		reader:     source.getReader(),
		useType:    useType,
		progress:   progress,
	}

	if DEBUG {
		log.Print("extracting")
	}
	err = ex.extract()

	if err != nil {
		log.Fatal("FATAL ERROR: " + err.Error())
	}

	var writer Writer
	lineChannel := make(chan string, 100)
	switch {
	case codeGenConvert:
		sWriter := new(stringWriter)
		writer = sWriter
		writer.open("", lineChannel)
		printGoStructVisitor := new(PrintGoStructVisitor)
		printGoStructVisitor.init(lineChannel, 9999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
		printGoStructVisitor.charDataFieldName = charDataFieldName
		printGoStructVisitor.Visit(ex.root)

		structSort(printGoStructVisitor)

		close(lineChannel)
		sWriter.close()

		xt := XMLType{NameType: ex.firstNode.makeType(namePrefix, nameSuffix),
			XMLName:      ex.firstNode.name,
			XMLNameUpper: capitalizeFirstLetter(ex.firstNode.name),
			XMLSpace:     ex.firstNode.space,
		}

		x := XmlInfo{
			BaseXML:         &xt,
			OneLevelDownXML: makeOneLevelDown(ex.root),
			Filename:        getFullPath(sourceName),
			Structs:         sWriter.s,
		}
		t := template.Must(template.New("chidleyGen").Parse(codeTemplate))

		err := t.Execute(os.Stdout, x)
		if err != nil {
			log.Println("executing template:", err)
		}

	case structsToStdout:
		writer = new(stdoutWriter)
		writer.open("", lineChannel)
		printGoStructVisitor := new(PrintGoStructVisitor)
		printGoStructVisitor.init(lineChannel, 999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
		printGoStructVisitor.charDataFieldName = charDataFieldName
		printGoStructVisitor.Visit(ex.root)
		structSort(printGoStructVisitor)
		close(lineChannel)
		writer.close()

	case writeJava:
		if len(userJavaPackageName) > 0 {
			javaAppName = userJavaPackageName
		}
		javaPackage := javaBasePackage + "." + javaAppName
		javaDir := baseJavaDir + "/" + mavenJavaBase + "/" + javaBasePackagePath + "/" + javaAppName

		os.RemoveAll(baseJavaDir)
		os.MkdirAll(javaDir+"/xml", 0755)
		date := time.Now()
		printJavaJaxbVisitor := PrintJavaJaxbVisitor{
			alreadyVisited:      make(map[string]bool),
			globalTagAttributes: ex.globalTagAttributes,
			nameSpaceTagMap:     ex.nameSpaceTagMap,
			useType:             useType,
			javaDir:             javaDir,
			javaPackage:         javaPackage,
			namePrefix:          namePrefix,
			Date:                date,
		}

		var onlyChild *Node
		for _, child := range ex.root.children {
			printJavaJaxbVisitor.Visit(child)
			// Bad: assume only one base element
			onlyChild = child
		}
		printJavaJaxbMain(onlyChild.makeJavaType(namePrefix, ""), javaDir, javaPackage, getFullPath(sourceName), date)
		printPackageInfo(onlyChild, javaDir, javaPackage, ex.globalTagAttributes, ex.nameSpaceTagMap)

		printMavenPom(baseJavaDir+"/pom.xml", javaAppName)
	}

}

func printPackageInfo(node *Node, javaDir string, javaPackage string, globalTagAttributes map[string][]*FQN, nameSpaceTagMap map[string]string) {

	//log.Printf("%+v\n", node)

	if node.space != "" {
		_ = findNameSpaces(globalTagAttributes[nk(node)])
		//attributes := findNameSpaces(globalTagAttributes[nk(node)])

		t := template.Must(template.New("package-info").Parse(jaxbPackageInfoTemplage))
		packageInfoPath := javaDir + "/xml/package-info.java"
		fi, err := os.Create(packageInfoPath)
		if err != nil {
			log.Print("Problem creating file: " + packageInfoPath)
			panic(err)
		}
		defer fi.Close()

		writer := bufio.NewWriter(fi)
		packageInfo := JaxbPackageInfo{
			BaseNameSpace: node.space,
			//AdditionalNameSpace []*FQN
			PackageName: javaPackage + ".xml",
		}
		err = t.Execute(writer, packageInfo)
		if err != nil {
			log.Println("executing template:", err)
		}
		bufio.NewWriter(writer).Flush()
	}

}

const XMLNS = "xmlns"

func findNameSpaces(attributes []*FQN) []*FQN {
	if attributes == nil || len(attributes) == 0 {
		return nil
	}
	xmlns := make([]*FQN, 0)
	return xmlns
}

func printMavenPom(pomPath string, javaAppName string) {
	t := template.Must(template.New("mavenPom").Parse(mavenPomTemplate))
	fi, err := os.Create(pomPath)
	if err != nil {
		log.Print("Problem creating file: " + pomPath)
		panic(err)
	}
	defer fi.Close()

	writer := bufio.NewWriter(fi)
	maven := JaxbMavenPomInfo{
		AppName: javaAppName,
	}
	err = t.Execute(writer, maven)
	if err != nil {
		log.Println("executing template:", err)
	}
	bufio.NewWriter(writer).Flush()
}

func printJavaJaxbMain(rootElementName string, javaDir string, javaPackage string, sourceXMLFilename string, date time.Time) {
	t := template.Must(template.New("chidleyJaxbGenClass").Parse(jaxbMainTemplate))
	writer, f, err := javaClassWriter(javaDir, javaPackage, "Main")
	defer f.Close()

	classInfo := JaxbMainClassInfo{
		PackageName:       javaPackage,
		BaseXMLClassName:  rootElementName,
		SourceXMLFilename: sourceXMLFilename,
		Date:              date,
	}
	err = t.Execute(writer, classInfo)
	if err != nil {
		log.Println("executing template:", err)
	}
	bufio.NewWriter(writer).Flush()

}

func makeSourceReader(sourceName string, url bool, standardIn bool) (Source, error) {
	var err error

	var source Source
	if url {
		source = new(UrlSource)
		if DEBUG {
			log.Print("Making UrlSource")
		}
	} else {
		if standardIn {
			source = new(StdinSource)
			if DEBUG {
				log.Print("Making StdinSource")
			}
		} else {
			source = new(FileSource)
			if DEBUG {
				log.Print("Making FileSource")
			}
		}
	}
	if DEBUG {
		log.Print("Making Source:[" + sourceName + "]")
	}
	err = source.newSource(sourceName)
	return source, err
}

func attributes(atts map[string]bool) string {
	ret := ": "
	for k, _ := range atts {
		ret = ret + k + ", "
	}
	return ret
}

func indent(d int) string {
	indent := ""
	for i := 0; i < d; i++ {
		indent = indent + "\t"
	}
	return indent
}

func countNumberOfBoolsSet(a []*bool) int {
	counter := 0
	for i := 0; i < len(a); i++ {
		if *a[i] {
			counter += 1
		}
	}
	return counter
}

func makeOneLevelDown(node *Node) []*XMLType {
	var children []*XMLType

	for _, np := range node.children {
		if np == nil {
			continue
		}
		for _, n := range np.children {
			if n == nil {
				continue
			}
			x := XMLType{NameType: n.makeType(namePrefix, nameSuffix),
				XMLName:      n.name,
				XMLNameUpper: capitalizeFirstLetter(n.name),
				XMLSpace:     n.space}
			children = append(children, &x)
		}
	}
	return children
}
func printChildrenChildren(node *Node) {
	for k, v := range node.children {
		log.Print(k)
		log.Printf("children: %+v\n", v.children)
	}
}

// Order Xml is encountered
func printStructsByXml(v *PrintGoStructVisitor) {
	orderNodes := make(map[int]*Node)
	var order []int

	for k := range v.alreadyVisitedNodes {
		nodeOrder := v.alreadyVisitedNodes[k].discoveredOrder
		orderNodes[nodeOrder] = v.alreadyVisitedNodes[k]
		order = append(order, nodeOrder)
	}
	sort.Ints(order)

	for o := range order {
		print(v, orderNodes[o])
	}
}

// Alphabetical order
func printStructsAlphabetical(v *PrintGoStructVisitor) {
	var keys []string
	for k := range v.alreadyVisitedNodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		print(v, v.alreadyVisitedNodes[k])
	}

}
