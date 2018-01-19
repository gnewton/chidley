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

//https://github.com/gnewton/chidley/issues/24
const mixedCaseSameNameXML_Issue24 = `
<?xml version="1.0" encoding="UTF-8"?>
<cluster ContinentType="Outlands" ClusterQuality="High">
		  <RareState state="0" weight="1556" />
		  <RareState state="1" weight="400" />
		  <RareState state="2" weight="40" />
		  <RareState state="3" weight="4" />
 </cluster>
<Cluster default="0.05" minimum="0" maximum="0.05"/>
`

const tagsContainHyphens = `
<?xml version="1.0" encoding="UTF-8"?>
<name>
	<first-name>Bill</first-name>
	<last-name>Smith</last-name>
</name>
`

//https://github.com/gnewton/chidley/issues/14
const githubIssue14 = `
<con1:actions>
   <con2:route xmlns:con2="http://www.bea.com/wli/sb/stages/routing/config">
     <con3:id xmlns:con3="http://www.bea.com/wli/sb/stages/config">_ActionId-3525062221263473230--64db5972.154ce25275c.-7fc5</con3:id>
     <con2:service ref="Something/Proxy/Something" xsi:type="ref:ProxyRef" xmlns:ref="http://www.bea.com/wli/sb/reference" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>
     <con2:operation passThrough="true"/>
     <con2:outboundTransform/>
     <con3:outboundTransform/>
     <con2:responseTransform/>
   </con2:route>
</con1:actions>
`
