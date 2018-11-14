package main

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
)

type Project struct {
	Vision string
	OKR    OKRs
	Groups []*Group

	graph *gographviz.Graph
}

func (c *Project) WriteGraph() {
	c.graph, _ = gographviz.Read([]byte(`digraph G {}`))
	if err := c.graph.SetName("G"); err != nil {
		panic(err)
	}

	c.graph.AddNode("G", "Vision", tableNode("Vision", c.Vision, c.OKR.Trs()))

	for _, group := range c.Groups {
		group.WriteGraph(c.graph, "Vision")
	}
	fmt.Printf(c.graph.String())
}
