package main

import "github.com/awalterschulze/gographviz"

type Member struct {
	Name string

	OKR OKRs
}

func (m *Member) WriteGraph(g *gographviz.Graph, srcNode string) {
	g.AddNode("G", m.Name, nil)
	g.AddNode("G", m.Name, tableNode(m.Name, "", m.OKR.Trs()))
	g.AddEdge(srcNode, m.Name, true, nil)
}
