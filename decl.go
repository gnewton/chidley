package main

import (
	"strings"
)

var DEBUG = false
var attributePrefix = "Attr"
var codeGenConvert = false
var classicStructNamesWithUnderscores = false
var nameSpaceInJsonName = false
var prettyPrint = false
var progress = false
var readFromStandardIn = false
var sortByXmlOrder = false
var structsToStdout = true
var validateFieldTemplate = false

var ignoreLowerCaseXmlTags = false
var ignoredXmlTags = ""
var ignoredXmlTagsMap *map[string]struct{}

var ignoreXmlDecodingErrors = false

var codeGenDir = "codegen"
var codeGenFilename = "CodeGenStructs.go"

// Java out
const javaBasePackage = "ca.gnewton.chidley"
const mavenJavaBase = "src/main/java"

var javaBasePackagePath = strings.Replace(javaBasePackage, ".", "/", -1)
var javaAppName = "jaxb"
var writeJava = false
var baseJavaDir = "java"
var userJavaPackageName = ""

var namePrefix = "C"
var nameSuffix = ""
var xmlName = false
var url = false
var useType = false
var addDbMetadata = false
var flattenStrings = false

//FIXXX: should not be global
var keepXmlFirstLetterCase = true

var lengthTagName = ""
var lengthTagPadding int64 = 0
var lengthTagAttribute = ""
var lengthTagSeparator = ":"

var cdataStringName = "Text"
var cdataNumberName = "Number"
var cdataBooleanName = "Flag"

type structSortFunc func(v *PrintGoStructVisitor)

var structSort = printStructsAlphabetical

var outputs = []*bool{
	&codeGenConvert,
	&structsToStdout,
	&writeJava,
}
