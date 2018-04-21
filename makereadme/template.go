package main

const endInline = `

` + "`" + `` + "`" + `` + "`"

const beginInline = ` 

` + endInline

const readmeTemplate = `
# ` + "`" + `chidley` + "`" + `

# NOTE: The below documentation is out of date with the most recent release. I will work to update these docs in the next week. 2018/10/17

###############################################################################################################################################

## ` + "`" + `chidley` + "`" + ` converts *any* XML to Go structs (and therefor to JSON)
* By *any*, any XML that can be read by the Go [xml package](http://golang.org/pkg/encoding/xml/) decoder. 
* Where *convert* means, generates Go code that when compiled, will convert the XML to JSON
* or will just generate the Go structs that represent the input XML
* or converts XML to XML (useful for validation) 

Author: Glen Newton
Language: Go

## New

# 2016.08.14
Added ability to sort structs into the same order the XML is encountered in the file. Useful for human readers comparing the Go structs to the original XML.
Use flag ` + "`" + `-X` + "`" + ` to invoke. Overrides default of sorting by alphabetical sorting.

# 2015.07.24
` + "`" + `chidley` + "`" + ` now supports the user naming of the resulting JAXB Java class package.
Previously the package name could only be ` + "`" + `ca/gnewton/chidley/jaxb` + "`" + `.
Now, using the ` + "`" + `-P name` + "`" + `, the ` + "`" + `jaxb` + "`" + ` default can be altered.

So ` + "`" + `chidley -J -P "foobar" sample.xml` + "`" + ` will result in Java classes with package name: ` + "`" + `ca/gnewton/chidley/foobar` + "`" + `.


### Previous

` + "`" + `chidley` + "`" + ` now has support for Java/JAXB. It generates appropriate Java/JAXB classes and associated maven pom.

See [Java/JAXB section below](#user-content-java) for usage.


## How does it work (with such a small memory footprint)
` + "`" + `chidley` + "`" + ` uses the input XML to build a model of each XML element (aka tag).
It examines each instance of a tag, and builds a (single) prototypical representation; that is, the union of all the attributes and all of the child elements of all instances of the tag.
So even if there are million instances of a specific tag, there is only one model tag representation.
Note that a tag is unique by its namespace+tagname combination (in the Go xml package parlance, [` + "`" + `space + local` + "`" + `](http://golang.org/pkg/encoding/xml/#Name).
### Types
` + "`" + `chidley` + "`" + ` by default makes all values (attributes, tag content) in the generated Go structs a string (which is always valid for XML attributes and content), but it has a flag (` + "`" + `-t` + "`" + `) where it will detect and use the most appropriate type. 
` + "`" + `chidley` + "`" + ` tries to fit the **smallest** Go type. 
For example, if all instances of a tag contain a number, and all instances are -128 to 127, then it will use an ` + "`" + `int8` + "`" + ` in the Go struct.

## ` + "`" + `chidley` + "`" + ` binary
Compiled for 64bit Linux Fedora18, go version go1.3 linux/amd64

## Usage` + beginInline + `
$ chidley -h
{{.ChidleyUsage}}
$` + endInline + `

### Specific Usages:
* ` + "`" + `chidley -W ` + "`" + `: writes Go code to standard out, so this output should be directed to a filename and subsequently be compiled. When compiled, the resulting binary will:
    * convert the XML file to JSON
    * or convert the XML file to XML (useful for validation)
    * or count the # of elements (space, local) in the XML file
* ` + "`" + `chidley -G ...` + "`" + `: writes just the Go structs that represent the input XML. For incorporation into the user's code base.


### Example
####  ` + "`" + `data/test.xml` + "`" + `:` + beginInline + `
{{.SimpleExampleXMLFile}}` + endInline + `

#### Generated Go structs:
` + beginInline + `
{{.SimpleExampleXMLChidleyGoStructs}}` + endInline + `

Note that ` + "`" + `chidley` + "`" + `prepends a ` + "`" + `C` + "`" + `in front of every Go struct that corresponds to an XML tag in the original XML. It also does not alter the tag part of the Go struct name (except where noted below in "Name Changes"). 
The reason this is done is 1) The Go XML and JSON libraries only operate on public fields (name must start with a capital); and 2) To avoid name collisions.
Here is a simple and contrived example of name collisions in XMK:
` + beginInline + `
<record>
    <name>Fred</name>
    <Name>Gone with the Wind</Name>
</record>
` + endInline + `
In order to work` + "`" + `<name>` + "`" + ` would have to be capitalized, causing a collision with ` + "`" + `<Name>` + "`" + `
Prefixing both with a capital c avoid this.

This prefix can be changed with the ` + "`" + `-e` + "`" + ` flag



Note that all XMl tags are converted to Go structs. 
However, for those that always correspond to a single element (no sub-tags & no XML attributes), like title:` + "`" + `<title>Footfall</title>` + "`" + ` this is a bit of a waste.
It is possible to have ` + "`" + `chidley` + "`" + `collapse these into inline strings with the ` + "`" + `-F` + "`" + `flag. However, if your example XML is not canonical (i.e. it does not exhibit all uses of all XML tags), it may result in Go structs that do not capture everything that is needed.

## Name changes: XML vs. Go structs
XML names can contain dots ` + "`" + `.` + "`" + ` and hyphens or dashes ` + "`" + `-` + "`" + `. These are not valid identifiers for Go structs or variables. These are mapped as:
* ` + "`" + `"-": "_"` + "`" + `
* ` + "`" + `".": "_dot_"` + "`" + `


Example:
` + beginInline + `
` + "`$" + `chidley -F data/test.xml` + "`" + `:
{{.SimpleExampleXMLChidleyGoStructsCollapsed}}` + endInline + `

#### Generating code to read XML` + `
` + "`chidley`" + ` can generate Go code that will read in the XML and output a number of things, including the equivalent JSON, XML (XML to XML useful for validation), and a count total for each XML tag in the source file

Generating code:
` + beginInline + `
$ mkdir gencode
$ chidley -W data/test.xml> gencode/main.go
` + endInline + `

#### Usage of generated code
` + beginInline + `
$ cd gencode
$ go build
$ ./gencode
{{.GeneratedUsage}}
` + endInline + `

##### Generated code: Convert XML to JSON ` + "`" + `-j` + "`" + beginInline + `
$ ./test1 -j -f ../../xml/test1.xml 
{{.GeneratedXMLToJson}}` + endInline + `

##### Generated code: Convert XML to XML ` + "`" + `-x` + "`" + `
` + beginInline + `$ ./test1 -x -f ../../xml/test1.xml 
{{.GeneratedXMLToXML}}` + endInline + `

##### Generated code: Count elements ` + "`" + `-c` + "`" + `
XML elements (or tags) are counted in the source file (space,local) and are printed-out, unsorted
` + beginInline + `$ gencode -c
{{.GeneratedCountElements}}` + endInline + `

**Note**: the underscore before the colon indicates there is no (or the default) namespace for the element. 

## Type example
` + "`" + `` + "`" + `` + "`" + `
{{.SimpleExampleXMLFile}}

` + "`" + `` + "`" + `` + "`" + `

#### Default
` + "`" + `` + "`" + `` + "`" + `
$ ./chidley -G xml/testType.xml
{{.SimpleExampleXMLChidleyGoStructs}}

` + "`" + `` + "`" + `` + "`" + `

#### Types turned on ` + "`" + `-t` + "`" + `
` + "`" + `` + "`" + `` + "`" + `
$ ./chidley -G -t data/xml.test
{{.SimpleExampleXMLChidleyGoStructsWithTypes}}

` + "`" + `` + "`" + `` + "`" + `

Notice:
* ` + "`" + `Text int8` + "`" + ` in ` + "`" + `Chi_age` + "`" + `
* ` + "`" + `Text bool` + "`" + ` in ` + "`" + `Chi_married` + "`" + `

## Go struct name prefix
` + "`" + `chidley` + "`" + ` by default prepends a prefix to Go struct type identifiers. The default is ` + "`" + `Chi` + "`" + ` but this can be changed with the ` + "`" + `-e` + "`" + ` flag. If changed from the default, the new prefix must start with a capital letter (for the XML annotation and decoder to work: the struct fields must be public).

## Warning
If you are going to use the ` + "`" + `chidley` + "`" + ` generated Go structs on XML other than the input XML, you need to make sure the input XML has examples of all tags, and tag attribute and tag child tag combinations. 

If the input does not have all of these, and you use new XML that has tags not found in the input XML, attributes not seen in tags in the input XML, or child tags not encountered in the input XML, these will not be *seen* by the xml decoder, as they will not be in the Go structs used by the xml decoder.

## Limitations
` + "`" + `chidley` + "`" + ` is constrained by the underlying Go [xml package](http://golang.org/pkg/encoding/xml/)
Some of these limitations include:
* The default encoding supported by ` + "`" + `encoder/xml` + "`" + ` is UTF-8. Right now ` + "`" + `chidley` + "`" + ` does not support additional charsets. 
An xml decoder that handles charsets other than UTF-8 is possible (see example https://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go). 
It is possible that this method might be used in the future to extend ` + "`" + `chidley` + "`" + ` to include a small set of popular charsets.
* For vanilla XML with no namespaces, there should be no problem using ` + "`" + `chidley` + "`" + `

### Go ` + "`" + `xml` + "`" + ` package Namespace issues
* There are a number of bugs open for the Go xml package that relate to XML namespaces: https://code.google.com/p/go/issues/list?can=2&q=xml+namespace  If the XML you are using uses namespaces in certain ways, these bugs will impact whether ` + "`" + `chidley` + "`" + ` can create correct structs for your XML
* For _most_ XML with namespaces, the JSON will be OK but if you convert XML to XML using the generated Go code, there will be a chance one of the above mentioned bugs may impact results. Here is an example I encountered: https://groups.google.com/d/msg/golang-nuts/drWStJSt0Pg/Z47JHeij7ToJ

## Name
` + "`" + `chidley` + "`" + ` is named after [Cape Chidley](https://en.wikipedia.org/wiki/Cape_Chidley), Canada

## Larger & more complex example
Using the file ` + "`" + `{{.PubmedXMLFileName}}` + "`" + `. Generated from a query to pubmed (similar but much larger than [http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=20598978,444444,455555&retmode=xml](http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=20598978,444444,455555&retmode=xml)), returning a document in the [MEDLINE/Pubmed XML format](http://www.nlm.nih.gov/bsd/licensee/data_elements_doc.html).
* Compressed size: 27M
* Uncompressed size: 337M

### Generate Go structs from {{.PubmedXMLFileName}}: ` + "`" + `-G` + "`" + `

` + "`" + `` + "`" + `` + "`" + `

{{.PubmedExampleXMLChidleyGoStructsWithTypes}}

` + "`" + `` + "`" + `` + "`" + `
####Run times
$ /usr/bin/time -f "%E %M" ./chidley -G xml/pubmed_xml_12750255.xml.bz2 
{{.PubmedExampleXMLChidleyGoStructsWithTypeTiming}}

*Note:* All timings from Dell laptop 16GB, regular disk, 8 core i7-3720QM CPU @ 2.60GHz)
` + "`" + `Linux 3.11.10-100.fc18.x86_64 #1 SMP Mon Dec 2 20:28:38 UTC 2013 x86_64 x86_64 x86_64 GNU/Linux` + "`" + `



### Generate Go program: ` + "`" + `-W` + "`" + `

` + "`" + `` + "`" + `` + "`" + `
$ ./chidley -W xml/pubmed_xml_12750255.xml.bz2 > examples/pubmed/ChiPubmed.go
$ cd examples/pubmed/
$ go build
$
` + "`" + `` + "`" + `` + "`" + `

#### Generated program: count tags
` + "`" + `` + "`" + `` + "`" + `
$ /usr/bin/time -f "%E %M" ./pubmed -c | sort -n
0:36.58 10460
{{.GeneratedPubmedCount}}

12121212

$
` + "`" + `` + "`" + `` + "`" + `

*Note:* The underscore before the colon indicates there is no (or the default) namespace for the element.

36 seconds for 337MB XML; resident size: 10.5MB

#### Generated program: convert XML to JSON

##### No streaming
` + "`" + `` + "`" + `` + "`" + `
$ /usr/bin/time -f "%E %M" ./pubmed -j > /dev/null
{{.GeneratedPubmedNoStreaming}}
0:57.26 2866408
$
` + "`" + `` + "`" + `` + "`" + `
57 seconds for 337MB XML; resident size: 2.9GB

##### With streaming ` + "`" + `-s` + "`" + `
Streaming decodes using the XML elements that are one level down from the top level container element.
` + "`" + `` + "`" + `` + "`" + `
{{.GeneratedPubmedStreaming}}
$ /usr/bin/time -f "%E %M" ./pubmed -j -s > /dev/null
0:58.72 15944
` + "`" + `` + "`" + `` + "`" + `
59 seconds for 337MB XML; resident size: 16MB

##### Sample of generated JSON
` + "`" + `` + "`" + `` + "`" + `
$ /usr/bin/time -f "%E %M" ./pubmed -j -s |head -310
{{.GeneratedPubmedXMLToJson}}
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
` + "`" + `` + "`" + `` + "`" + `
 
##### Sample of generated XML to XML
` + "`" + `` + "`" + `` + "`" + `
$ ./pubmed -x -s |head -100
{{.GeneratedPubmedXMLToXML}}
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
` + "`" + `` + "`" + `` + "`" + `

## <a name="java">Java/JAXB</a>
` + "`" + `chidley` + "`" + ` now supports the production of Java/JAXB code. It generates a class-per-element, with classes mapping to the Go structs generated by the XML extraction.
Its only dependency is Google [Gson](https://code.google.com/p/google-gson/), for JSON generation.


### Usage
` + "`" + `chidley` + "`" + ` creates a maven project in ` + "`" + `./java` + "`" + ` (settable using the ` + "`" + `-D` + "`" + ` flag) and creates Java JAXB files in ` + "`" + `src/main/java/ca/gnewton/chidley/jaxb/xml` + "`" + `. 
It creates a ` + "`" + `Main.java` + "`" + ` in ` + "`" + `src/main/java/ca/gnewton/chidley/jaxb` + "`" + `
` + "`" + `` + "`" + `` + "`" + `
$ chidley -J xml/test1.xml
{{.ChidleyGenerateJava}}
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiDocs.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiDoc.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiTitle.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiAuthor.java
2014/09/02 10:22:27 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiLast_name.java
2014/09/02 10:22:28 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/xml/ChiFirstName.java
2014/09/02 10:22:28 printJavaJaxbVisitor.go:100: Writing java Class file: java/src/main/java/ca/gnewton/chidley/jaxb/Main.java
` + "`" + `` + "`" + `` + "`" + `

### New
Changing the package with "-P":
` + "`" + `` + "`" + `` + "`" + `
$ chidley -J -P testFoo xml/test1.xml
{{.ChidleyGenerateJavaChangePackageName}}
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiDocs.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiDoc.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiTitle.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiAuthor.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiLast_name.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/xml/ChiFirstName.java
2015/07/24 15:56:52 printJavaJaxbVisitor.go:103: Writing java Class file: java/src/main/java/ca/gnewton/chidley/testFoo/Main.java
$
` + "`" + `` + "`" + `` + "`" + `

#### Build Java package
` + "`" + `` + "`" + `` + "`" + `
$ cd java
$ ls
pom.xml  src
$ $ mvn package
{{.ChidleyGenerateJavaMavenBuild}}
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
` + "`" + `` + "`" + `` + "`" + ` 

#### Running
` + "`" + `` + "`" + `` + "`" + `
$ export CLASSPATH=./target/chidley-jaxb-1.0-SNAPSHOT.jar:/home/myhome/.m2/repository/com/google/code/gson/gson/2.3/gson-2.3.jar:$CLASSPATH
$ java ca.gnewton.chidley.jaxb.Main
zzzzzz {{.ChidleyGenerateJavaRun}} mmmmm
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
` + "`" + `` + "`" + `` + "`" + `

### Limitations
- Can handle vanilla XML (no namespaces) OK
- Can handle top level namespaces OK
` + "`" + `` + "`" + `` + "`" + `
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
n</docs>
` + "`" + `` + "`" + `` + "`" + `
- *Cannot* handle element- or attribute-level namespaces (*soon*), like:
` + "`" + `` + "`" + `` + "`" + `
<?xml version="1.0"?>
<docs>
  <doc>number one</doc>
  <p1:doc xmlns:p1="http://test.org/1">number two</p1:doc>
</docs>
` + "`" + `` + "`" + `` + "`" + `
- *Cannot* read ` + "`" + `gz` + "`" + ` or ` + "`" + `bz2` + "`" + ` compressed XML (soon)
- *Cannot* do [XML streaming](https://stackoverflow.com/questions/1134189/can-jaxb-parse-large-xml-files-in-chunks) (thus limited to smaller XML files) (perhaps soon?)


Copyright 2014,2015,2016 Glen Newton
`
