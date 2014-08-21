# `chidley`
## `chidley` converts *any* XML to JSON.
* By *any*, any XML that can be read by the Go [xml package](http://golang.org/pkg/encoding/xml/) decoder. 
* Where *convert* means, generates Go code that when compiled, will convert the XML to JSON
* Also converts XML to XML (useful for validation) 

## Usage
```
$ chidley -h
Usage of ./chidley:
  -G=false: Only write generated Go structs to stdout
  -W=false: Generate Go code to convert XML to JSON or XML (latter useful for validation) and write it to stdout
  -a="Attr": Prefix to attribute names
  -c=false: Read XML from standard input
  -d=false: Debug; prints out much information
  -e="Chi": Prefix to element names
  -p=false: Pretty-print json in generated code (if applicable)
  -r=false: Progress: every 50000 elements
  -s="Type": Suffix to element names
  -t=false: Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete
  -u=false: Filename interpreted as an URL
```
`chidley` writes Go code to standard out, so this output should be directed to a filename and subsequently be compiled.

###Example:
```
$ chidley -W 
```

##Usage of compiled code
If the output of `chidley` is directed to `chidCodeGen/C.go`, and `cd chidCodeGen; go build` is run, the compiled Go binary `chidCodeGen` is created.




##Name
`chidley` is names after [Cape Chidley](https://en.wikipedia.org/wiki/Cape_Chidley)

