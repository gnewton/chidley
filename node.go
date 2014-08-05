package main

import (
//"log"
)

type Node struct {
	name         string
	space        string
	spaceTag     string
	parent       *Node
	children     map[string]*Node
	childCount   map[string]int
	repeats      bool
	nodeTypeInfo *NodeTypeInfo
}

type NodeVisitor interface {
	Visit(n *Node) bool
	AlreadyVisited(n *Node) bool
	SetAlreadyVisited(n *Node)
}

func (n *Node) initialize(name string, space string, spaceTag string, parent *Node) {
	n.parent = parent
	n.name = name
	n.space = space
	n.spaceTag = spaceTag
	n.children = make(map[string]*Node)
	n.childCount = make(map[string]int)
	n.nodeTypeInfo = new(NodeTypeInfo)
	n.nodeTypeInfo.initialize()
}

func (n *Node) makeName() string {
	spaceTag := ""
	if n.spaceTag != "" {
		spaceTag = "_" + n.spaceTag
	}
	return capitalizeFirstLetter(cleanName(n.name)) + spaceTag
}

func (n *Node) makeType(prefix string, suffix string) string {
	return makeTypeGeneric(n.name, n.spaceTag, prefix, suffix)
}

func makeTypeGeneric(name string, space string, prefix string, suffix string) string {
	spaceTag := ""
	if space != "" {
		spaceTag = "__" + space
	}
	return prefix + cleanName(name) + spaceTag + suffix

}
