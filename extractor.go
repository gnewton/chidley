package main

import (
	"encoding/xml"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

var nameMapper = map[string]string{
	"-": "_",
	".": "_dot_",
}

type Extractor struct {
	//globalTagAttributes map[string](map[string]string)
	globalTagAttributes    map[string]([]*FQN)
	globalTagAttributesMap map[string]bool
	globalNodeMap          map[string]*Node
	namePrefix             string
	nameSpaceTagMap        map[string]string
	nameSuffix             string
	reader                 io.Reader
	root                   *Node
	firstNode              *Node
	hasStartElements       bool
	useType                bool
	progress               bool
}

func (ex *Extractor) extract() error {
	ex.globalTagAttributes = make(map[string]([]*FQN))
	ex.globalTagAttributesMap = make(map[string]bool)
	ex.nameSpaceTagMap = make(map[string]string)
	ex.globalNodeMap = make(map[string]*Node)

	decoder := xml.NewDecoder(ex.reader)

	ex.root = new(Node)
	ex.root.initialize("root", "", "", nil)

	ex.hasStartElements = false

	tokenChannel := make(chan xml.Token, 3000)
	handleTokensDoneChannel := make(chan bool)

	go handleTokens(tokenChannel, ex, handleTokensDoneChannel)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		if token == nil {
			break
		}
		tokenChannel <- token
	}
	close(tokenChannel)
	_ = <-handleTokensDoneChannel
	return nil
}

func handleTokens(tChannel chan xml.Token, ex *Extractor, handleTokensDoneChannel chan bool) {
	depth := 0
	thisNode := ex.root
	first := true
	var progressCounter int64 = 0

	for token := range tChannel {
		switch element := token.(type) {
		case xml.Comment:
			if DEBUG {
				log.Print(thisNode.name)
				log.Printf("Comment: %+v\n", string(element))
			}

		case xml.ProcInst:
			if DEBUG {
				log.Printf("ProcInst: %+v\n", element)
			}

		case xml.Directive:
			if DEBUG {
				log.Printf("Directive: %+v\n", string(element))
			}

		case xml.StartElement:
			thisNode.tempCharData = ""
			progressCounter += 1
			if DEBUG {
				log.Printf("StartElement: %+v\n", element)
			}
			ex.hasStartElements = true

			if element.Name.Local == "" {
				continue
			}
			thisNode = ex.handleStartElement(element, thisNode)
			if first {
				first = false
				ex.firstNode = thisNode
			}
			depth += 1
			if progress {
				if progressCounter%50000 == 0 {
					log.Print(progressCounter)
				}
			}

		case xml.CharData:
			if DEBUG {
				log.Print(thisNode.name)
				log.Printf("CharData: [%+v]\n", string(element))
			}
			/*
				if !isJustSpacesAndLinefeeds(string(element)) {
					thisNode.nodeTypeInfo.checkFieldType(string(element))
					thisNode.hasCharData = true
				}
			*/
			if !thisNode.hasCharData {
				thisNode.tempCharData += strings.TrimSpace(string(element))
			}
			if thisNode.name == "RDF" {
				log.Print("[" + string(element) + "]")
			}

		case xml.EndElement:
			if !thisNode.hasCharData && !isJustSpacesAndLinefeeds(thisNode.tempCharData) {
				thisNode.nodeTypeInfo.checkFieldType(thisNode.tempCharData)
				thisNode.hasCharData = true
				thisNode.tempCharData = ""
			}

			depth -= 1
			if DEBUG {
				log.Printf("EndElement: %+v\n", element)
			}
			for key, c := range thisNode.childCount {
				if c > 1 {
					thisNode.children[key].repeats = true
				}
				thisNode.childCount[key] = 0
			}
			//if thisNode.parent != nil {
			if thisNode.peekParent() != nil {
				//thisNode = thisNode.parent
				thisNode = thisNode.popParent()
			}
		}
	}
	handleTokensDoneChannel <- true
	close(handleTokensDoneChannel)
}

func space(n int) string {
	s := strconv.Itoa(n) + ":"
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

func (ex *Extractor) printTree(n *Node, lineChannel chan string, d int, startName string, foundStartString bool) {
	if n.name == startName {
		foundStartString = true
	}
	repeats := ""
	if n.repeats {
		repeats = "*"
	}
	if foundStartString {
		lineChannel <- indent(d) + n.name + repeats
		d += 1
	}

	for _, v := range n.children {
		ex.printTree(v, lineChannel, d, startName, foundStartString)
	}
}

func (ex *Extractor) printStruct(n *Node, lineChannel chan string, startName string, foundStartString bool, alreadyPrinted map[string]bool) {
	_, ok := alreadyPrinted[nk(n)]
	if ok {
		return
	}
	alreadyPrinted[nk(n)] = true

	if nk(n) == startName {
		foundStartString = true
	}
	if foundStartString {
		if len(n.children) > 0 {
			lineChannel <- "type " + n.makeType(namePrefix, nameSuffix) + " struct {"
			ex.printInternalFields(n, lineChannel)
			lineChannel <- "}\n"
		}
	}

	for _, v := range n.children {
		ex.printStruct(v, lineChannel, startName, foundStartString, alreadyPrinted)
	}
}

func (ex *Extractor) printInternalFields(n *Node, lineChannel chan string) {
	attributes := ex.globalTagAttributes[nk(n)]
	fields := makeAttributes(attributes, ex.nameSpaceTagMap)
	var xmlName string
	if n.space != "" {
		xmlName = "\tXMLName  xml.Name `xml:\"" + n.space + " " + n.name + ",omitempty\" json:\",omitempty\"`"
	} else {
		xmlName = "\tXMLName  xml.Name `xml:\"" + n.name + ",omitempty\" json:\",omitempty\"`"
	}
	fields = append(fields, xmlName)

	var field string

	for _, v := range n.children {
		field = "\t " + v.makeType(namePrefix, nameSuffix) + " "
		//childAttributes := ex.globalTagAttributes[nk(v)]
		//if len(v.children) == 0 && len(childAttributes) == 0 {
		if len(v.children) == 0 {
			if v.repeats {
				field += "[]"
			}
			field += findType(v.nodeTypeInfo, ex.useType)

		} else {
			if v.repeats {
				field += "[]"
			} else {
				field += "*"
			}
			field += v.makeType(namePrefix, nameSuffix)
		}
		var xmlString string
		if v.space != "" {
			xmlString = " `xml:\"" + v.space + " " + v.name + ",omitempty\" json:\",omitempty\"`"
		} else {
			xmlString = " `xml:\"" + v.name + ",omitempty\" json:\",omitempty\"`"
		}
		field += xmlString
		fields = append(fields, field)
	}

	if len(n.children) == 0 && len(attributes) > 0 || n.hasCharData {
		xmlString := " `xml:\",chardata\" json:\",omitempty\"`"
		//charField := "\t" + capitalizeFirstLetter(n.name) + " string" + xmlString
		charField := "\t" + "Text" + " string" + xmlString
		fields = append(fields, charField)
	}
	sort.Strings(fields)

	for i := 0; i < len(fields); i++ {
		lineChannel <- fields[i]
	}
}

func (ex *Extractor) findNewNameSpaces(attrs []xml.Attr) {
	for _, attr := range attrs {
		if attr.Name.Space == "xmlns" {
			ex.nameSpaceTagMap[attr.Value] = attr.Name.Local
		}
	}
}

var full struct{}

func (ex *Extractor) handleStartElement(startElement xml.StartElement, thisNode *Node) *Node {
	name := startElement.Name.Local
	space := startElement.Name.Space

	ex.findNewNameSpaces(startElement.Attr)

	var child *Node
	var attributes []*FQN
	key := nks(space, name)
	child, ok := thisNode.children[key]
	// Does thisNode node already exist as child
	//fmt.Println(space, name)
	if ok {
		thisNode.childCount[key] += 1
		attributes, ok = ex.globalTagAttributes[key]
	} else {
		// if thisNode node does not already exist as child, it may still exist as child on other node:
		child, ok = ex.globalNodeMap[key]
		if !ok {
			child = new(Node)
			ex.globalNodeMap[key] = child
			spaceTag, _ := ex.nameSpaceTagMap[space]
			child.initialize(name, space, spaceTag, thisNode)
			thisNode.childCount[key] = 1

			attributes = make([]*FQN, 0, 2)
			ex.globalTagAttributes[key] = attributes
		} else {
			attributes = ex.globalTagAttributes[key]
		}
		thisNode.children[key] = child
	}
	child.pushParent(thisNode)

	for _, attr := range startElement.Attr {
		bigKey := key + "_" + attr.Name.Space + "_" + attr.Name.Local
		_, ok := ex.globalTagAttributesMap[bigKey]
		//_, _ = ex.globalTagAttributesMap[bigKey]
		if !ok {
			fqn := new(FQN)
			fqn.name = attr.Name.Local
			fqn.space = attr.Name.Space
			attributes = append(attributes, fqn)
			ex.globalTagAttributesMap[bigKey] = true
		}
	}
	ex.globalTagAttributes[key] = attributes
	return child
}

func isJustSpacesAndLinefeeds(s string) bool {
	s = strings.Replace(s, "\\n", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return len(strings.TrimSpace(s)) == 0
}
