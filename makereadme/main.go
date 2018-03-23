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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var ri ReadmeInfo
	var err error
	var tmp []byte

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

	//log.Println(out)

	t := template.Must(template.New("readmeTemplate").Parse(readmeTemplate))

	err = t.Execute(os.Stdout, ri)
	if err != nil {
		log.Println("executing template:", err)
	}

	//os.Setenv("PATH", "/home/newtong/go/src/github.com/gnewton/chidley")

}
