package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var nameSuffix = "_Type"
var namePrefix = "Chi_"
var attributePrefix = "Attr_"
var VERIFY = false
var DEBUG = false

type Writer interface {
	open(s string, lineChannel chan string) error
	close()
}

func init() {
	flag.StringVar(&nameSuffix, "s", nameSuffix, "Suffix to element names")
	flag.StringVar(&namePrefix, "p", namePrefix, "Prefix to element names")
	flag.StringVar(&attributePrefix, "a", attributePrefix, "Prefix to attribute names")
	flag.BoolVar(&VERIFY, "V", VERIFY, "Do full code generation & see if it can decode the original source file")
	flag.BoolVar(&DEBUG, "D", DEBUG, "Debug; prints out much information")
}

func handleParameters() bool {
	flag.Parse()
	return true
}

func main() {
	handleParameters()

	if len(flag.Args()) != 1 {
		fmt.Println("chidley <flags> xmlFileName")
		flag.Usage()
		return
	}
	xmlFilename := flag.Args()[0]

	extractor := Extractor{
		namePrefix:  namePrefix,
		nameSuffix:  nameSuffix,
		verify:      VERIFY,
		xmlFilename: xmlFilename,
	}

	err := extractor.extract()

	if err != nil {
		log.Fatal(err)
		fmt.Println("**********")
		os.Exit(42)
	}

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
