package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type Member struct {
	Name string
	OKR  OKRs

	// Cron
}

func (m *Member) WriteGraph(g *gographviz.Graph, srcNode string) {
	g.AddNode("G", nodeName(srcNode, m.Name), tableNode(fmt.Sprintf("Member: %s", m.Name), "", m.OKR.Trs()))
	g.AddEdge(srcNode, nodeName(srcNode, m.Name), true, nil)
}
