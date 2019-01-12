package main

import (
	"encoding/xml"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
)

const XML_NAMESPACE_ACRONYM = "xmlns"

var nameMapper = map[string]string{
	"-": "Hyphen",
	".": "Dot",
}

var DiscoveredOrder = 0

type Extractor struct {
	globalTagAttributes     map[string]([]*FQN)
	globalTagAttributesMap  map[string]bool
	globalNodeMap           map[string]*Node
	namePrefix              string
	nameSpaceTagMap         map[string]string
	nameSuffix              string
	root                    *Node
	firstNode               *Node
	hasStartElements        bool
	useType                 bool
	progress                bool
	ignoreXmlDecodingErrors bool
	initted                 bool
	tokenChannel            chan xml.Token
	handleTokensDoneChannel chan bool
}

const RootName = "ChidleyRoot314159"

func (ex *Extractor) init() {
	ex.globalTagAttributes = make(map[string]([]*FQN))
	ex.globalTagAttributesMap = make(map[string]bool)
	ex.nameSpaceTagMap = make(map[string]string)
	ex.globalNodeMap = make(map[string]*Node)
	ex.root = new(Node)
	//ex.root.initialize(RootName, "", "", nil)
	ex.root.initialize("", "", "", nil)
	ex.hasStartElements = false
	ex.initted = true
	ex.tokenChannel = make(chan xml.Token, 100)
	ex.handleTokensDoneChannel = make(chan bool)
	go handleTokens(ex)
}

func (ex *Extractor) done() {
	close(ex.tokenChannel)
	_ = <-ex.handleTokensDoneChannel
}

func (ex *Extractor) extract(reader io.Reader) error {
	if ex.initted == false {
		return errors.New("extractor not properly initted: must run extractor.init() first")
	}
	decoder := xml.NewDecoder(reader)

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				// OK
				break
			}
			log.Println(err)
			if !ex.ignoreXmlDecodingErrors {
				return err
			}
		}
		if token == nil {
			log.Println("Empty token")
			break
		}
		ex.tokenChannel <- xml.CopyToken(token)
	}
	return nil
}

func handleTokens(ex *Extractor) {
	tChannel := ex.tokenChannel
	handleTokensDoneChannel := ex.handleTokensDoneChannel
	depth := 0
	thisNode := ex.root
	first := true
	var progressCounter int64 = 0

	for token := range tChannel {
		//log.Println(token)
		switch element := token.(type) {

		case xml.Comment:
			if DEBUG {
				//log.Print(thisNode.name)
				//log.Printf("Comment: %+v\n", string(element))
			}

		case xml.ProcInst:
			if DEBUG {
				//log.Println("ProcInst: Target=" + element.Target + "  Inst=[" + string(element.Inst) + "]")
			}

		case xml.Directive:
			if DEBUG {
				//log.Printf("Directive: %+v\n", string(element))
			}

		case xml.StartElement:
			progressCounter += 1
			if DEBUG {
				//log.Printf("StartElement: %+v\n", element)
			}
			ex.hasStartElements = true

			if element.Name.Local == "" {
				continue
			}
			thisNode = ex.handleStartElement(element, thisNode)
			thisNode.tempCharData = ""
			thisNode.ignoredTag = isIgnoredTag(element.Name.Local)
			thisNode.minDepth = depth

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
				//log.Print(thisNode.name)
				//log.Printf("CharData: [%+v]\n", string(element))
			}

			//if !thisNode.hasCharData {
			charData := string(element)
			thisNode.tempCharData += charData //strings.TrimSpace(charData)
			thisNode.charDataCount += int64(len(charData))
		//}

		case xml.EndElement:
			//if ignoredTag(element.Name.Local) {
			//continue
			//}

			thisNode.nodeTypeInfo.checkFieldType(thisNode.tempCharData)
			thisNode.nodeTypeInfo.addFieldLength(thisNode.charDataCount)
			thisNode.charDataCount = 0

			if DEBUG {
				//log.Printf("EndElement: %+v\n", element)
				//log.Printf("[[" + thisNode.tempCharData + "]]")
				//log.Println("Char is empty: ", isJustSpacesAndLinefeeds(thisNode.tempCharData))
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
					thisNode.children[key].maxNumInstances = c
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
	s := strconv.Itoa(n) + "Space"
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

func (ex *Extractor) findNewNameSpaces(attrs []xml.Attr) {
	for _, attr := range attrs {
		if strings.HasPrefix(attr.Name.Space, XML_NAMESPACE_ACRONYM) {
			//log.Println("mmmmmmmmmmmmmmmmmmmmmmm", attr)
			//log.Println("+++++++++++++++++++++++++++", attr.Value, "|", attr.Name.Local, "|", attr.Name.Space)
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
			DiscoveredOrder += 1
			child.discoveredOrder = DiscoveredOrder
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

	// Extract attributes
	for _, attr := range startElement.Attr {
		bigKey := key + "_" + attr.Name.Space + "_" + attr.Name.Local
		_, ok := ex.globalTagAttributesMap[bigKey]
		if ok {
			fqn := findThisAttribute(attr.Name.Local, attr.Name.Space, ex.globalTagAttributes[key])
			if fqn == nil {
				log.Println("This should not be happening: fqn is nil")
				continue
			}
			lenValue := len(attr.Value)
			if lenValue > fqn.maxLength {
				fqn.maxLength = lenValue
			}
		} else {
			fqn := new(FQN)
			fqn.name = attr.Name.Local
			fqn.space = attr.Name.Space
			fqn.maxLength = len(attr.Value)
			//log.Println(name, "|", fqn.name, "||", fqn.space)
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
