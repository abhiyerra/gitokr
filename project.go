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
	ID     string `yaml:"ID"`
	Vision string `yaml:"Vision"`

	Type string `yaml:"Type"`

	OKR OKRs `yaml:"OKR"`

	ExternalProjects []*ExternalProject `yaml:"ExternalProjects"`
	Projects         []*Project         `yaml:"Projects"`
	Members          []*Member          `yaml:"Members"`
	Crons            Crons              `yaml:"Crons"`
	Tasks            Tasks              `yaml:"Tasks"`
	SOPs             SOPs               `yaml:"SOPs"`
	Links            Links              `yaml:"Links"`
}

func (c *Project) NodeName() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Name)
}

func (c *Project) WriteGraph(g *gographviz.Graph, srcNode string) {
	if c.Type == "" {
		c.Type = DefaultProjectType
	}

	currentNodeName := nodeName(srcNode, c.Name)
	if c.ID != "" {
		currentNodeName = c.ID
	}

	g.AddNode("G", currentNodeName, tableNode(c.NodeName(), c.Vision, c.OKR.Trs(), nil))
	if srcNode != "" {
		g.AddEdge(srcNode, currentNodeName, true, nil)
	}

	for _, e := range c.ExternalProjects {
		if proj := e.GetProject(); proj != nil {
			c.Projects = append(c.Projects, proj)
		}
	}

	for _, project := range c.Projects {
		project.WriteGraph(g, currentNodeName)
	}

	for _, member := range c.Members {
		member.WriteGraph(g, currentNodeName)
	}

	for _, sop := range c.SOPs {
		sop.WriteGraph(g, currentNodeName)
	}

	for _, link := range c.Links {
		link.WriteGraph(g, currentNodeName)
	}
}

type Score struct {
	KeyResultsFinished float32
	KeyResultsTotal    float32
}

func (c *Project) WriteScore() {
	scores := c.Score(make(map[string]Score))

	for k := range scores {
		log.Printf("Score: %s - %0.2f%%", k, scores[k].KeyResultsFinished/scores[k].KeyResultsTotal*100)
	}
}

func (c *Project) Score(scores map[string]Score) map[string]Score {
	for okr := range c.OKR {
		score := scores[okr]
		for _, i := range c.OKR[okr].KeyResults {
			if i.Done {
				score.KeyResultsFinished += 1.0
			}
			score.KeyResultsTotal += 1.0
		}
		scores[okr] = score
	}

	for _, i := range c.Projects {
		i.Score(scores)
	}

	return scores
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
