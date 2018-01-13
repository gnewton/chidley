package main

// Copyright 2014,2015,2016 Glen Newton
// glen.newton@gmail.com

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"text/template"
	"time"
)

func init() {

	flag.BoolVar(&DEBUG, "d", DEBUG, "Debug; prints out much information")
	flag.BoolVar(&addDbMetadata, "B", addDbMetadata, "Add database metadata to created Go structs")
	flag.BoolVar(&classicStructNamesWithUnderscores, "C", classicStructNamesWithUnderscores, "Structs have underscores instead of CamelCase; how chidley used to produce output; includes name spaces (see -n)")
	flag.BoolVar(&codeGenConvert, "W", codeGenConvert, "Generate Go code to convert XML to JSON or XML (latter useful for validation) and write it to stdout")
	flag.BoolVar(&flattenStrings, "F", flattenStrings, "Assume complete representative XML and collapse tags with only a single string and no attributes")
	flag.BoolVar(&ignoreXmlDecodingErrors, "I", ignoreXmlDecodingErrors, "If XML decoding error encountered, continue")
	flag.BoolVar(&nameSpaceInJsonName, "n", nameSpaceInJsonName, "Use the XML namespace prefix as prefix to JSON name")
	flag.BoolVar(&prettyPrint, "p", prettyPrint, "Pretty-print json in generated code (if applicable)")
	flag.BoolVar(&progress, "r", progress, "Progress: every 50000 input tags (elements)")
	flag.BoolVar(&readFromStandardIn, "c", readFromStandardIn, "Read XML from standard input")
	flag.BoolVar(&sortByXmlOrder, "X", sortByXmlOrder, "Sort output of structs in Go code by order encounered in source XML (default is alphabetical order)")
	flag.BoolVar(&structsToStdout, "G", structsToStdout, "Only write generated Go structs to stdout")
	flag.BoolVar(&url, "u", url, "Filename interpreted as an URL")
	flag.BoolVar(&useType, "t", useType, "Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete")
	flag.BoolVar(&writeJava, "J", writeJava, "Generated Java code for Java/JAXB")
	flag.BoolVar(&xmlName, "x", xmlName, "Add XMLName (Space, Local) for each XML element, to JSON")

	flag.BoolVar(&keepXmlFirstLetterCase, "K", keepXmlFirstLetterCase, "Do not change the case of the first letter of the XML tag names")

	flag.StringVar(&attributePrefix, "a", attributePrefix, "Prefix to attribute names")
	flag.StringVar(&baseJavaDir, "D", baseJavaDir, "Base directory for generated Java code (root of maven project)")
	flag.StringVar(&javaAppName, "k", javaAppName, "App name for Java code (appended to ca.gnewton.chidley Java package name))")
	flag.StringVar(&namePrefix, "e", namePrefix, "Prefix to struct (element) names; must start with a capital")
	flag.StringVar(&userJavaPackageName, "P", userJavaPackageName, "Java package name (rightmost in full package name")

	flag.StringVar(&lengthTagName, "N", lengthTagName, "The tag name to use for the max length Go annotations")
	flag.StringVar(&lengthTagAttribute, "A", lengthTagAttribute, "The tag name attribute to use for the max length Go annotations")
	flag.StringVar(&lengthTagSeparator, "S", lengthTagSeparator, "The tag name separator to use for the max length Go annotations")
	flag.Int64Var(&lengthTagPadding, "Z", lengthTagPadding, "The padding on the max length tag attribute")

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

	if lengthTagName == "" && lengthTagAttribute == "" || lengthTagName != "" && lengthTagAttribute != "" {
		return nil
	}

	return errors.New("Both lengthTagName and lengthTagAttribute must be set")
}

func main() {

	//EXP()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := handleParameters()

	if err != nil {
		log.Println(err)
		flag.Usage()
		return
	}

	if len(flag.Args()) == 1 && !readFromStandardIn {
		fmt.Println("chidley <flags> xmlFileName|url")
		fmt.Println("xmlFileName can be .gz or .bz2: uncompressed transparently")
		flag.Usage()
		return
	}

	var sourceNames []string

	if !readFromStandardIn {
		sourceNames = flag.Args()
		log.Println("FLAGS", flag.Args())
	}
	if !url && !readFromStandardIn {
		for i, _ := range sourceNames {
			sourceNames[i], err = filepath.Abs(sourceNames[i])
			if err != nil {
				log.Fatal("FATAL ERROR: " + err.Error())
			}
		}
	}

	sources, err := makeSourceReaders(sourceNames, url, readFromStandardIn)
	if err != nil {
		log.Fatal("FATAL ERROR: " + err.Error())
	}

	ex := Extractor{
		namePrefix:              namePrefix,
		nameSuffix:              nameSuffix,
		useType:                 useType,
		progress:                progress,
		ignoreXmlDecodingErrors: ignoreXmlDecodingErrors,
		initted:                 false,
	}

	if DEBUG {
		log.Print("extracting")
	}

	m := &ex
	m.init()

	for i, _ := range sources {
		err = m.extract(sources[i].getReader())

		if err != nil {
			log.Println("ERROR: " + err.Error())
			if !ignoreXmlDecodingErrors {
				log.Fatal("FATAL ERROR: " + err.Error())
			}
		}
	}

	ex.done()

	switch {
	case codeGenConvert:
		generateGoCode(os.Stdout, sourceNames, &ex)

	case structsToStdout:
		generateGoStructs(os.Stdout, sourceNames[0], &ex)

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
		printJavaJaxbMain(onlyChild.makeJavaType(namePrefix, ""), javaDir, javaPackage, getFullPath(sourceNames[0]), date)
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

func makeSourceReaders(sourceNames []string, url bool, standardIn bool) ([]Source, error) {

	var sources []Source
	for i, _ := range sourceNames {

		if url {
			sources[i] = new(UrlSource)
			if DEBUG {
				log.Print("Making UrlSource")
			}
		} else {
			if standardIn {
				sources[i] = new(StdinSource)
				if DEBUG {
					log.Print("Making StdinSource")
				}
			} else {
				sources[i] = new(FileSource)
				if DEBUG {
					log.Print("Making FileSource")
				}
			}
		}
		if DEBUG {
			log.Print("Making Source:[" + sourceNames[i] + "]")
		}
	}

	return sources, nil
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

func makeOneLevelDown(node *Node, globalTagAttributes map[string]([]*FQN)) []*XMLType {
	var children []*XMLType

	for _, np := range node.children {
		if np == nil {
			continue
		}
		for _, n := range np.children {
			if n == nil {
				continue
			}
			if flattenStrings && isStringOnlyField(n, len(globalTagAttributes[nk(n)])) {
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

func generateGoStructs(out io.Writer, sourceName string, ex *Extractor) {
	lineChannel := make(chan string, 100)
	var writer Writer
	writer = new(stdoutWriter)
	writer.open("", lineChannel)
	printGoStructVisitor := new(PrintGoStructVisitor)
	printGoStructVisitor.init(lineChannel, 999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
	printGoStructVisitor.Visit(ex.root)
	structSort(printGoStructVisitor)
	close(lineChannel)
	writer.close()
}

func generateGoCode(out io.Writer, sourceNames []string, ex *Extractor) {
	lineChannel := make(chan string, 100)
	var writer Writer
	sWriter := new(stringWriter)
	writer = sWriter
	writer.open("", lineChannel)
	printGoStructVisitor := new(PrintGoStructVisitor)
	printGoStructVisitor.init(lineChannel, 9999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
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
		OneLevelDownXML: makeOneLevelDown(ex.root, ex.globalTagAttributes),
		Filenames:       getFullPaths(sourceNames),
		Filename:        getFullPath(sourceNames[0]),
		Structs:         sWriter.s,
	}
	t := template.Must(template.New("chidleyGen").Parse(codeTemplate))

	err := t.Execute(out, x)
	if err != nil {
		log.Println("executing template:", err)
	}

}
