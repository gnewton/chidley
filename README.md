#`chidley`



##`chidley` converts *any* XML to Go structs (and therefor to JSON)
* By *any*, any XML that can be read by the Go [xml package](http://golang.org/pkg/encoding/xml/) decoder. 
* Where *convert* means, generates Go code that when compiled, will convert the XML to JSON
* or will just generate the Go structs that represent the input XML
* or converts XML to XML (useful for validation) 

Author: Glen Newton
Language: Go

##New

#2015.07.24
`chidley` now supports the user naming of the resulting JAXB Java class package.
Previously the package name could only be `ca/gnewton/chidley/jaxb`.
Now, using the `-P name`, the `jaxb` default can be altered.

So `chidley -J -P "foobar" sample.xml` will result in Java classes with package name: `ca/gnewton/chidley/foobar`.


###Previous

`chidley` now has support for Java/JAXB. It generates appropriate Java/JAXB classes and associated maven pom.

See [Java/JAXB section below](#user-content-java) for usage.


##How does it work (with such a small memory footprint)
`chidley` uses the input XML to build a model of each XML element (aka tag).
It examines each instance of a tag, and builds a (single) prototypical representation; that is, the union of all the attributes and all of the child elements of all instances of the tag.
So even if there are million instances of a specific tag, there is only one model tag representation.
Note that a tag is unique by its namespace+tagname combination (in the Go xml package parlance, [`space + local`](http://golang.org/pkg/encoding/xml/#Name).
###Types
`chidley` by default makes all values (attributes, tag content) in the generated Go structs a string (which is always valid for XML attributes and content), but it has a flag (`-t`) where it will detect and use the most appropriate type. 
`chidley` tries to fit the **smallest** Go type. 
For example, if all instances of a tag contain a number, and all instances are -128 to 127, then it will use an `int8` in the Go struct.

##`chidley` binary
Compiled for 64bit Linux Fedora18, go version go1.3 linux/amd64

##Usage
```

Usage of ./chidley:
  -B	Add database metadata to created Go structs
  -D string
    	Base directory for generated Java code (root of maven project) (default "java")
  -G	Only write generated Go structs to stdout
  -J	Generated Java code for Java/JAXB
  -P string
    	Java package name (rightmost in full package name
  -W	Generate Go code to convert XML to JSON or XML (latter useful for validation) and write it to stdout
  -X	Sort output of structs in Go code by order encounered in source XML  (default is alphabetical order)
  -a string
    	Prefix to attribute names (default "Attr_")
  -c	Read XML from standard input
  -d	Debug; prints out much information
  -e string
    	Prefix to struct (element) names; must start with a capital (default "Chi")
  -k string
    	App name for Java code (appended to ca.gnewton.chidley Java package name)) (default "jaxb")
  -n	Use the XML namespace prefix as prefix to JSON name; prefix followed by 2 underscores (__)
  -p	Pretty-print json in generated code (if applicable)
  -r	Progress: every 50000 input tags (elements)
  -t	Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete
  -u	Filename interpreted as an URL
  -x	Add XMLName (Space, Local) for each XML element, to JSON
$
```



###Specific Usages:
* `chidley -W ...`: writes Go code to standard out, so this output should be directed to a filename and subsequently be compiled. When compiled, the resulting binary will:
    * convert the XML file to JSON
    * or convert the XML file to XML (useful for validation)
    * or count the # of elements (space, local) in the XML file
* `chidley -G ...`: writes just the Go structs that represent the input XML. For incorporation into the user's code base.


###Example `chidley -W`:
```
$ chidley -W xml/test1.xml > examples/test1/ChidTest1.go
```
####Usage of generated code
```
$ cd examples/test1
$ go build
$ ./test1
Usage of ./test1:
  -c=false: Count each instance of XML tags
  -f="/home/gnewton/work/chidley/xml/test1.xml": XML file or URL to read in
  -h=false: Usage
  -j=false: Convert to JSON
  -s=false: Stream XML by using XML elements one down from the root tag. Good for huge XML files (see http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
  -x=false: Convert to XML
```

#####Generated code: Convert XML to JSON `-j`
```
$ ./test1 -j -f ../../xml/test1.xml 
{
 "doc": [
  {
   "Attr_type": "book",
   "author": {
    "firstName": {
     "Text": "Frank"
    },
    "last-name": {
     "Text": "Herbert"
    }
   },
   "title": {
    "Text": "Dune"
   }
  },
  {
   "Attr_type": "article",
   "author": {
    "firstName": {
     "Text": "Aldous"
    },
    "last-name": {
     "Text": "Huxley"
    }
   },
   "title": {
    "Text": "Brave New Wold"
   }
  }
 ]
}
```

#####Generated code: Convert XML to XML `-x`
```
$ ./test1 -x -f ../../xml/test1.xml 
  <Chi_docs>
      <doc type="book">
          <author>
              <firstName>Frank</firstName>
              <last-name>Herbert</last-name>
          </author>
          <title>Dune</title>
      </doc>
      <doc type="article">
          <author>
              <firstName>Aldous</firstName>
              <last-name>Huxley</last-name>
          </author>
          <title>Brave New Wold</title>
      </doc>
  </Chi_docs>
```

#####Generated code: Count elements `-c`
XML elements (or tags) are counted in the source file (space,local) and are printed-out, unsorted

```
$ ./test1 -c
1 _:docs
2 _:doc
2 _:title
2 _:author
2 _:last-name
2 _:firstName
```

**Note**: the underscore before the colon indicates there is no (or the default) namespace for the element. 

###Example `chidley -G`:
Just prints out the Go structs to standard out:
```
$ chidley -G xml/test1.xml
type Chi_root struct {
	Chi_docs *Chi_docs `xml:" docs,omitempty" json:"docs,omitempty"`
}

type Chi_docs struct {
	Chi_doc []*Chi_doc `xml:" doc,omitempty" json:"doc,omitempty"`
}

type Chi_doc struct {
	Attr_type string `xml:" type,attr"  json:",omitempty"`
	Chi_author *Chi_author `xml:" author,omitempty" json:"author,omitempty"`
	Chi_title *Chi_title `xml:" title,omitempty" json:"title,omitempty"`
}

type Chi_title struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_author struct {
	Chi_firstName *Chi_firstName `xml:" firstName,omitempty" json:"firstName,omitempty"`
	Chi_last_name *Chi_last_name `xml:" last-name,omitempty" json:"last-name,omitempty"`
}

type Chi_last_name struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_firstName struct {
	Text string `xml:",chardata" json:",omitempty"`
}
```

##Name changes: XML vs. Go structs
XML names can contain dots `.` and hyphens or dashes `-`. These are not valid identifiers for Go structs or variables. These are mapped as:
* `"-": "_"`
* `".": "_dot_"`

Note that the original XML name strings are used in the struct XML and JSON annotations for the element.

##Type example
```
<people>

  <person>
    <name>bill</name>
    <age>37</age>
    <married>true</married>
  </person>

  <person>
    <name>sarah</name>
    <age>24</age>
    <married>false</married>
  </person>

</people>
```

####Default
```
$ ./chidley -G xml/testType.xml
type Chi_root struct {
	Chi_people *Chi_people `xml:" people,omitempty" json:"people,omitempty"`
}

type Chi_people struct {
	Chi_person []*Chi_person `xml:" person,omitempty" json:"person,omitempty"`
}

type Chi_person struct {
	Chi_age *Chi_age `xml:" age,omitempty" json:"age,omitempty"`
	Chi_married *Chi_married `xml:" married,omitempty" json:"married,omitempty"`
	Chi_name *Chi_name `xml:" name,omitempty" json:"name,omitempty"`
}

type Chi_name struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_age struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_married struct {
	Text string `xml:",chardata" json:",omitempty"`
}
```

####Types turned on `-t`
```
$ ./chidley -G -t xml/testType.xml
$ type Chi_root struct {
	Chi_people *Chi_people `xml:" people,omitempty" json:"people,omitempty"`
}

type Chi_people struct {
	Chi_person []*Chi_person `xml:" person,omitempty" json:"person,omitempty"`
}

type Chi_person struct {
	Chi_age *Chi_age `xml:" age,omitempty" json:"age,omitempty"`
	Chi_married *Chi_married `xml:" married,omitempty" json:"married,omitempty"`
	Chi_name *Chi_name `xml:" name,omitempty" json:"name,omitempty"`
}

type Chi_name struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_age struct {
	Text int8 `xml:",chardata" json:",omitempty"`
}

type Chi_married struct {
	Text bool `xml:",chardata" json:",omitempty"`
}
```

Notice:
* `Text int8` in `Chi_age`
* `Text bool` in `Chi_married`

##Go struct name prefix
`chidley` by default prepends a prefix to Go struct type identifiers. The default is `Chi` but this can be changed with the `-e` flag. If changed from the default, the new prefix must start with a capital letter (for the XML annotation and decoder to work: the struct fields must be public).

##Warning
If you are going to use the `chidley` generated Go structs on XML other than the input XML, you need to make sure the input XML has examples of all tags, and tag attribute and tag child tag combinations. 

If the input does not have all of these, and you use new XML that has tags not found in the input XML, attributes not seen in tags in the input XML, or child tags not encountered in the input XML, these will not be *seen* by the xml decoder, as they will not be in the Go structs used by the xml decoder.

##Limitations
`chidley` is constrained by the underlying Go [xml package](http://golang.org/pkg/encoding/xml/)
Some of these limitations include:
* The default encoding supported by `encoder/xml` is UTF-8. Right now `chidley` does not support additional charsets. 
An xml decoder that handles charsets other than UTF-8 is possible (see example https://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go). 
It is possible that this method might be used in the future to extend `chidley` to include a small set of popular charsets.
* For vanilla XML with no namespaces, there should be no problem using `chidley`

###Go `xml` package Namespace issues
* There are a number of bugs open for the Go xml package that relate to XML namespaces: https://code.google.com/p/go/issues/list?can=2&q=xml+namespace  If the XML you are using uses namespaces in certain ways, these bugs will impact whether `chidley` can create correct structs for your XML
* For _most_ XML with namespaces, the JSON will be OK but if you convert XML to XML using the generated Go code, there will be a chance one of the above mentioned bugs may impact results. Here is an example I encountered: https://groups.google.com/d/msg/golang-nuts/drWStJSt0Pg/Z47JHeij7ToJ

##Name
`chidley` is named after [Cape Chidley](https://en.wikipedia.org/wiki/Cape_Chidley), Canada

##Larger & more complex example
Using the file `xml/pubmed_xml_12750255.xml.bz2`. Generated from a query to pubmed (similar but much larger than [http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=20598978,444444,455555&retmode=xml](http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=20598978,444444,455555&retmode=xml)), returning a document in the [MEDLINE/Pubmed XML format](http://www.nlm.nih.gov/bsd/licensee/data_elements_doc.html).
* Compressed size: 27M
* Uncompressed size: 337M

###Generate Go structs from pubmed_xml_12750255.xml.bz2: `-G`

```
$ /usr/bin/time -f "%E %M" ./chidley -G xml/pubmed_xml_12750255.xml.bz2 
type Chi_root struct {
	Chi_PubmedArticleSet *Chi_PubmedArticleSet `xml:" PubmedArticleSet,omitempty" json:"PubmedArticleSet,omitempty"`
}

type Chi_PubmedArticleSet struct {
	Chi_PubmedArticle []*Chi_PubmedArticle `xml:" PubmedArticle,omitempty" json:"PubmedArticle,omitempty"`
	Chi_PubmedBookArticle []*Chi_PubmedBookArticle `xml:" PubmedBookArticle,omitempty" json:"PubmedBookArticle,omitempty"`
}

type Chi_PubmedArticle struct {
	Chi_MedlineCitation *Chi_MedlineCitation `xml:" MedlineCitation,omitempty" json:"MedlineCitation,omitempty"`
	Chi_PubmedData *Chi_PubmedData `xml:" PubmedData,omitempty" json:"PubmedData,omitempty"`
}

type Chi_MedlineCitation struct {
	Attr_Owner string `xml:" Owner,attr"  json:",omitempty"`
	Attr_Status string `xml:" Status,attr"  json:",omitempty"`
	Attr_VersionDate string `xml:" VersionDate,attr"  json:",omitempty"`
	Attr_VersionID string `xml:" VersionID,attr"  json:",omitempty"`
	Chi_Article *Chi_Article `xml:" Article,omitempty" json:"Article,omitempty"`
	Chi_ChemicalList *Chi_ChemicalList `xml:" ChemicalList,omitempty" json:"ChemicalList,omitempty"`
	Chi_CitationSubset []*Chi_CitationSubset `xml:" CitationSubset,omitempty" json:"CitationSubset,omitempty"`
	Chi_CommentsCorrectionsList *Chi_CommentsCorrectionsList `xml:" CommentsCorrectionsList,omitempty" json:"CommentsCorrectionsList,omitempty"`
	Chi_DateCompleted *Chi_DateCompleted `xml:" DateCompleted,omitempty" json:"DateCompleted,omitempty"`
	Chi_DateCreated *Chi_DateCreated `xml:" DateCreated,omitempty" json:"DateCreated,omitempty"`
	Chi_DateRevised *Chi_DateRevised `xml:" DateRevised,omitempty" json:"DateRevised,omitempty"`
	Chi_GeneSymbolList *Chi_GeneSymbolList `xml:" GeneSymbolList,omitempty" json:"GeneSymbolList,omitempty"`
	Chi_GeneralNote []*Chi_GeneralNote `xml:" GeneralNote,omitempty" json:"GeneralNote,omitempty"`
	Chi_InvestigatorList *Chi_InvestigatorList `xml:" InvestigatorList,omitempty" json:"InvestigatorList,omitempty"`
	Chi_KeywordList *Chi_KeywordList `xml:" KeywordList,omitempty" json:"KeywordList,omitempty"`
	Chi_MedlineJournalInfo *Chi_MedlineJournalInfo `xml:" MedlineJournalInfo,omitempty" json:"MedlineJournalInfo,omitempty"`
	Chi_MeshHeadingList *Chi_MeshHeadingList `xml:" MeshHeadingList,omitempty" json:"MeshHeadingList,omitempty"`
	Chi_NumberOfReferences *Chi_NumberOfReferences `xml:" NumberOfReferences,omitempty" json:"NumberOfReferences,omitempty"`
	Chi_OtherAbstract *Chi_OtherAbstract `xml:" OtherAbstract,omitempty" json:"OtherAbstract,omitempty"`
	Chi_OtherID []*Chi_OtherID `xml:" OtherID,omitempty" json:"OtherID,omitempty"`
	Chi_PMID *Chi_PMID `xml:" PMID,omitempty" json:"PMID,omitempty"`
	Chi_PersonalNameSubjectList *Chi_PersonalNameSubjectList `xml:" PersonalNameSubjectList,omitempty" json:"PersonalNameSubjectList,omitempty"`
	Chi_SpaceFlightMission []*Chi_SpaceFlightMission `xml:" SpaceFlightMission,omitempty" json:"SpaceFlightMission,omitempty"`
	Chi_SupplMeshList *Chi_SupplMeshList `xml:" SupplMeshList,omitempty" json:"SupplMeshList,omitempty"`
}

type Chi_PMID struct {
	Attr_Version string `xml:" Version,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_DateCreated struct {
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_Day struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Year struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Month struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_DateRevised struct {
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_MedlineJournalInfo struct {
	Chi_Country *Chi_Country `xml:" Country,omitempty" json:"Country,omitempty"`
	Chi_ISSNLinking *Chi_ISSNLinking `xml:" ISSNLinking,omitempty" json:"ISSNLinking,omitempty"`
	Chi_MedlineTA *Chi_MedlineTA `xml:" MedlineTA,omitempty" json:"MedlineTA,omitempty"`
	Chi_NlmUniqueID *Chi_NlmUniqueID `xml:" NlmUniqueID,omitempty" json:"NlmUniqueID,omitempty"`
}

type Chi_Country struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_MedlineTA struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_NlmUniqueID struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ISSNLinking struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_CitationSubset struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_NumberOfReferences struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_KeywordList struct {
	Attr_Owner string `xml:" Owner,attr"  json:",omitempty"`
	Chi_Keyword []*Chi_Keyword `xml:" Keyword,omitempty" json:"Keyword,omitempty"`
}

type Chi_Keyword struct {
	Attr_MajorTopicYN string `xml:" MajorTopicYN,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_OtherAbstract struct {
	Attr_Type string `xml:" Type,attr"  json:",omitempty"`
	Attr_Language string `xml:" Language,attr"  json:",omitempty"`
	Chi_AbstractText []*Chi_AbstractText `xml:" AbstractText,omitempty" json:"AbstractText,omitempty"`
	Chi_CopyrightInformation *Chi_CopyrightInformation `xml:" CopyrightInformation,omitempty" json:"CopyrightInformation,omitempty"`
}

type Chi_AbstractText struct {
	Attr_Label string `xml:" Label,attr"  json:",omitempty"`
	Attr_NlmCategory string `xml:" NlmCategory,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_CopyrightInformation struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PersonalNameSubjectList struct {
	Chi_PersonalNameSubject []*Chi_PersonalNameSubject `xml:" PersonalNameSubject,omitempty" json:"PersonalNameSubject,omitempty"`
}

type Chi_PersonalNameSubject struct {
	Chi_ForeName *Chi_ForeName `xml:" ForeName,omitempty" json:"ForeName,omitempty"`
	Chi_Initials *Chi_Initials `xml:" Initials,omitempty" json:"Initials,omitempty"`
	Chi_LastName *Chi_LastName `xml:" LastName,omitempty" json:"LastName,omitempty"`
	Chi_Suffix *Chi_Suffix `xml:" Suffix,omitempty" json:"Suffix,omitempty"`
}

type Chi_LastName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ForeName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Initials struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Suffix struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_GeneSymbolList struct {
	Chi_GeneSymbol []*Chi_GeneSymbol `xml:" GeneSymbol,omitempty" json:"GeneSymbol,omitempty"`
}

type Chi_GeneSymbol struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_SupplMeshList struct {
	Chi_SupplMeshName []*Chi_SupplMeshName `xml:" SupplMeshName,omitempty" json:"SupplMeshName,omitempty"`
}

type Chi_SupplMeshName struct {
	Attr_Type string `xml:" Type,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_DateCompleted struct {
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_Article struct {
	Attr_PubModel string `xml:" PubModel,attr"  json:",omitempty"`
	Chi_Abstract *Chi_Abstract `xml:" Abstract,omitempty" json:"Abstract,omitempty"`
	Chi_ArticleDate *Chi_ArticleDate `xml:" ArticleDate,omitempty" json:"ArticleDate,omitempty"`
	Chi_ArticleTitle *Chi_ArticleTitle `xml:" ArticleTitle,omitempty" json:"ArticleTitle,omitempty"`
	Chi_AuthorList []*Chi_AuthorList `xml:" AuthorList,omitempty" json:"AuthorList,omitempty"`
	Chi_DataBankList *Chi_DataBankList `xml:" DataBankList,omitempty" json:"DataBankList,omitempty"`
	Chi_ELocationID []*Chi_ELocationID `xml:" ELocationID,omitempty" json:"ELocationID,omitempty"`
	Chi_GrantList *Chi_GrantList `xml:" GrantList,omitempty" json:"GrantList,omitempty"`
	Chi_Journal *Chi_Journal `xml:" Journal,omitempty" json:"Journal,omitempty"`
	Chi_Language []*Chi_Language `xml:" Language,omitempty" json:"Language,omitempty"`
	Chi_Pagination *Chi_Pagination `xml:" Pagination,omitempty" json:"Pagination,omitempty"`
	Chi_PublicationTypeList *Chi_PublicationTypeList `xml:" PublicationTypeList,omitempty" json:"PublicationTypeList,omitempty"`
	Chi_VernacularTitle *Chi_VernacularTitle `xml:" VernacularTitle,omitempty" json:"VernacularTitle,omitempty"`
}

type Chi_Pagination struct {
	Chi_MedlinePgn *Chi_MedlinePgn `xml:" MedlinePgn,omitempty" json:"MedlinePgn,omitempty"`
}

type Chi_MedlinePgn struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_AuthorList struct {
	Attr_CompleteYN string `xml:" CompleteYN,attr"  json:",omitempty"`
	Attr_Type string `xml:" Type,attr"  json:",omitempty"`
	Chi_Author []*Chi_Author `xml:" Author,omitempty" json:"Author,omitempty"`
}

type Chi_Author struct {
	Attr_ValidYN string `xml:" ValidYN,attr"  json:",omitempty"`
	Chi_Affiliation *Chi_Affiliation `xml:" Affiliation,omitempty" json:"Affiliation,omitempty"`
	Chi_CollectiveName *Chi_CollectiveName `xml:" CollectiveName,omitempty" json:"CollectiveName,omitempty"`
	Chi_ForeName *Chi_ForeName `xml:" ForeName,omitempty" json:"ForeName,omitempty"`
	Chi_Identifier *Chi_Identifier `xml:" Identifier,omitempty" json:"Identifier,omitempty"`
	Chi_Initials *Chi_Initials `xml:" Initials,omitempty" json:"Initials,omitempty"`
	Chi_LastName *Chi_LastName `xml:" LastName,omitempty" json:"LastName,omitempty"`
	Chi_Suffix *Chi_Suffix `xml:" Suffix,omitempty" json:"Suffix,omitempty"`
}

type Chi_Affiliation struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_CollectiveName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Identifier struct {
	Attr_Source string `xml:" Source,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Language struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Journal struct {
	Chi_ISOAbbreviation *Chi_ISOAbbreviation `xml:" ISOAbbreviation,omitempty" json:"ISOAbbreviation,omitempty"`
	Chi_ISSN *Chi_ISSN `xml:" ISSN,omitempty" json:"ISSN,omitempty"`
	Chi_JournalIssue *Chi_JournalIssue `xml:" JournalIssue,omitempty" json:"JournalIssue,omitempty"`
	Chi_Title *Chi_Title `xml:" Title,omitempty" json:"Title,omitempty"`
}

type Chi_ISSN struct {
	Attr_IssnType string `xml:" IssnType,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_JournalIssue struct {
	Attr_CitedMedium string `xml:" CitedMedium,attr"  json:",omitempty"`
	Chi_Issue *Chi_Issue `xml:" Issue,omitempty" json:"Issue,omitempty"`
	Chi_PubDate *Chi_PubDate `xml:" PubDate,omitempty" json:"PubDate,omitempty"`
	Chi_Volume *Chi_Volume `xml:" Volume,omitempty" json:"Volume,omitempty"`
}

type Chi_Volume struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Issue struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PubDate struct {
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_MedlineDate *Chi_MedlineDate `xml:" MedlineDate,omitempty" json:"MedlineDate,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Season *Chi_Season `xml:" Season,omitempty" json:"Season,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_MedlineDate struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Season struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Title struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ISOAbbreviation struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ArticleTitle struct {
	Attr_book string `xml:" book,attr"  json:",omitempty"`
	Attr_part string `xml:" part,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Abstract struct {
	Chi_AbstractText []*Chi_AbstractText `xml:" AbstractText,omitempty" json:"AbstractText,omitempty"`
	Chi_CopyrightInformation *Chi_CopyrightInformation `xml:" CopyrightInformation,omitempty" json:"CopyrightInformation,omitempty"`
}

type Chi_PublicationTypeList struct {
	Chi_PublicationType []*Chi_PublicationType `xml:" PublicationType,omitempty" json:"PublicationType,omitempty"`
}

type Chi_PublicationType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_VernacularTitle struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ELocationID struct {
	Attr_EIdType string `xml:" EIdType,attr"  json:",omitempty"`
	Attr_ValidYN string `xml:" ValidYN,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_DataBankList struct {
	Attr_CompleteYN string `xml:" CompleteYN,attr"  json:",omitempty"`
	Chi_DataBank []*Chi_DataBank `xml:" DataBank,omitempty" json:"DataBank,omitempty"`
}

type Chi_DataBank struct {
	Chi_AccessionNumberList *Chi_AccessionNumberList `xml:" AccessionNumberList,omitempty" json:"AccessionNumberList,omitempty"`
	Chi_DataBankName *Chi_DataBankName `xml:" DataBankName,omitempty" json:"DataBankName,omitempty"`
}

type Chi_DataBankName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_AccessionNumberList struct {
	Chi_AccessionNumber []*Chi_AccessionNumber `xml:" AccessionNumber,omitempty" json:"AccessionNumber,omitempty"`
}

type Chi_AccessionNumber struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ArticleDate struct {
	Attr_DateType string `xml:" DateType,attr"  json:",omitempty"`
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_GrantList struct {
	Attr_CompleteYN string `xml:" CompleteYN,attr"  json:",omitempty"`
	Chi_Grant []*Chi_Grant `xml:" Grant,omitempty" json:"Grant,omitempty"`
}

type Chi_Grant struct {
	Chi_Acronym *Chi_Acronym `xml:" Acronym,omitempty" json:"Acronym,omitempty"`
	Chi_Agency *Chi_Agency `xml:" Agency,omitempty" json:"Agency,omitempty"`
	Chi_Country *Chi_Country `xml:" Country,omitempty" json:"Country,omitempty"`
	Chi_GrantID *Chi_GrantID `xml:" GrantID,omitempty" json:"GrantID,omitempty"`
}

type Chi_GrantID struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Acronym struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Agency struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_MeshHeadingList struct {
	Chi_MeshHeading []*Chi_MeshHeading `xml:" MeshHeading,omitempty" json:"MeshHeading,omitempty"`
}

type Chi_MeshHeading struct {
	Chi_DescriptorName *Chi_DescriptorName `xml:" DescriptorName,omitempty" json:"DescriptorName,omitempty"`
	Chi_QualifierName []*Chi_QualifierName `xml:" QualifierName,omitempty" json:"QualifierName,omitempty"`
}

type Chi_DescriptorName struct {
	Attr_MajorTopicYN string `xml:" MajorTopicYN,attr"  json:",omitempty"`
	Attr_Type string `xml:" Type,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_QualifierName struct {
	Attr_MajorTopicYN string `xml:" MajorTopicYN,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ChemicalList struct {
	Chi_Chemical []*Chi_Chemical `xml:" Chemical,omitempty" json:"Chemical,omitempty"`
}

type Chi_Chemical struct {
	Chi_NameOfSubstance *Chi_NameOfSubstance `xml:" NameOfSubstance,omitempty" json:"NameOfSubstance,omitempty"`
	Chi_RegistryNumber *Chi_RegistryNumber `xml:" RegistryNumber,omitempty" json:"RegistryNumber,omitempty"`
}

type Chi_RegistryNumber struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_NameOfSubstance struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_CommentsCorrectionsList struct {
	Chi_CommentsCorrections []*Chi_CommentsCorrections `xml:" CommentsCorrections,omitempty" json:"CommentsCorrections,omitempty"`
}

type Chi_CommentsCorrections struct {
	Attr_RefType string `xml:" RefType,attr"  json:",omitempty"`
	Chi_Note *Chi_Note `xml:" Note,omitempty" json:"Note,omitempty"`
	Chi_PMID *Chi_PMID `xml:" PMID,omitempty" json:"PMID,omitempty"`
	Chi_RefSource *Chi_RefSource `xml:" RefSource,omitempty" json:"RefSource,omitempty"`
}

type Chi_RefSource struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Note struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_OtherID struct {
	Attr_Source string `xml:" Source,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_GeneralNote struct {
	Attr_Owner string `xml:" Owner,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_InvestigatorList struct {
	Chi_Investigator []*Chi_Investigator `xml:" Investigator,omitempty" json:"Investigator,omitempty"`
}

type Chi_Investigator struct {
	Attr_ValidYN string `xml:" ValidYN,attr"  json:",omitempty"`
	Chi_Affiliation *Chi_Affiliation `xml:" Affiliation,omitempty" json:"Affiliation,omitempty"`
	Chi_ForeName *Chi_ForeName `xml:" ForeName,omitempty" json:"ForeName,omitempty"`
	Chi_Initials *Chi_Initials `xml:" Initials,omitempty" json:"Initials,omitempty"`
	Chi_LastName *Chi_LastName `xml:" LastName,omitempty" json:"LastName,omitempty"`
	Chi_Suffix *Chi_Suffix `xml:" Suffix,omitempty" json:"Suffix,omitempty"`
}

type Chi_SpaceFlightMission struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PubmedData struct {
	Chi_ArticleIdList *Chi_ArticleIdList `xml:" ArticleIdList,omitempty" json:"ArticleIdList,omitempty"`
	Chi_History *Chi_History `xml:" History,omitempty" json:"History,omitempty"`
	Chi_PublicationStatus *Chi_PublicationStatus `xml:" PublicationStatus,omitempty" json:"PublicationStatus,omitempty"`
}

type Chi_History struct {
	Chi_PubMedPubDate []*Chi_PubMedPubDate `xml:" PubMedPubDate,omitempty" json:"PubMedPubDate,omitempty"`
}

type Chi_PubMedPubDate struct {
	Attr_PubStatus string `xml:" PubStatus,attr"  json:",omitempty"`
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_Hour *Chi_Hour `xml:" Hour,omitempty" json:"Hour,omitempty"`
	Chi_Minute *Chi_Minute `xml:" Minute,omitempty" json:"Minute,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_Hour struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Minute struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PublicationStatus struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ArticleIdList struct {
	Chi_ArticleId []*Chi_ArticleId `xml:" ArticleId,omitempty" json:"ArticleId,omitempty"`
}

type Chi_ArticleId struct {
	Attr_IdType string `xml:" IdType,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PubmedBookArticle struct {
	Chi_BookDocument *Chi_BookDocument `xml:" BookDocument,omitempty" json:"BookDocument,omitempty"`
	Chi_PubmedBookData *Chi_PubmedBookData `xml:" PubmedBookData,omitempty" json:"PubmedBookData,omitempty"`
}

type Chi_BookDocument struct {
	Chi_Abstract *Chi_Abstract `xml:" Abstract,omitempty" json:"Abstract,omitempty"`
	Chi_ArticleIdList *Chi_ArticleIdList `xml:" ArticleIdList,omitempty" json:"ArticleIdList,omitempty"`
	Chi_ArticleTitle *Chi_ArticleTitle `xml:" ArticleTitle,omitempty" json:"ArticleTitle,omitempty"`
	Chi_AuthorList []*Chi_AuthorList `xml:" AuthorList,omitempty" json:"AuthorList,omitempty"`
	Chi_Book *Chi_Book `xml:" Book,omitempty" json:"Book,omitempty"`
	Chi_ContributionDate *Chi_ContributionDate `xml:" ContributionDate,omitempty" json:"ContributionDate,omitempty"`
	Chi_DateRevised *Chi_DateRevised `xml:" DateRevised,omitempty" json:"DateRevised,omitempty"`
	Chi_ItemList []*Chi_ItemList `xml:" ItemList,omitempty" json:"ItemList,omitempty"`
	Chi_KeywordList *Chi_KeywordList `xml:" KeywordList,omitempty" json:"KeywordList,omitempty"`
	Chi_Language []*Chi_Language `xml:" Language,omitempty" json:"Language,omitempty"`
	Chi_LocationLabel []*Chi_LocationLabel `xml:" LocationLabel,omitempty" json:"LocationLabel,omitempty"`
	Chi_PMID *Chi_PMID `xml:" PMID,omitempty" json:"PMID,omitempty"`
	Chi_Sections *Chi_Sections `xml:" Sections,omitempty" json:"Sections,omitempty"`
}

type Chi_LocationLabel struct {
	Attr_Type string `xml:" Type,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Sections struct {
	Chi_Section []*Chi_Section `xml:" Section,omitempty" json:"Section,omitempty"`
}

type Chi_Section struct {
	Chi_LocationLabel []*Chi_LocationLabel `xml:" LocationLabel,omitempty" json:"LocationLabel,omitempty"`
	Chi_Section []*Chi_Section `xml:" Section,omitempty" json:"Section,omitempty"`
	Chi_SectionTitle []*Chi_SectionTitle `xml:" SectionTitle,omitempty" json:"SectionTitle,omitempty"`
}

type Chi_SectionTitle struct {
	Attr_book string `xml:" book,attr"  json:",omitempty"`
	Attr_part string `xml:" part,attr"  json:",omitempty"`
	Attr_sec string `xml:" sec,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Book struct {
	Chi_AuthorList []*Chi_AuthorList `xml:" AuthorList,omitempty" json:"AuthorList,omitempty"`
	Chi_BeginningDate *Chi_BeginningDate `xml:" BeginningDate,omitempty" json:"BeginningDate,omitempty"`
	Chi_BookTitle *Chi_BookTitle `xml:" BookTitle,omitempty" json:"BookTitle,omitempty"`
	Chi_CollectionTitle *Chi_CollectionTitle `xml:" CollectionTitle,omitempty" json:"CollectionTitle,omitempty"`
	Chi_Edition *Chi_Edition `xml:" Edition,omitempty" json:"Edition,omitempty"`
	Chi_EndingDate *Chi_EndingDate `xml:" EndingDate,omitempty" json:"EndingDate,omitempty"`
	Chi_Isbn []*Chi_Isbn `xml:" Isbn,omitempty" json:"Isbn,omitempty"`
	Chi_Medium *Chi_Medium `xml:" Medium,omitempty" json:"Medium,omitempty"`
	Chi_PubDate *Chi_PubDate `xml:" PubDate,omitempty" json:"PubDate,omitempty"`
	Chi_Publisher *Chi_Publisher `xml:" Publisher,omitempty" json:"Publisher,omitempty"`
}

type Chi_BookTitle struct {
	Attr_book string `xml:" book,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_BeginningDate struct {
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_Edition struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Publisher struct {
	Chi_PublisherLocation *Chi_PublisherLocation `xml:" PublisherLocation,omitempty" json:"PublisherLocation,omitempty"`
	Chi_PublisherName *Chi_PublisherName `xml:" PublisherName,omitempty" json:"PublisherName,omitempty"`
}

type Chi_PublisherName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PublisherLocation struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_EndingDate struct {
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_Medium struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_Isbn struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_CollectionTitle struct {
	Attr_book string `xml:" book,attr"  json:",omitempty"`
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_ContributionDate struct {
	Chi_Day *Chi_Day `xml:" Day,omitempty" json:"Day,omitempty"`
	Chi_Month *Chi_Month `xml:" Month,omitempty" json:"Month,omitempty"`
	Chi_Year *Chi_Year `xml:" Year,omitempty" json:"Year,omitempty"`
}

type Chi_ItemList struct {
	Attr_ListType string `xml:" ListType,attr"  json:",omitempty"`
	Chi_Item *Chi_Item `xml:" Item,omitempty" json:"Item,omitempty"`
}

type Chi_Item struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Chi_PubmedBookData struct {
	Chi_ArticleIdList *Chi_ArticleIdList `xml:" ArticleIdList,omitempty" json:"ArticleIdList,omitempty"`
	Chi_History *Chi_History `xml:" History,omitempty" json:"History,omitempty"`
	Chi_PublicationStatus *Chi_PublicationStatus `xml:" PublicationStatus,omitempty" json:"PublicationStatus,omitempty"`
}

0:48.21 15044
$
```

48 seconds for 337MB XML; resident size: 15MB

*Note:* All timings from Dell laptop 16GB, regular disk, 8 core i7-3720QM CPU @ 2.60GHz)

`Linux 3.11.10-100.fc18.x86_64 #1 SMP Mon Dec 2 20:28:38 UTC 2013 x86_64 x86_64 x86_64 GNU/Linux`



###Generate Go program: `-W`

```
$ /usr/bin/time -f "%E %M" ./chidley -W xml/pubmed_xml_12750255.xml.bz2 > examples/pubmed/ChiPubmed.go
0:47.75 17064
$ cd examples/pubmed/
$ go build
$
```
48 seconds for 337MB XML; resident size: 17MB

####Generated program: count tags
```
$ /usr/bin/time -f "%E %M" ./pubmed -c | sort -n
0:36.58 10460
1 _:Edition
1 _:Identifier
1 _:PubmedArticleSet
2 _:EndingDate
2 _:Item
2 _:ItemList
3 _:BeginningDate
3 _:CollectionTitle
3 _:ContributionDate
3 _:Medium
7 _:Book
7 _:BookDocument
7 _:BookTitle
7 _:Isbn
7 _:Publisher
7 _:PublisherLocation
7 _:PublisherName
7 _:PubmedBookArticle
7 _:PubmedBookData
7 _:Sections
19 _:SpaceFlightMission
33 _:Note
38 _:LocationLabel
69 _:InvestigatorList
78 _:SupplMeshList
82 _:Section
82 _:SectionTitle
87 _:SupplMeshName
101 _:GeneSymbolList
128 _:OtherAbstract
214 _:GeneSymbol
285 _:CollectiveName
366 _:PersonalNameSubjectList
419 _:PersonalNameSubject
424 _:Season
491 _:DataBankList
499 _:AccessionNumberList
499 _:DataBank
499 _:DataBankName
615 _:GeneralNote
861 _:Suffix
2532 _:Investigator
2618 _:CopyrightInformation
3006 _:NumberOfReferences
3730 _:MedlineDate
3830 _:AccessionNumber
4399 _:GrantList
5233 _:KeywordList
6201 _:CommentsCorrectionsList
7190 _:VernacularTitle
7783 _:ArticleDate
8687 _:OtherID
9324 _:Acronym
9805 _:ELocationID
9864 _:GrantID
10120 _:Agency
10120 _:Grant
18050 _:Keyword
23156 _:ChemicalList
28079 _:Affiliation
28656 _:Abstract
38168 _:DateRevised
43353 _:MeshHeadingList
43405 _:Issue
45839 _:AbstractText
46216 _:DateCompleted
46926 _:ISSNLinking
46998 _:Volume
47031 _:AuthorList
47376 _:ISSN
48000 _:Article
48000 _:DateCreated
48000 _:ISOAbbreviation
48000 _:Journal
48000 _:JournalIssue
48000 _:MedlineCitation
48000 _:MedlineJournalInfo
48000 _:MedlinePgn
48000 _:MedlineTA
48000 _:NlmUniqueID
48000 _:Pagination
48000 _:PublicationTypeList
48000 _:PubmedArticle
48000 _:PubmedData
48000 _:Title
48005 _:ArticleTitle
48007 _:History
48007 _:PubDate
48007 _:PublicationStatus
48014 _:ArticleIdList
48107 _:Language
50862 _:CitationSubset
57755 _:Country
82258 _:ArticleId
85217 _:PublicationType
87553 _:Chemical
87553 _:NameOfSubstance
87553 _:RegistryNumber
119444 _:Hour
119444 _:Minute
123069 _:CommentsCorrections
123069 _:RefSource
165664 _:PubMedPubDate
170077 _:Author
170804 _:PMID
171242 _:ForeName
172469 _:Initials
172743 _:LastName
241729 _:QualifierName
315764 _:Day
343262 _:Month
350116 _:Year
439020 _:DescriptorName
439020 _:MeshHeading
$
```

*Note:* The underscore before the colon indicates there is no (or the default) namespace for the element.

36 seconds for 337MB XML; resident size: 10.5MB

####Generated program: convert XML to JSON

#####No streaming
```
$ /usr/bin/time -f "%E %M" ./pubmed -j > /dev/null
0:57.26 2866408
$
```
57 seconds for 337MB XML; resident size: 2.9GB

#####With streaming `-s`
Streaming decodes using the XML elements that are one level down from the top level container element.
```
$ /usr/bin/time -f "%E %M" ./pubmed -j -s > /dev/null
0:58.72 15944
```
59 seconds for 337MB XML; resident size: 16MB

#####Sample of generated JSON
```
$ /usr/bin/time -f "%E %M" ./pubmed -j -s |head -310
{
 "MedlineCitation": {
  "Attr_Owner": "NLM",
  "Attr_Status": "MEDLINE",
  "Article": {
   "Attr_PubModel": "Print",
   "Abstract": {
    "AbstractText": [
     {
      "Text": "A review on the operative methods for prophylaxis of urological complications (fistulas and strictures) due to radical hysterectomy with systemic dissection of lymph nodes is described. The authors recommend the method of T. H. Green as the most effective method. A new method for protection of the ureter with flaps, formed as a \"leg\" from omentum majus, is proposed. The modification has been used in 20 patients without postoperative complications. The method is recommended in cases, when postoperative stenosis or strictures of the ureters are expected as well as when postoperative irradiation is forthcoming."
     }
    ]
   },
   "ArticleTitle": {
    "Text": "[A method for preventing the urologic complications connected with the surgical treatment of cancer of the cervix uteri]."
   },
   "AuthorList": [
    {
     "Attr_CompleteYN": "Y",
     "Author": [
      {
       "Attr_ValidYN": "Y",
       "ForeName": {
        "Text": "T"
       },
       "Initials": {
        "Text": "T"
       },
       "LastName": {
        "Text": "Kŭrlov"
       }
      },
      {
       "Attr_ValidYN": "Y",
       "ForeName": {
        "Text": "N"
       },
       "Initials": {
        "Text": "N"
       },
       "LastName": {
        "Text": "Vasilev"
       }
      }
     ]
    }
   ],
   "Journal": {
    "ISOAbbreviation": {
     "Text": "Akush Ginekol (Sofiia)"
    },
    "ISSN": {
     "Attr_IssnType": "Print",
     "Text": "0324-0959"
    },
    "JournalIssue": {
     "Attr_CitedMedium": "Print",
     "Issue": {
      "Text": "1"
     },
     "PubDate": {
      "Year": {
       "Text": "1990"
      }
     },
     "Volume": {
      "Text": "29"
     }
    },
    "Title": {
     "Text": "Akusherstvo i ginekologii͡a"
    }
   },
   "Language": [
    {
     "Text": "bul"
    }
   ],
   "Pagination": {
    "MedlinePgn": {
     "Text": "55-7"
    }
   },
   "PublicationTypeList": {
    "PublicationType": [
     {
      "Text": "English Abstract"
     },
     {
      "Text": "Journal Article"
     }
    ]
   },
   "VernacularTitle": {
    "Text": "Metod za profilaktika na urologichnite uslozhneniia, svŭrzani s operativnoto lechenie na raka na matochnata shiĭka."
   }
  },
  "CitationSubset": [
   {
    "Text": "IM"
   }
  ],
  "DateCompleted": {
   "Day": {
    "Text": "22"
   },
   "Month": {
    "Text": "08"
   },
   "Year": {
    "Text": "1990"
   }
  },
  "DateCreated": {
   "Day": {
    "Text": "22"
   },
   "Month": {
    "Text": "08"
   },
   "Year": {
    "Text": "1990"
   }
  },
  "DateRevised": {
   "Day": {
    "Text": "15"
   },
   "Month": {
    "Text": "11"
   },
   "Year": {
    "Text": "2006"
   }
  },
  "MedlineJournalInfo": {
   "Country": {
    "Text": "BULGARIA"
   },
   "ISSNLinking": {
    "Text": "0324-0959"
   },
   "MedlineTA": {
    "Text": "Akush Ginekol (Sofiia)"
   },
   "NlmUniqueID": {
    "Text": "0370455"
   }
  },
  "MeshHeadingList": {
   "MeshHeading": [
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Female"
     }
    },
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Humans"
     }
    },
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Hysterectomy"
     },
     "QualifierName": [
      {
       "Attr_MajorTopicYN": "N",
       "Text": "methods"
      }
     ]
    },
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Lymph Node Excision"
     },
     "QualifierName": [
      {
       "Attr_MajorTopicYN": "N",
       "Text": "methods"
      }
     ]
    },
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Postoperative Complications"
     },
     "QualifierName": [
      {
       "Attr_MajorTopicYN": "N",
       "Text": "etiology"
      },
      {
       "Attr_MajorTopicYN": "Y",
       "Text": "prevention \u0026 control"
      }
     ]
    },
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Urologic Diseases"
     },
     "QualifierName": [
      {
       "Attr_MajorTopicYN": "N",
       "Text": "etiology"
      },
      {
       "Attr_MajorTopicYN": "Y",
       "Text": "prevention \u0026 control"
      }
     ]
    },
    {
     "DescriptorName": {
      "Attr_MajorTopicYN": "N",
      "Text": "Uterine Cervical Neoplasms"
     },
     "QualifierName": [
      {
       "Attr_MajorTopicYN": "N",
       "Text": "complications"
      },
      {
       "Attr_MajorTopicYN": "Y",
       "Text": "surgery"
      }
     ]
    }
   ]
  },
  "PMID": {
   "Attr_Version": "1",
   "Text": "2372101"
  }
 },
 "PubmedData": {
  "ArticleIdList": {
   "ArticleId": [
    {
     "Attr_IdType": "pubmed",
     "Text": "2372101"
    }
   ]
  },
  "History": {
   "PubMedPubDate": [
    {
     "Attr_PubStatus": "pubmed",
     "Day": {
      "Text": "1"
     },
     "Month": {
      "Text": "1"
     },
     "Year": {
      "Text": "1990"
     }
    },
    {
     "Attr_PubStatus": "medline",
     "Day": {
      "Text": "1"
     },
     "Hour": {
      "Text": "0"
     },
     "Minute": {
      "Text": "1"
     },
     "Month": {
      "Text": "1"
     },
     "Year": {
      "Text": "1990"
     }
    },
    {
     "Attr_PubStatus": "entrez",
     "Day": {
      "Text": "1"
     },
     "Hour": {
      "Text": "0"
     },
     "Minute": {
      "Text": "0"
     },
     "Month": {
      "Text": "1"
     },
     "Year": {
      "Text": "1990"
     }
    }
   ]
  },
  "PublicationStatus": {
   "Text": "ppublish"
  }
 }
}
{
 "MedlineCitation": {
```
 
#####Sample of generated XML to XML
```
$ ./pubmed -x -s |head -100
  <Chi_PubmedArticle>
      <MedlineCitation Owner="NLM" Status="MEDLINE" VersionDate="" VersionID="">
          <Article PubModel="Print">
              <Abstract>
                  <AbstractText Label="" NlmCategory="">A review on the operative methods for prophylaxis of urological complications (fistulas and strictures) due to radical hysterectomy with systemic dissection of lymph nodes is described. The authors recommend the method of T. H. Green as the most effective method. A new method for protection of the ureter with flaps, formed as a &#34;leg&#34; from omentum majus, is proposed. The modification has been used in 20 patients without postoperative complications. The method is recommended in cases, when postoperative stenosis or strictures of the ureters are expected as well as when postoperative irradiation is forthcoming.</AbstractText>
              </Abstract>
              <ArticleTitle>[A method for preventing the urologic complications connected with the surgical treatment of cancer of the cervix uteri].</ArticleTitle>
              <AuthorList CompleteYN="Y" Type="">
                  <Author ValidYN="Y">
                      <ForeName>T</ForeName>
                      <Initials>T</Initials>
                      <LastName>Kŭrlov</LastName>
                  </Author>
                  <Author ValidYN="Y">
                      <ForeName>N</ForeName>
                      <Initials>N</Initials>
                      <LastName>Vasilev</LastName>
                  </Author>
              </AuthorList>
              <Journal>
                  <ISOAbbreviation>Akush Ginekol (Sofiia)</ISOAbbreviation>
                  <ISSN IssnType="Print">0324-0959</ISSN>
                  <JournalIssue CitedMedium="Print">
                      <Issue>1</Issue>
                      <PubDate>
                          <Year>1990</Year>
                      </PubDate>
                      <Volume>29</Volume>
                  </JournalIssue>
                  <Title>Akusherstvo i ginekologii͡a</Title>
              </Journal>
              <Language>bul</Language>
              <Pagination>
                  <MedlinePgn>55-7</MedlinePgn>
              </Pagination>
              <PublicationTypeList>
                  <PublicationType>English Abstract</PublicationType>
                  <PublicationType>Journal Article</PublicationType>
              </PublicationTypeList>
              <VernacularTitle>Metod za profilaktika na urologichnite uslozhneniia, svŭrzani s operativnoto lechenie na raka na matochnata shiĭka.</VernacularTitle>
          </Article>
          <CitationSubset>IM</CitationSubset>
          <DateCompleted>
              <Day>22</Day>
              <Month>08</Month>
              <Year>1990</Year>
          </DateCompleted>
          <DateCreated>
              <Day>22</Day>
              <Month>08</Month>
              <Year>1990</Year>
          </DateCreated>
          <DateRevised>
              <Day>15</Day>
              <Month>11</Month>
              <Year>2006</Year>
          </DateRevised>
          <MedlineJournalInfo>
              <Country>BULGARIA</Country>
              <ISSNLinking>0324-0959</ISSNLinking>
              <MedlineTA>Akush Ginekol (Sofiia)</MedlineTA>
              <NlmUniqueID>0370455</NlmUniqueID>
          </MedlineJournalInfo>
          <MeshHeadingList>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Female</DescriptorName>
              </MeshHeading>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Humans</DescriptorName>
              </MeshHeading>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Hysterectomy</DescriptorName>
                  <QualifierName MajorTopicYN="N">methods</QualifierName>
              </MeshHeading>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Lymph Node Excision</DescriptorName>
                  <QualifierName MajorTopicYN="N">methods</QualifierName>
              </MeshHeading>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Postoperative Complications</DescriptorName>
                  <QualifierName MajorTopicYN="N">etiology</QualifierName>
                  <QualifierName MajorTopicYN="Y">prevention &amp; control</QualifierName>
              </MeshHeading>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Urologic Diseases</DescriptorName>
                  <QualifierName MajorTopicYN="N">etiology</QualifierName>
                  <QualifierName MajorTopicYN="Y">prevention &amp; control</QualifierName>
              </MeshHeading>
              <MeshHeading>
                  <DescriptorName MajorTopicYN="N" Type="">Uterine Cervical Neoplasms</DescriptorName>
                  <QualifierName MajorTopicYN="N">complications</QualifierName>
                  <QualifierName MajorTopicYN="Y">surgery</QualifierName>
              </MeshHeading>
          </MeshHeadingList>
          <PMID Version="1">2372101</PMID>
      </MedlineCitation>
      <PubmedData>
          <ArticleIdList>
              <ArticleId IdType="pubmed">2372101</ArticleId>
          </ArticleIdList>
```

## <a name="java">Java/JAXB</a>
`chidley` now supports the production of Java/JAXB code. It generates a class-per-element, with classes mapping to the Go structs generated by the XML extraction.
Its only dependency is Google [Gson](https://code.google.com/p/google-gson/), for JSON generation.


###Usage
`chidley` creates a maven project in `./java` (settable using the `-D` flag) and creates Java JAXB files in `src/main/java/ca/gnewton/chidley/jaxb/xml`. 
It creates a `Main.java` in `src/main/java/ca/gnewton/chidley/jaxb`
```
$ chidley -J xml/test1.xml
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiDocs.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiDoc.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiTitle.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiAuthor.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiLast_name.java
2014/09/02 10:22:28 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiFirstName.java
2014/09/02 10:22:28 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/Main.java
```

###New
Changing the package with "-P":
```
$ chidley -J -P testFoo xml/test1.xml
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiDocs.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiDoc.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiTitle.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiAuthor.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiLast_name.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiFirstName.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/Main.java
$
```

####Build Java package
```
$ cd java
$ ls
pom.xml  src
$ $ mvn package
[INFO] Scanning for projects...
[INFO]                                                                         
[INFO] ------------------------------------------------------------------------
[INFO] Building chidley-jaxb 1.0-SNAPSHOT
[INFO] ------------------------------------------------------------------------
[INFO] 
[INFO] --- maven-resources-plugin:2.6:resources (default-resources) @ chidley-jaxb ---
[WARNING] Using platform encoding (UTF-8 actually) to copy filtered resources, i.e. build is platform dependent!
[INFO] skip non existing resourceDirectory /home/newtong/work/chidley/java/src/main/resources
[INFO] 
[INFO] --- maven-compiler-plugin:2.5.1:compile (default-compile) @ chidley-jaxb ---
[WARNING] File encoding has not been set, using platform encoding UTF-8, i.e. build is platform dependent!
[INFO] Compiling 7 source files to /home/newtong/work/chidley/java/target/classes
[INFO] 
[INFO] --- maven-resources-plugin:2.6:testResources (default-testResources) @ chidley-jaxb ---
[WARNING] Using platform encoding (UTF-8 actually) to copy filtered resources, i.e. build is platform dependent!
[INFO] skip non existing resourceDirectory /home/newtong/work/chidley/java/src/test/resources
[INFO] 
[INFO] --- maven-compiler-plugin:2.5.1:testCompile (default-testCompile) @ chidley-jaxb ---
[INFO] No sources to compile
[INFO] 
[INFO] --- maven-surefire-plugin:2.12.4:test (default-test) @ chidley-jaxb ---
[INFO] No tests to run.
[INFO] 
[INFO] --- maven-jar-plugin:2.4:jar (default-jar) @ chidley-jaxb ---
[INFO] Building jar: /home/newtong/work/chidley/java/target/chidley-jaxb-1.0-SNAPSHOT.jar
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time: 3.514s
[INFO] Finished at: Tue Sep 02 10:29:18 EDT 2014
[INFO] Final Memory: 11M/240M
[INFO] ------------------------------------------------------------------------
$
``` 

#### Running
```
$ export CLASSPATH=./target/chidley-jaxb-1.0-SNAPSHOT.jar:/home/myhome/.m2/repository/com/google/code/gson/gson/2.3/gson-2.3.jar:$CLASSPATH
$ java ca.gnewton.chidley.jaxb.Main
{
  "language": "eng",
  "doc": [
    {
      "type": "book",
      "title": {
        "tagValue": "Dune"
      },
      "author": {
        "last-name": {
          "tagValue": "Herbert"
        },
        "firstName": {
          "tagValue": "Frank"
        }
      }
    },
    {
      "type": "article",
      "title": {
        "tagValue": "Brave New Wold"
      },
      "author": {
        "last-name": {
          "tagValue": "Huxley"
        },
        "firstName": {
          "tagValue": "Aldous"
        }
      }
    }
  ]
}
$ 
```

###Limitations
- Can handle vanilla XML (no namespaces) OK
- Can handle top level namespaces OK
```
<?xml version="1.0"?>
<docs 
    xmlns:book="http://fake.org/book"
    xmlns:article="http://fake.org/article"
>
  <book:doc>
    A book entry
  </book:doc>
  <article:doc>
    A article entry
  </article:doc>
</docs>
```
- *Cannot* handle element- or attribute-level namespaces (*soon*), like:
```
<?xml version="1.0"?>
<docs>
  <doc>number one</doc>
  <p1:doc xmlns:p1="http://test.org/1">number two</p1:doc>
</docs>
```
- *Cannot* read `gz` or `bz2` compressed XML (soon)
- *Cannot* do [XML streaming](https://stackoverflow.com/questions/1134189/can-jaxb-parse-large-xml-files-in-chunks) (thus limited to smaller XML files) (perhaps soon?)


Copyright 2014,2015,2016 Glen Newton