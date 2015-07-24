package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
)

func genericReader(filename string) (io.Reader, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
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
		name = strings.Replace(name, old, new, -1)
	}
	return name
}

func findType(nti *NodeTypeInfo, useType bool) string {
	if !useType {
		return "string"
	}

	if nti.alwaysBool {
		return "bool"
	}

	if nti.alwaysInt08 {
		return "int8"
	}
	if nti.alwaysInt16 {
		return "int16"
	}
	if nti.alwaysInt32 {
		return "int32"
	}
	if nti.alwaysInt64 {
		return "int64"
	}

	if nti.alwaysInt0 {
		return "int"
	}

	if nti.alwaysFloat32 {
		return "float32"
	}
	if nti.alwaysFloat64 {
		return "float64"
	}
	return "string"
}

func makeAttributes(lineChannel chan string, attributes []*FQN, nameSpaceTagMap map[string]string) {

	for _, fqn := range attributes {
		name := fqn.name
		space := fqn.space

		spaceTag, ok := nameSpaceTagMap[space]
		if ok && spaceTag != "" {
			spaceTag = spaceTag + "_"
		}

		lineChannel <- "\t" + attributePrefix + "_" + spaceTag + cleanName(name) + " string `xml:\"" + space + " " + name + ",attr\"  json:\",omitempty\"`"
	}
}

// node key
func nk(n *Node) string {
	return nks(n.space, n.name)
}

func nks(space, name string) string {
	return space + "___" + name
}

func getFullPath(filename string) string {
	if filename == "" {
		return ""
	}
	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Print("Error opening: " + filename)
		log.Fatal(err)
	}
	return file.Name()
}
