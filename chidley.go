package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"
)

var DEBUG = false
var progress = false
var attributePrefix = "Attr"
var structsToStdout = false
var prettyPrint = false
var codeGenCounter = false
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
	&codeGenCounter,
	&codeGenConvert,
}

func init() {
	flag.BoolVar(&DEBUG, "D", DEBUG, "Debug; prints out much information")
	flag.BoolVar(&codeGenConvert, "W", codeGenConvert, "Generate Go code to convert XML to JSON or XML (latter useful for validation)")
	flag.BoolVar(&codeGenCounter, "C", codeGenCounter, "Generate Go code that counts bumber of each unique XML tag in XML file")
	flag.BoolVar(&prettyPrint, "P", prettyPrint, "Pretty-print json in generated code (if applicable)")
	flag.BoolVar(&progress, "R", progress, "Progress: every 50000 elements")
	flag.BoolVar(&structsToStdout, "c", structsToStdout, "Write generated Go structs to stdout")
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

	sourceName = flag.Args()[0]

	source, err := makeSourceReader(sourceName, url)
	if err != nil {
		return
	}
	defer source.Close()

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
		os.Exit(42)
	}

	var writer Writer
	lineChannel := make(chan string)

	switch {
	case codeGenCounter:
		alreadyPrinted := make(map[string]bool)
		var codegen *CodeGenerator
		codegen = new(CodeGenerator)
		newSource, err := source.copySource()
		writer, err = codegen.init(ex.firstNode, ex.globalNodeMap, namePrefix, ex.nameSpaceTagMap, nameSuffix, codeGenDir, codeGenFilename, newSource, lineChannel)
		codegen.generateCodePre(ex.hasStartElements, url, false, false)
		if err != nil {
			log.Fatal(err)
			return
		}
		ex.printStruct(ex.firstNode, lineChannel, "", true, alreadyPrinted)
		codegen.generateVerifyCode(ex.hasStartElements, ex.globalTagAttributes, url)
		break

	case codeGenConvert:
		sWriter := new(stringWriter)
		writer = sWriter
		writer.open("", lineChannel)
		printStructVisitor := new(PrintStructVisitor)
		printStructVisitor.init(lineChannel)
		printStructVisitor.Visit(ex.root)

		x := XmlType{
			NameType:     ex.firstNode.makeType(namePrefix, nameSuffix),
			XMLName:      ex.firstNode.name,
			XMLNameUpper: capitalizeFirstLetter(ex.firstNode.name),
			XMLSpace:     ex.firstNode.space,
			Filename:     sourceName,
			Structs:      sWriter.s,
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
		printStructVisitor.init(lineChannel)
		//visitNode(ex.root, printStructVisitor)
		printStructVisitor.Visit(ex.root)
		break
	}
	close(lineChannel)
	//writer.close()

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
