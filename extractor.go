package main

import (
	"encoding/xml"
	"io"
	"log"
	"strconv"
	"strings"
)

var nameMapper = map[string]string{
	"-": "_",
	".": "_dot_",
}

type Extractor struct {
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

	tokenChannel := make(chan xml.Token, 100)
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
		tokenChannel <- xml.CopyToken(token)
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
				log.Println("ProcInst: Target=" + element.Target + "  Inst=[" + string(element.Inst) + "]")
			}

		case xml.Directive:
			if DEBUG {
				log.Printf("Directive: %+v\n", string(element))
			}

		case xml.StartElement:
			progressCounter += 1
			if DEBUG {
				log.Printf("StartElement: %+v\n", element)
			}
			ex.hasStartElements = true

			if element.Name.Local == "" {
				continue
			}
			thisNode = ex.handleStartElement(element, thisNode)
			thisNode.tempCharData = ""
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

			//if !thisNode.hasCharData {
			thisNode.tempCharData += strings.TrimSpace(string(element))
		//}

		case xml.EndElement:
			thisNode.nodeTypeInfo.checkFieldType(thisNode.tempCharData)

			if DEBUG {
				log.Printf("EndElement: %+v\n", element)
				log.Printf("[[" + thisNode.tempCharData + "]]")
				log.Printf("Char is empty: ", isJustSpacesAndLinefeeds(thisNode.tempCharData))
			}
			if !thisNode.hasCharData && !isJustSpacesAndLinefeeds(thisNode.tempCharData) {
				thisNode.hasCharData = true

			} else {

			}
			thisNode.tempCharData = ""
			depth -= 1

			for key, c := range thisNode.childCount {
				if c > 1 {
					thisNode.children[key].repeats = true
				}
				thisNode.childCount[key] = 0
			}
			if thisNode.peekParent() != nil {
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
