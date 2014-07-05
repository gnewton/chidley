package main

type Node struct {
	name         string
	space        string
	parent       *Node
	children     map[string]*Node
	childCount   map[string]int
	repeats      bool
	nodeTypeInfo *NodeTypeInfo
}

type Attributes map[string]bool

func (n *Node) initialize(name string, space string, parent *Node) {
	n.parent = parent
	n.name = name
	n.space = space
	n.children = make(map[string]*Node)
	n.childCount = make(map[string]int)
	n.nodeTypeInfo = new(NodeTypeInfo)
	n.nodeTypeInfo.initialize()
}
