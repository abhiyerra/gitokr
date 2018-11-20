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

	ExternalProjects []*ExternalProject `yaml:"ExternalProjects"`
	Projects         []*Project         `yaml:"Projects"`
	Members          []*Member          `yaml:"Members"`
	Crons            Crons              `yaml:"Crons"`
	Tasks            Tasks              `yaml:"Tasks"`
	SOPs             SOPs               `yaml:"SOPs"`
}

func (c *Project) NodeName() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Name)
}

func (c *Project) WriteGraph(g *gographviz.Graph, srcNode string) {
	if c.Type == "" {
		c.Type = DefaultProjectType
	}
	g.AddNode("G", nodeName(srcNode, c.Name), tableNode(c.NodeName(), c.Vision, c.OKR.Trs(), nil))
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

	for _, sop := range c.SOPs {
		sop.WriteGraph(g, nodeName(srcNode, c.Name))
	}
}

func (c *Project) RunCrons(srcNode string) {
	for _, cron := range c.Crons {
		cron.RunCron(nodeName(srcNode, c.Name))
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

	return project
}
