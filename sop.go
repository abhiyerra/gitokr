package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/awalterschulze/gographviz"
)

type SOPs []*SOP

type SOP struct {
	Name   string `yaml:"Name"`
	File   string `yaml:"File"`
	Github struct {
		Owner string `yaml:"Owner"`
		Repo  string `yaml:"Repo"`
	} `yaml:"Github"`

	fileContent string
}

func (m *SOP) githubLink() string {
	v := url.Values{}
	v.Set("title", m.Name)
	v.Set("body", m.fileContent)

	return fmt.Sprintf("https://github.com/%s/%s/issues/new?", m.Github.Owner, m.Github.Repo) + v.Encode()
}

func (m *SOP) WriteGraph(g *gographviz.Graph, srcNode string) {
	b, err := ioutil.ReadFile(m.File)
	if err != nil {
		log.Println(err)
	}

	m.fileContent = html.EscapeString(string(b))

	g.AddNode("G", nodeName(srcNode, m.Name), map[string]string{
		"label": `"SOP: ` + m.Name + `"`,
		"color": `"blue"`,
		"URL":   `"` + m.githubLink() + `"`,
	})
	g.AddEdge(srcNode, nodeName(srcNode, m.Name), true, map[string]string{
		"style": `"dotted"`,
	})
}
