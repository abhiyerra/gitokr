package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

const (
	DefaultGroupType = "Group"
	SystemGroupType  = "System"
	APIGroupType     = "API"
)

type Group struct {
	Name        string
	Description string

	OKR  OKRs
	Type string

	// Tasks

	Members []*Member
	Groups  []*Group
}

func (c *Group) NodeName() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Name)
}

func (c *Group) WriteGraph(g *gographviz.Graph, srcNode string) {
	if c.Type == "" {
		c.Type = DefaultGroupType
	}
	g.AddNode("G", nodeName(c.Name), tableNode(c.NodeName(), c.Description, c.OKR.Trs()))
	g.AddEdge(srcNode, nodeName(c.Name), true, nil)

	for _, member := range c.Members {
		member.WriteGraph(g, c.Name)
	}
	for _, group := range c.Groups {
		group.WriteGraph(g, c.Name)
	}
}
