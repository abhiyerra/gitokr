package main

import (
	"log"

	"github.com/awalterschulze/gographviz"
)

type Links []*Link

type Link struct {
	ID string `yaml:"ID"`
}

func (m *Link) WriteGraph(g *gographviz.Graph, srcNode string) {
	log.Println(m)
	g.AddEdge(srcNode, m.ID, true, map[string]string{
		//		"dir":   "both",
		"style": `"dotted"`,
	})
}
