#!/usr/bin/env bash

pubmed_filename="pubmedsample18n0001.xml.gz"
pubmed_filename="pubmed18n0001.xml.gz"
pubmed_url="ftp://ftp.ncbi.nlm.nih.gov/pubmed/baseline/"
#pubmed_url="ftp://ftp.ncbi.nlm.nih.gov/pubmed/baseline-2018-sample/"
pubmed_url+=$pubmed_filename

test_filename="../data/test.xml"

#echo "" > foo

if [ ! -e "$pubmed_filename" ]
then
  wget $pubmed_url
fi

function out(){
    echo "$1"
}

function exe(){
    out "\`\`\`"
    out "$($1 2>&1)"
    out "\`\`\`"
}


############start
out "# \`chidley\` converts *any* XML to Go structs (and therefor to JSON)"
out "* By *any*, any XML that can be read by the Go [xml package](http://golang.org/pkg/encoding/xml/) decoder. "
out "* Where *convert* means, generates Go code that when compiled, will convert the XML to JSON"
out "* or will just generate the Go structs that represent the input XML"
out "* or converts XML to XML (useful for validation) "
out ""
out "Author: Glen Newton"
out "Language: Go"
out ""
out "## How does it work (with such a small memory footprint)"
out "\`chidley\` uses the input XML to build a model of each XML element (aka tag)."
out "It examines each instance of a tag, and builds a (single) prototypical representation; that is, the union of all the attributes and all of the child elements of all instances of the tag."
out "So even if there are million instances of a specific tag, there is only one model tag representation."
out "Note that a tag is unique by its namespace+tagname combination (in the Go xml package parlance, [\`space + local\`](http://golang.org/pkg/encoding/xml/#Name)."
out "### Types"
out "\`chidley\` by default makes all values (attributes, tag content) in the generated Go structs a string (which is always valid for XML attributes and content), but it has a flag (\`-t\`) where it will detect and use the most appropriate type. "
out "\`chidley\` tries to fit the **smallest** Go type. "
out "For example, if all instances of a tag contain a number, and all instances are -128 to 127, then it will use an \`int8\` in the Go struct."

out ""
out "## Usage"
exe "../chidley -h"

out ""
out "### Specific Usages:"
out "* \`chidley -W ...\`: writes Go code to standard out, so this output should be directed to a filename and subsequently be compiled. When compiled, the resulting binary will:"
out "    * convert the XML file to JSON"
out "    * or convert the XML file to XML (useful for validation)"
out "    * or count the # of elements (space, local) in the XML file"
out "* \`chidley -G ...\`: writes just the Go structs that represent the input XML. For incorporation into the user's code base."

out "### Example:"
out "Using filename \`test.xml\`"
exe "cat ../test.xml"
cmd="../chidley ./${test_filename}"

out $cmd
out "### Generated Go structs"
exe "$cmd"

if [ ! -e "test1" ]
then
  mkdir test1
fi



eval "../chidley -W ./${test_filename} > test1/t.go"

cd test1
go build


out "### Usage -W"
exe "./test1 -h"

out "##### Generated code:: xml -> json"
exe "./test1 -j -s"

out "##### Generated code: xml -> xml"
exe "./test1 -x -s"

out "##### Generated code: Count elements -c"
exe "./test1 -c"

cd ..
rm test1/*




out "### Example chidley -G:"
out "#### Default"

out "#### Types turned on -t"
exe "../chidley -t ./${test_filename}"

out "Note the \`Number int16\` in \`Chiyear\`"

out "## Larger  more complex example"


out "Using the large pubmed XML file, $pubmed_url "
cmd="../chidley -t -F ./${pubmed_filename}"
out "\`\$ ${cmd}"
exe "$cmd"


out "Timings"
out "### Generate Go program: -W"
out "#### Generated program: count tags"

out "#### Generated program: convert XML to JSON"
out "##### No streaming"
out "##### With streaming -s"
out "##### Sample of generated JSON"

out "##### Sample of generated XML to XML"

out "## <a name=java>Java/JAXB</a>"
out "### Usage"


out "### New"

out "#### Build Java package"
out "#### Running"


out "### Limitations"


