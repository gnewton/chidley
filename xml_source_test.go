package main

const sameNameDifferentNameSpaceXML = `
<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:gx="http://www.google.com/kml/ext/2.2">
	<Document>
		<open>1</open>
		<atom:author><atom:name>Google Crisis Response</atom:name></atom:author>
		<author>ogle Crisis Responsex</author>
	</Document> 
</kml>
`

const mixedCaseSameNameXML = `
<?xml version="1.0" encoding="UTF-8"?>
<kml>
	<name>MODIS Imagery Overlays</name>
	<Name>MODIS Imagery Overlays</Name>
</kml>
`
