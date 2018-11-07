package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type Group struct {
	Name string
	OKR  OKRs

	Members []*Member
	Groups  []*Group
}

func (c *Group) WriteGraph(g *gographviz.Graph, srcNode string) {
	g.AddNode("G", c.Name, tableNode(fmt.Sprintf("Group: %s", c.Name), "", c.OKR.Trs()))
	g.AddEdge(srcNode, c.Name, true, nil)

	for _, member := range c.Members {
		member.WriteGraph(g, c.Name)
	}
	for _, group := range c.Groups {
		group.WriteGraph(g, c.Name)
	}
}
