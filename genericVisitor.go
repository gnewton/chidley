package main

import (
	"log"
)

type GenericVisitor struct {
	alreadyVisited      map[string]bool
	alreadyVisitedNodes map[string]*Node
	globalTagAttributes map[string]([]*FQN)

	maxDepth        int
	depth           int
	nameSpaceTagMap map[string]string
	useType         bool
	nodeInfoList    []*NodeInfo
}

type NodeInfo struct {
	TypeName       string
	AttributeNames []string
	SubElements    []*SubElement
	IsCompressed   bool
	HasCharData    bool
}

type SubElement struct {
	Name         string
	TypeName     string
	Length       string
	IsList       bool
	IsPointer    bool
	IsCompressed bool
}

func (v *GenericVisitor) init(maxDepth int, globalTagAttributes map[string]([]*FQN), nameSpaceTagMap map[string]string, useType bool, nameSpaceInJsonName bool) {
	v.alreadyVisited = make(map[string]bool)
	v.alreadyVisitedNodes = make(map[string]*Node)
	v.globalTagAttributes = make(map[string]([]*FQN))
	v.globalTagAttributes = globalTagAttributes
	v.maxDepth = maxDepth
	v.depth = 0
	v.nameSpaceTagMap = nameSpaceTagMap
	v.useType = useType
	v.nodeInfoList = make([]*NodeInfo, 0)
}

func (v *GenericVisitor) Visit(node *Node) bool {
	v.depth += 1

	if v.AlreadyVisited(node) || node.ignoredTag {
		v.depth += 1
		return false
	}

	v.SetAlreadyVisited(node)
	ni := new(NodeInfo)
	v.nodeInfoList = append(v.nodeInfoList, ni)
	ni.TypeName = node.name
	ni.AttributeNames = fqnNames(v.globalTagAttributes[nk(node)])

	ni.HasCharData = node.hasCharData
	log.Println(node.name)
	log.Println(isStringOnlyField(node, len(v.globalTagAttributes[nk(node)])))
	if node.hasCharData {
		log.Println("***************************")
	}

	ni.HasCharData = isStringOnlyField(node, len(v.globalTagAttributes[nk(node)]))
	if !ni.HasCharData {
		for i, _ := range node.children {
			log.Println(i)
			log.Println(" > " + node.children[i].name)

		}
		ni.SubElements = make([]*SubElement, 0)
		for i, _ := range node.children {
			if !node.children[i].ignoredTag {
				child := node.children[i]
				if contains(collapsedXmlTagsList, child.name) {
					//log.Println("Collapsed")
					ni.makeCollapsed(child, v.globalTagAttributes)
				} else {
					v.Visit(child)
					sub := new(SubElement)
					ni.SubElements = append(ni.SubElements, sub)
					sub.Name = child.name
					if isStringOnlyField(child, len(v.globalTagAttributes[nk(child)])) {
						sub.TypeName = findType(child.nodeTypeInfo, true)
						//sub.TypeName = "string"
					} else {
						sub.TypeName = child.makeType("", "")
					}
					if child.repeats {
						sub.IsList = true
					}
				}
			}

		}
	}
	v.depth += 1
	return true
}

func (v *GenericVisitor) AlreadyVisited(n *Node) bool {
	_, ok := v.alreadyVisited[nk(n)]
	return ok
}

func (v *GenericVisitor) SetAlreadyVisited(n *Node) {
	v.alreadyVisited[nk(n)] = true
	v.alreadyVisitedNodes[nk(n)] = n
}

func (ni *NodeInfo) makeCollapsed(node *Node, globalTagAttributes map[string]([]*FQN)) {
	attrs := fqnNames(globalTagAttributes[nk(node)])

	sub := new(SubElement)
	ni.SubElements = append(ni.SubElements, sub)
	sub.Name = node.name
	sub.TypeName = "string"

	for i, _ := range attrs {
		sub := new(SubElement)
		ni.SubElements = append(ni.SubElements, sub)
		sub.Name = node.name + "ATT" + attrs[i]
		sub.TypeName = "att string"
	}
}
