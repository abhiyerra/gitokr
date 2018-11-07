package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type Member struct {
	Name string

	OKR OKRs
}

func (m *Member) WriteGraph(g *gographviz.Graph, srcNode string) {
	n := fmt.Sprintf("%s%s", srcNode, m.Name)
	g.AddNode("G", n, nil)
	g.AddNode("G", n, tableNode(fmt.Sprintf("Member: %s", m.Name), "", m.OKR.Trs()))
	g.AddEdge(srcNode, n, true, nil)
}
