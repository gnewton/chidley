package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var DEBUG = false
var attributePrefix = "Attr_"
var structsToStdout = true
var prettyPrint = false
var codeGenVerify = false
var codeGenConvertToJson = false
var codeGenConvertToXML = false
var codeGenDir = "codegen"
var codeGenFilename = "CodeGenStructs.go"
var namePrefix = "Chi_"
var nameSuffix = "_Type"
var url = false
var useType = false

type Writer interface {
	open(s string, lineChannel chan string) error
	close()
}

func init() {
	flag.BoolVar(&DEBUG, "D", DEBUG, "Debug; prints out much information")
	flag.BoolVar(&codeGenConvertToJson, "J", codeGenConvertToJson, "Do Go code generation to convert XML to JSON")
	flag.BoolVar(&codeGenConvertToXML, "X", codeGenConvertToXML, "Do Go code generation to convert XML to XML (useful for validation)")
	flag.BoolVar(&codeGenVerify, "V", codeGenVerify, "Do code generation of code that reads XML and counts every tag instance: space_prefix:tag")
	flag.BoolVar(&prettyPrint, "P", prettyPrint, "Pretty-print json in generated code (if applicable)")
	flag.BoolVar(&structsToStdout, "c", structsToStdout, "Write generated Go structs to stdout")
	flag.BoolVar(&url, "u", url, "Filename interpreted as an URL")
	flag.BoolVar(&useType, "t", useType, "Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete")
	flag.StringVar(&attributePrefix, "a", attributePrefix, "Prefix to attribute names")
	flag.StringVar(&namePrefix, "p", namePrefix, "Prefix to element names")
	flag.StringVar(&nameSuffix, "s", nameSuffix, "Suffix to element names")
}

func handleParameters() bool {
	flag.Parse()
	return true
}

func main() {
	handleParameters()

	if len(flag.Args()) != 1 {
		fmt.Println("chidley <flags> xmlFileName|url")
		fmt.Println("xmlFileName can be .gz or .bz2: uncompressed transparently")
		flag.Usage()
		return
	}

	x := make(map[string]interface{})
	x["goo"] = x

	var sourceName string

	sourceName = flag.Args()[0]

	source, err := makeSourceReader(sourceName, url)
	if err != nil {

	}
	defer source.Close()

	ex := Extractor{
		namePrefix: namePrefix,
		nameSuffix: nameSuffix,
		reader:     source.getReader(),
		useType:    useType,
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
	case codeGenVerify:
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

	case codeGenConvertToJson:
		alreadyPrinted := make(map[string]bool)
		var codegen *CodeGenerator
		codegen = new(CodeGenerator)
		newSource, err := source.copySource()
		writer, err = codegen.init(ex.firstNode, ex.globalNodeMap, namePrefix, ex.nameSpaceTagMap, nameSuffix, codeGenDir, codeGenFilename, newSource, lineChannel)
		if err != nil {
			log.Fatal(err)
			return
		}
		codegen.generateCodePre(ex.hasStartElements, url, true, false)
		ex.printStruct(ex.firstNode, lineChannel, "", true, alreadyPrinted)
		codegen.generateConvertToCode(ex.firstNode, codeGenConvertToJson, codeGenConvertToXML, prettyPrint)
		break

	case codeGenConvertToXML:
		alreadyPrinted := make(map[string]bool)
		var codegen *CodeGenerator
		codegen = new(CodeGenerator)
		newSource, err := source.copySource()
		writer, err = codegen.init(ex.firstNode, ex.globalNodeMap, namePrefix, ex.nameSpaceTagMap, nameSuffix, codeGenDir, codeGenFilename, newSource, lineChannel)
		if err != nil {
			log.Fatal(err)
			return
		}
		codegen.generateCodePre(ex.hasStartElements, url, false, true)
		ex.printStruct(ex.firstNode, lineChannel, "", true, alreadyPrinted)
		codegen.generateConvertToCode(ex.firstNode, codeGenConvertToJson, codeGenConvertToXML, prettyPrint)
		break
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
	writer.close()

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
