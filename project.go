package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/awalterschulze/gographviz"
	"gopkg.in/yaml.v2"
)

const (
	DefaultProjectType = "Project"
	SystemProjectType  = "System"
	APIProjectType     = "API"
)

type Project struct {
	Name   string `yaml:"Name"`
	Vision string `yaml:"Vision"`

	Type string `yaml:"Type"`

	OKR OKRs `yaml:"OKR"`

	ExternalProjects []ExternalProject `yaml:"ExternalProjects"`
	Projects         []*Project        `yaml:"Projects"`
	Members          []*Member         `yaml:"Members"`
}

func (c *Project) NodeName() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Name)
}

func (c *Project) WriteGraph(g *gographviz.Graph, srcNode string) {
	if c.Type == "" {
		c.Type = DefaultProjectType
	}
	g.AddNode("G", nodeName(srcNode, c.Name), tableNode(c.NodeName(), c.Vision, c.OKR.Trs()))
	if srcNode != "" {
		g.AddEdge(srcNode, nodeName(srcNode, c.Name), true, nil)
	}

	for _, e := range c.ExternalProjects {
		if proj := e.GetProject(); proj != nil {
			c.Projects = append(c.Projects, proj)
		}
	}

	for _, project := range c.Projects {
		project.WriteGraph(g, nodeName(srcNode, c.Name))
	}

	for _, member := range c.Members {
		member.WriteGraph(g, nodeName(srcNode, c.Name))
	}
}

func NewProject(b []byte) *Project {
	project := &Project{}

	err := json.Unmarshal(b, project)
	if err != nil {
		log.Fatal(err)
	}

	return project
}

func NewProjectFromYaml(b []byte) *Project {
	project := &Project{}

	err := yaml.Unmarshal(b, project)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(project)

	return project
}
