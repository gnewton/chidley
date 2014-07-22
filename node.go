package main

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

type Attributes map[string]bool

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
	spaceTag := ""
	if n.spaceTag != "" {
		spaceTag = "__" + n.spaceTag + "_"
	}
	//return prefix + capitalizeFirstLetter(n.name) + spaceTag + nameSuffix
	return prefix + cleanName(n.name) + spaceTag + nameSuffix
}
