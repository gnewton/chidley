package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	//"strconv"
	"strings"
	"unicode"

	"github.com/xi2/xz"
)

func genericReader(filename string) (io.Reader, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasSuffix(filename, "xz") {
		reader, err := xz.NewReader(bufio.NewReader(file), 0)
		if err != nil {
			return nil, nil, err
		}
		return bufio.NewReader(reader), file, err
	}

	if strings.HasSuffix(filename, "bz2") {
		return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err
	}

	if strings.HasSuffix(filename, "gz") {
		reader, err := gzip.NewReader(bufio.NewReader(file))
		if err != nil {
			return nil, nil, err
		}
		return bufio.NewReader(reader), file, err
	}
	return bufio.NewReader(file), file, err
}

func cleanName(name string) string {
	for old, new := range nameMapper {
		parts := strings.Split(name, old)
		if len(parts) == 1 {
			continue
		} else {
			name := ""
			l := len(parts)
			for i := 0; i < l; i++ {
				if i == 0 {
					name += parts[i]
				} else {
					name += capitalizeFirstLetter(parts[i])
				}
				if i+1 < l {
					name += new
				}
			}
		}
	}
	//return capitalizeFirstLetter(name)
	return name
}

const BoolType = "bool"
const StringType = "string"
const IntType = "int"
const Int8Type = "int8"
const Int16Type = "int16"
const Int32Type = "int32"
const Int64Type = "int64"
const Float32Type = "float32"
const Float64Type = "float64"

func findType(nti *NodeTypeInfo, useType bool) string {
	if !useType {
		return StringType
	}

	if nti.alwaysBool {
		return BoolType
	}

	if nti.alwaysInt08 {
		return Int8Type
	}
	if nti.alwaysInt16 {
		return Int16Type
	}
	if nti.alwaysInt32 {
		return Int32Type
	}
	if nti.alwaysInt64 {
		return Int64Type
	}

	if nti.alwaysInt0 {
		return IntType
	}

	if nti.alwaysFloat32 {
		return Float32Type
	}
	if nti.alwaysFloat64 {
		return Float64Type
	}
	return StringType
}

type fqnSorter []*FQN

func isStringOnlyField(n *Node, nattributes int) bool {
	return (len(n.children) == 0 && nattributes == 0)
}

func makeAttributes(writer io.Writer, attributes []*FQN, nameSpaceTagMap map[string]string) {
	sort.Sort(fqnSorter(attributes))

	for _, fqn := range attributes {

		name := fqn.name
		nameSpace := fqn.space

		nameSpaceTag, ok := nameSpaceTagMap[nameSpace]
		if ok && nameSpaceTag != "" {
			nameSpaceTag = nameSpaceTag + "Space"
		} else {
			nameSpaceTag = nameSpace
		}

		nameSpaceTag = goVariableNameSanitize(nameSpaceTag)
		if len(nameSpace) > 0 {
			nameSpace = nameSpace + " "
		}

		//variableName := attributePrefix + capitalizeFirstLetter(nameSpaceTag) + cleanName(name)
		variableName := goVariableNameSanitize(attributePrefix + capitalizeFirstLetter(nameSpaceTag) + cleanName(name))
		variableType := "string"

		//lineChannel <- "\t" + variableName + " " + variableType + "`xml:\"" + nameSpace + name + ",attr\"  json:\",omitempty\"`" + "  // maxLength=" + strconv.Itoa(fqn.maxLength)
		//fmt.Fprintln(writer, "\t"+variableName+" "+variableType+"`xml:\""+nameSpace+name+",attr\"  json:\",omitempty\"`"+"  // maxLength="+strconv.Itoa(fqn.maxLength))
		fmt.Fprintln(writer, "\t"+variableName+" "+variableType+"`xml:\""+nameSpace+name+",attr\"  json:\",omitempty\"`")
	}
}

func goVariableNameSanitize(s string) string {
	s = strings.Replace(s, ":", "_colon_", -1)
	s = strings.Replace(s, "/", "_slash_", -1)
	s = strings.Replace(s, ".", "_dot_", -1)
	s = strings.Replace(s, "-", "_dash_", -1)
	s = strings.Replace(s, " ", "_space_", -1)
	s = strings.Replace(s, "-", "_dash_", -1)

	return s
}

// Len is part of sort.Interface.
func (s fqnSorter) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s fqnSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s fqnSorter) Less(i, j int) bool {
	return strings.Compare(cleanName(s[i].name), cleanName(s[j].name)) < 0
}

// node key
func nk(n *Node) string {
	return nks(n.space, n.name)
}

func nks(space, name string) string {
	return space + "NS" + name
}

func getFullPaths(filenames []string) ([]string, error) {
	fps := make([]string, len(filenames))
	var err error
	for i, _ := range filenames {
		fps[i], err = getFullPath(filenames[i])
		if err != nil {
			return nil, err
		}
	}

	return fps, nil
}

func getFullPath(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}
	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Print("Error opening: " + filename)
		return "", err
	}
	return file.Name(), nil
}

type alterer func(string) string

func alterFirstLetter(s string, f alterer) string {
	switch len(s) {
	case 0:
		return s
	case 1:
		v := f(s[0:1])
		return v

	default:
		return f(s[0:1]) + s[1:]
	}
}

func capitalizeFirstLetter(s string) string {
	return alterFirstLetter(s, strings.ToUpper)
}

func lowerFirstLetter(s string) string {
	return alterFirstLetter(s, strings.ToLower)
}

func findThisAttribute(local, nameSpace string, attrs []*FQN) *FQN {
	for _, attr := range attrs {
		if attr.name == local && attr.space == nameSpace {
			return attr
		}
	}
	return nil
}

func makeAttributeName(key, namespace, local string) string {
	return key + "_" + local + "_" + namespace

}

func containsUnicodeSpace(s string) bool {
	if s == "" {
		return false
	}

	for _, rune := range s {
		//log.Printf("*** %#U ****", rune)
		if unicode.IsSpace(rune) {
			return true
		}
	}
	return false
}

func isIgnoredTag(tag string) bool {
	var ignored bool
	if ignoredXmlTagsMap == nil {
		return false
	}
	_, ignored = (*ignoredXmlTagsMap)[tag]

	if ignored || (ignoreLowerCaseXmlTags && tag == strings.ToLower(tag)) {
		return true
	}

	return false

}

func extractExcludedTags(tagsString string) (*map[string]struct{}, error) {
	ignoredMap := make(map[string]struct{})
	if tagsString == "" {
		return &ignoredMap, nil
	}

	tags := strings.Split(tagsString, ",")

	for i, _ := range tags {
		tag := strings.TrimSpace(tags[i])
		if containsUnicodeSpace(tag) {
			return nil, errors.New("Excluded tag contains space: [" + tag + "] in list of excluded tags:" + tagsString + "]")
		}
		ignoredMap[tag] = struct{}{}

	}
	return &ignoredMap, nil
}

func findFieldNameFromTypeInfo(t string) string {
	switch t {
	case IntType, Int8Type, Int16Type, Int32Type, Int64Type, Float32Type, Float64Type:
		return cdataNumberName
	case BoolType:
		return cdataBooleanName
	}
	return cdataStringName
}
