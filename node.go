package main

import (
//"log"
)

type Node struct {
	name         string
	nameSpace        string
	nameSpaceTag     string
	parent       *Node
	parents      []*Node
	children     map[string]*Node
	childCount   map[string]int
	repeats      bool
	nodeTypeInfo *NodeTypeInfo
	hasCharData  bool
	tempCharData string
}

type NodeVisitor interface {
	Visit(n *Node) bool
	AlreadyVisited(n *Node) bool
	SetAlreadyVisited(n *Node)
}

func (n *Node) initialize(name string, nameSpace string, nameSpaceTag string, parent *Node) {
	n.parent = parent
	n.parents = make([]*Node, 0, 0)
	n.pushParent(parent)
	n.name = name
	n.nameSpace = nameSpace
	n.nameSpaceTag = nameSpaceTag
	n.children = make(map[string]*Node)
	n.childCount = make(map[string]int)
	n.nodeTypeInfo = new(NodeTypeInfo)
	n.nodeTypeInfo.initialize()
	n.hasCharData = false
}

func (n *Node) makeName() string {
	nameSpaceTag := ""
	if n.nameSpaceTag != "" {
		nameSpaceTag = "_" + n.nameSpaceTag
	}
	return capitalizeFirstLetter(cleanName(n.name)) + nameSpaceTag
}

func (n *Node) makeType(prefix string, suffix string) string {
	return capitalizeFirstLetter(makeTypeGeneric(n.name, n.nameSpaceTag, prefix, suffix, false))
}

func (n *Node) makeJavaType(prefix string, suffix string) string {
	return capitalizeFirstLetter(makeTypeGeneric(n.name, n.nameSpaceTag, prefix, suffix, true))
}

func (n *Node) peekParent() *Node {
	if len(n.parents) == 0 {
		return nil
	}
	a := n.parents
	return a[len(a)-1]
}

func (n *Node) pushParent(parent *Node) {
	n.parents = append(n.parents, parent)
}

func (n *Node) popParent() *Node {
	if len(n.parents) == 0 {
		return nil
	}
	var poppedNode *Node
	a := n.parents
	poppedNode, n.parents = a[len(a)-1], a[:len(a)-1]
	return poppedNode
}

func makeTypeGeneric(name string, nameSpace string, prefix string, suffix string, capitalizeName bool) string {
	nameSpaceTag := ""
	if nameSpace != "" {
		nameSpaceTag = nameSpace + "_"
	}
	if capitalizeName {
		name = capitalizeFirstLetter(name)
	}
	return prefix + nameSpaceTag + cleanName(name) + suffix

}
