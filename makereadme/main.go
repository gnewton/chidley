package main

// Copyright 2014,2015,2016 Glen Newton
// glen.newton@gmail.com

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type ReadmeInfo struct {
	ChidleyUsage                                   string
	GeneratedUsage                                 string
	GeneratedXMLToJson                             string
	GeneratedXMLToXML                              string
	GeneratedCountElements                         string
	ChidleyOnlyStructOutput                        string
	SimpleExampleXMLFile                           string
	SimpleExampleXMLChidleyGoStructs               string
	SimpleExampleXMLChidleyGoStructsCollapsed      string
	SimpleExampleXMLChidleyGoStructsWithTypes      string
	PubmedXMLFileName                              string
	PubmedExampleXMLChidleyGoStructsWithTypes      string
	PubmedExampleXMLChidleyGoStructsWithTypeTiming string
	GeneratedPubmedCount                           string
	GeneratedPubmedNoStreaming                     string
	GeneratedPubmedStreaming                       string
	GeneratedPubmedXMLToJson                       string
	GeneratedPubmedXMLToXML                        string
	ChidleyGenerateJava                            string
	ChidleyGenerateJavaChangePackageName           string
	ChidleyGenerateJavaMavenBuild                  string
	ChidleyGenerateJavaRun                         string
}

const empty = `
    ***********EMPTY`

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	var tmp []byte

	ri := ReadmeInfo{
		ChidleyUsage:                                   empty,
		GeneratedUsage:                                 empty,
		GeneratedXMLToJson:                             empty,
		GeneratedXMLToXML:                              empty,
		GeneratedCountElements:                         empty,
		ChidleyOnlyStructOutput:                        empty,
		SimpleExampleXMLFile:                           empty,
		SimpleExampleXMLChidleyGoStructs:               empty,
		SimpleExampleXMLChidleyGoStructsCollapsed:      empty,
		SimpleExampleXMLChidleyGoStructsWithTypes:      empty,
		PubmedXMLFileName:                              empty,
		PubmedExampleXMLChidleyGoStructsWithTypes:      empty,
		PubmedExampleXMLChidleyGoStructsWithTypeTiming: empty,
		GeneratedPubmedCount:                           empty,
		GeneratedPubmedNoStreaming:                     empty,
		GeneratedPubmedStreaming:                       empty,
		GeneratedPubmedXMLToJson:                       empty,
		GeneratedPubmedXMLToXML:                        empty,
		ChidleyGenerateJava:                            empty,
		ChidleyGenerateJavaChangePackageName:           empty,
		ChidleyGenerateJavaMavenBuild:                  empty,
		ChidleyGenerateJavaRun:                         empty,
	}

	ri.ChidleyUsage, err = runCaptureStdout("..", "./chidley", "-n")

	if err != nil {
		log.Fatal(err)
	}

	tmp, err = ioutil.ReadFile("../data/test.xml")
	if err != nil {
		log.Fatal(err)
	}
	ri.SimpleExampleXMLFile = string(tmp)

	ri.SimpleExampleXMLChidleyGoStructs, err = runCaptureStdout("..", "./chidley", "-G", "data/test.xml")

	if err != nil {
		log.Fatal(err)
	}

	ri.SimpleExampleXMLChidleyGoStructsCollapsed = string(tmp)

	ri.SimpleExampleXMLChidleyGoStructsCollapsed, err = runCaptureStdout("..", "./chidley", "-G", "-F", "data/test.xml")
	if err != nil {
		log.Fatal(err)
	}

	ri.SimpleExampleXMLChidleyGoStructsWithTypes, err = runCaptureStdout("..", "./chidley", "-G", "-t", "data/test.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = generateSimpleGoCode()

	if err != nil {
		log.Fatal(err)
	}

	// Build generated code
	var output string
	output, err = runCaptureStdout("gencode", "go", "build")
	if err != nil {
		log.Println(output)
		log.Fatal(err)
	}

	// Run generated code, usage
	ri.GeneratedUsage, err = runCaptureStdout("gencode", "./gencode", "-h")
	if err != nil {
		log.Fatal(err)
	}

	// Run generated code, convert XML to JSON
	ri.GeneratedXMLToJson, err = runCaptureStdout("gencode", "./gencode", "-j")
	if err != nil {
		log.Fatal(err)
	}

	// Run generated code, convert XML to XML
	ri.GeneratedXMLToXML, err = runCaptureStdout("gencode", "./gencode", "-x")
	if err != nil {
		log.Fatal(err)
	}

	// Run generated code, count XML tags
	ri.GeneratedCountElements, err = runCaptureStdout("gencode", "./gencode", "-c")
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(out)

	if true {
		t := template.Must(template.New("readmeTemplate").Parse(readmeTemplate))

		err = t.Execute(os.Stdout, ri)

		if err != nil {
			log.Println("executing template:", err)
		}
	}
	//os.Setenv("PATH", "/home/newtong/go/src/github.com/gnewton/chidley")

}

func generateSimpleGoCode() error {
	//tmp, err := runCaptureStdout(".", "mkdir", "gencode")
	_, _ = runCaptureStdout(".", "mkdir", "gencode")
	tmp, err := runCaptureStdout(".", "../chidley", "-W", "../data/test.xml") //> gencode/main.go)

	if err != nil {
		log.Println(err)
		return err
	}

	//file, err := os.Create("gencode/main.go")
	file, _ := os.Create("gencode/main.go")
	file.WriteString(tmp)
	file.Sync()
	file.Close()

	//_, _ = runCaptureStdout(".", "rmdir", "gencode")
	return err
}
