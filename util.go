package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
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

func makeAttributes(attributes []*FQN, nameSpaceTagMap map[string]string) []string {
	all := make([]string, 0)
	for _, fqn := range attributes {
		name := fqn.name
		space := fqn.space

		spaceTag, ok := nameSpaceTagMap[space]
		if ok && spaceTag != "" {
			spaceTag = "__" + spaceTag
		}

		attStr := "\t" + attributePrefix + cleanName(name) + spaceTag + " string `xml:\"" + space + " " + name + ",attr\"  json:\",omitempty\"`"
		all = append(all, attStr)
	}
	return all
}

// node key
func nk(n *Node) string {
	return nks(n.space, n.name)
}

func nks(space, name string) string {
	return space + "___" + name
}
