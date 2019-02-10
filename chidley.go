package main

// Copyright 2014,2015,2016 Glen Newton
// glen.newton@gmail.com

import (
	"bufio"
	"bytes"
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
	flag.BoolVar(&writeOtira, "O", writeOtira, "Generated Go code for Otira")
	flag.BoolVar(&xmlName, "x", xmlName, "Add XMLName (Space, Local) for each XML element, to JSON")
	flag.BoolVar(&keepXmlFirstLetterCase, "K", keepXmlFirstLetterCase, "Do not change the case of the first letter of the XML tag names")
	flag.BoolVar(&validateFieldTemplate, "m", validateFieldTemplate, "Validate the field template. Useful to make sure the template defined with -T is valid")

	flag.BoolVar(&ignoreLowerCaseXmlTags, "L", ignoreLowerCaseXmlTags, "Ignore lower case XML tags")

	flag.StringVar(&attributePrefix, "a", attributePrefix, "Prefix to attribute names")
	flag.StringVar(&baseJavaDir, "D", baseJavaDir, "Base directory for generated Java code (root of maven project)")
	flag.StringVar(&cdataStringName, "M", cdataStringName, "Set name of CDATA string field")
	flag.StringVar(&fieldTemplateString, "T", fieldTemplateString, "Field template for the struct field definition. Can include annotations. Default is for XML and JSON")
	flag.StringVar(&javaAppName, "k", javaAppName, "App name for Java code (appended to ca.gnewton.chidley Java package name))")
	flag.StringVar(&lengthTagAttribute, "A", lengthTagAttribute, "The tag name attribute to use for the max length Go annotations")
	flag.StringVar(&lengthTagName, "N", lengthTagName, "The tag name to use for the max length Go annotations")
	flag.StringVar(&lengthTagSeparator, "S", lengthTagSeparator, "The tag name separator to use for the max length Go annotations")
	flag.StringVar(&namePrefix, "e", namePrefix, "Prefix to struct (element) names; must start with a capital")
	flag.StringVar(&userJavaPackageName, "P", userJavaPackageName, "Java package name (rightmost in full package name")

	flag.StringVar(&ignoredXmlTags, "h", ignoredXmlTags, "List of XML tags to ignore; comma separated")
	flag.StringVar(&collapsedXmlTags, "Y", collapsedXmlTags, "Collapse down one level these tags. I.e. <PMID Version=\"1\">30516271</PMID> becomes fields pmid_version, pmid")

	flag.StringVar(&noDeeperThanTags, "H", noDeeperThanTags, "List of tags bekow which no tags are followed")

}

func handleParameters() error {
	flag.Parse()

	if codeGenConvert || writeJava {
		structsToStdout = false
	}

	numBoolsSet := countNumberOfBoolsSet(outputs)
	if numBoolsSet > 1 {
		log.Print("  ERROR: Only one of -O -W -J -X -V -c can be set")
	} else if numBoolsSet == 0 {
		log.Print("  ERROR: At least one of -O -W -J -X -V -c must be set")
	}
	if sortByXmlOrder {
		structSort = printStructsByXml
	}

	var err error

	ignoredXmlTagsMap, err = extractExcludedTags(ignoredXmlTags)
	if err != nil {
		return err
	}

	collapsedXmlTagsList, err = extractCollapsedTags(collapsedXmlTags)
	if err != nil {
		return err
	}

	noDeeperThanTagsMap, err := extractCollapsedTags(noDeeperThanTags)
	if err != nil {
		return err
	}

	log.Println(collapsedXmlTagsList)
	log.Println(noDeeperThanTagsMap)

	if lengthTagName == "" && lengthTagAttribute == "" || lengthTagName != "" && lengthTagAttribute != "" {
		return nil
	}

	return errors.New("Both lengthTagName and lengthTagAttribute must be set")
}

func main() {
	//log.Println(fieldTemplateString)

	//EXP()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := handleParameters()

	if err != nil {
		log.Println(err)
		flag.Usage()
		return
	}

	err = runValidateFieldTemplate(validateFieldTemplate)
	if err != nil {

		return
	}
	if validateFieldTemplate {
		return
	}

	if len(flag.Args()) == 0 && !readFromStandardIn {
		fmt.Println("chidley <flags> xmlFileName|url")
		fmt.Println("xmlFileName can be .gz or .bz2: uncompressed transparently")
		flag.Usage()
		return
	}

	var sourceNames []string

	if !readFromStandardIn {
		sourceNames = flag.Args()
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

	for source := range sources {
		if DEBUG {
			log.Println("READER", source)
		}
		err = m.extract(source.getReader())

		if err != nil {
			log.Println("ERROR: " + err.Error())
			if !ignoreXmlDecodingErrors {
				log.Fatal("FATAL ERROR: " + err.Error())
			}
		}
		if DEBUG {
			log.Println("DONE READER", source)
		}
	}

	ex.done()

	switch {

	case writeOtira:
		fmt.Println("Otira")
		generateOtiraCode(os.Stdout, sourceNames, &ex)

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
		fullPath, err := getFullPath(sourceNames[0])
		if err != nil {
			log.Fatal(err)
		}
		printJavaJaxbMain(onlyChild.makeJavaType(namePrefix, ""), javaDir, javaPackage, fullPath, date)
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

//func makeSourceReaders(sourceNames []string, url bool, standardIn bool) ([]Source, error) {
func makeSourceReaders(sourceNames []string, url bool, standardIn bool) (chan Source, error) {
	var err error
	//sources := make([]Source, len(sourceNames))
	sources := make(chan Source, 1)

	go func() {
		var newSource Source
		for i, _ := range sourceNames {
			if url {
				newSource = new(UrlSource)
				if DEBUG {
					log.Print("Making UrlSource")
				}
			} else {
				if standardIn {
					newSource = new(StdinSource)
					if DEBUG {
						log.Print("Making StdinSource")
					}
				} else {
					newSource = new(FileSource)
					if DEBUG {
						//log.Print("Making FileSource")
					}
				}
			}

			err = newSource.newSource(sourceNames[i])
			if err != nil {
				log.Fatal(err)
			}
			sources <- newSource
			if DEBUG {
				log.Print("Making Source:[" + sourceNames[i] + "]")
			}
		}
		close(sources)
	}()
	return sources, err
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
func printStructsByXml(v *PrintGoStructVisitor) error {
	orderNodes := make(map[int]*Node)
	var order []int

	for k := range v.alreadyVisitedNodes {
		nodeOrder := v.alreadyVisitedNodes[k].discoveredOrder
		orderNodes[nodeOrder] = v.alreadyVisitedNodes[k]
		order = append(order, nodeOrder)
	}
	sort.Ints(order)

	for o := range order {
		err := print(v, orderNodes[o])
		if err != nil {
			return err
		}
	}
	return nil
}

type ByDepth []*Node

func (a ByDepth) Len() int           { return len(a) }
func (a ByDepth) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDepth) Less(i, j int) bool { return a[i].minDepth > a[j].minDepth }

func printOtiraByDepth(v *PrintOtiraVisitor) error {
	nodes := make([]*Node, len(v.alreadyVisitedNodes))

	i := 0
	for k := range v.alreadyVisitedNodes {
		nodes[i] = v.alreadyVisitedNodes[k]
		i++
	}
	sort.Sort(ByDepth(nodes))

	one2mT := template.Must(template.New(one2manyTemplateName).Parse(one2manyTemplate))
	m2mT := template.Must(template.New(many2manyTemplateName).Parse(many2manyTemplate))
	stringFieldT := template.Must(template.New(stringFieldTemplateName).Parse(stringFieldTemplate))

	fieldCounter := 0
	for o := range nodes {
		err := printOtiraTables(v, nodes[o], &fieldCounter, stringFieldT, collapsedXmlTagsList)
		if err != nil {
			return err
		}
	}

	for o := range nodes {
		err := printOtiraRelations(v, nodes[o], one2mT, m2mT, collapsedXmlTagsList)
		if err != nil {
			return err
		}
	}
	return nil
}

// Alphabetical order
func printStructsAlphabetical(v *PrintGoStructVisitor) error {
	var keys []string
	for k := range v.alreadyVisitedNodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		err := print(v, v.alreadyVisitedNodes[k])
		if err != nil {
			return err
		}
	}
	return nil

}

func generateGoStructs(out io.Writer, sourceName string, ex *Extractor) {
	printGoStructVisitor := new(PrintGoStructVisitor)

	printGoStructVisitor.init(os.Stdout, 999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
	printGoStructVisitor.Visit(ex.root)
	structSort(printGoStructVisitor)
}

func generateOtiraCode(out io.Writer, sourceNames []string, ex *Extractor) error {
	// buf := bytes.NewBufferString("")
	// printGoStructVisitor := new(PrintGoStructVisitor)
	// printGoStructVisitor.init(buf, 9999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
	// printGoStructVisitor.Visit(ex.root)

	// fmt.Println(buf)

	//printOtiraVisitor := new(PrintOtiraVisitor)
	printOtiraVisitor := new(GenericVisitor)
	//printOtiraVisitor.init(os.Stdout, 999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
	printOtiraVisitor.init(999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
	printOtiraVisitor.Visit(ex.root)
	//printOtiraByDepth(printOtiraVisitor)
	for i, _ := range printOtiraVisitor.nodeInfoList {
		ni := printOtiraVisitor.nodeInfoList[i]
		log.Printf("%+v\n", ni)
		for j, _ := range ni.SubElements {
			log.Printf("\t\t%+v\n", ni.SubElements[j])
		}
	}
	return nil
}

//Writes structs to a string then uses this in a template to generate Go code
func generateGoCode(out io.Writer, sourceNames []string, ex *Extractor) error {
	buf := bytes.NewBufferString("")
	printGoStructVisitor := new(PrintGoStructVisitor)
	printGoStructVisitor.init(buf, 9999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType, nameSpaceInJsonName)
	printGoStructVisitor.Visit(ex.root)

	structSort(printGoStructVisitor)

	xt := XMLType{NameType: ex.firstNode.makeType(namePrefix, nameSuffix),
		XMLName:      ex.firstNode.name,
		XMLNameUpper: capitalizeFirstLetter(ex.firstNode.name),
		XMLSpace:     ex.firstNode.space,
	}

	fullPath, err := getFullPath(sourceNames[0])
	if err != nil {
		return err
	}

	fullPaths, err := getFullPaths(sourceNames)
	if err != nil {
		return err
	}
	x := XmlInfo{
		BaseXML:         &xt,
		OneLevelDownXML: makeOneLevelDown(ex.root, ex.globalTagAttributes),
		Filenames:       fullPaths,
		Filename:        fullPath,
		Structs:         buf.String(),
	}
	x.init()
	t := template.Must(template.New("chidleyGen").Parse(codeTemplate))

	err = t.Execute(out, x)
	if err != nil {
		log.Println("executing template:", err)
		return err
	}
	return err
}
