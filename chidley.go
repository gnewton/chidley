package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var DEBUG = false
var progress = false
var attributePrefix = "Attr"
var structsToStdout = false
var prettyPrint = false
var codeGenConvert = false
var codeGenDir = "codegen"
var codeGenFilename = "CodeGenStructs.go"
var namePrefix = "Chi"
var nameSuffix = "Type"
var url = false
var useType = false

type Writer interface {
	open(s string, lineChannel chan string) error
	close()
}

var outputs = []*bool{
	&codeGenConvert,
	&structsToStdout,
}

func init() {
	flag.BoolVar(&DEBUG, "D", DEBUG, "Debug; prints out much information")
	flag.BoolVar(&codeGenConvert, "W", codeGenConvert, "Generate Go code to convert XML to JSON or XML (latter useful for validation)")
	flag.BoolVar(&structsToStdout, "c", structsToStdout, "Write generated Go structs to stdout")

	flag.BoolVar(&prettyPrint, "P", prettyPrint, "Pretty-print json in generated code (if applicable)")
	flag.BoolVar(&progress, "R", progress, "Progress: every 50000 elements")
	flag.BoolVar(&url, "u", url, "Filename interpreted as an URL")
	flag.BoolVar(&useType, "t", useType, "Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete")
	flag.StringVar(&attributePrefix, "a", attributePrefix, "Prefix to attribute names")
	flag.StringVar(&namePrefix, "p", namePrefix, "Prefix to element names")
	flag.StringVar(&nameSuffix, "s", nameSuffix, "Suffix to element names")
}

func handleParameters() error {
	flag.Parse()

	numBoolsSet := countNumberOfBoolsSet(outputs)
	if numBoolsSet > 1 {
		log.Print("  ERROR: Only one of -W -J -X -V -c can be set")
		return nil
	}
	if numBoolsSet == 0 {
		log.Print("  ERROR: At least one of -W -J -X -V -c must be set")
		return nil
	}
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := handleParameters()

	if err != nil {
		flag.Usage()
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Println("chidley <flags> xmlFileName|url")
		fmt.Println("xmlFileName can be .gz or .bz2: uncompressed transparently")
		flag.Usage()
		return
	}

	var sourceName string

	sourceName, err = filepath.Abs(flag.Args()[0])
	if err != nil {
		log.Fatal("FATAL ERROR: " + err.Error())
	}

	source, err := makeSourceReader(sourceName, url)
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
		printStructVisitor := new(PrintStructVisitor)
		printStructVisitor.init(lineChannel, 9999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType)
		printStructVisitor.Visit(ex.root)
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
		t := template.Must(template.New("letter").Parse(codeTemplate))

		err := t.Execute(os.Stdout, x)
		if err != nil {
			log.Println("executing template:", err)
		}
		break

	case structsToStdout:
		writer = new(stdoutWriter)
		writer.open("", lineChannel)
		printStructVisitor := new(PrintStructVisitor)
		printStructVisitor.init(lineChannel, 999, ex.globalTagAttributes, ex.nameSpaceTagMap, useType)
		printStructVisitor.Visit(ex.root)
		close(lineChannel)
		break
	}

}

func makeSourceReader(sourceName string, url bool) (Source, error) {
	var source Source
	if url {
		source = new(UrlSource)
		if DEBUG {
			log.Print("Making UrlSource")
		}
	} else {
		source = new(FileSource)
		if DEBUG {
			log.Print("Making FileSource")
		}
	}
	if DEBUG {
		log.Print("Making Source:[" + sourceName + "]")
	}
	err := source.newSource(sourceName)

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

func capitalizeFirstLetter(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
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
