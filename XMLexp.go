package main

import (
	"encoding/xml"
	"log"
	"strings"
)

const modx = `
<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:gx="http://www.google.com/kml/ext/2.2">
	<Document>
		<name>MODIS Imagery Overlays</name>
		<Snippet maxLines="1">Deepwater Horizon Oil Spill</Snippet>
		<open>1</open>
		<atom:author><atom:name>Google Crisis Response</atom:name></atom:author>
		<Author>ogle Crisis Response</Author>
		<author>ogle Crisis Response</author>
		<name>2010-06-19 - MODIS</name>
	</Document> 
</kml>
`

func EXP() {
	log.Println(modx)
	decoder := xml.NewDecoder(strings.NewReader(modx))

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				// OK
				break
			}
			log.Println(err)
		}
		if token == nil {
			log.Println("Empty token")
			break
		}
		log.Println("TTT", token)

		switch element := token.(type) {
		case xml.Comment:
			log.Println("O", string(element))

		case xml.ProcInst:
			log.Println("P", element.Target)
			log.Println("P---", string(element.Inst))

		case xml.Directive:
			log.Println("D", string(element))

		case xml.StartElement:
			log.Println("START", element.Name)
			log.Println("START-", element.Attr)

		case xml.CharData:
			log.Println("C", string(element))

		case xml.EndElement:
			log.Println("E", element.Name)
		}
	}
}
