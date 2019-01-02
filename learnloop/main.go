package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/google/go-github/github"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
)

func nodeName(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(reg.ReplaceAllString(input, ""), "_", "", -1)
}

func tableNode(l *Loop) (v map[string]string) {
	v = make(map[string]string)

	f := fmt.Sprintf(`<table border="0" cellspacing="0" cellborder="1">
    <tr>
     <td colspan="2" bgcolor="orange"><b>%s</b></td>
    </tr>
     <tr>
     <td colspan="2">%s</td>
     </tr></table>
    `, l.Learn, strings.Join(l.Do, "<br>"))

	v = map[string]string{
		"shape": "plaintext",
		"label": "<" + f + ">",
	}

	return
}

var (
	githubAccessToken string
)

func isYaml(fileName string) bool {
	return strings.HasSuffix(fileName, "yml") || strings.HasSuffix(fileName, "yaml")
}

func githubClient() (ctx context.Context, githubClient *github.Client) {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient = github.NewClient(tc)

	return ctx, githubClient
}

type Loop struct {
	Learn   string   `yaml:"Learn"`
	Do      []string `yaml:"Do"`
	Advance bool     `yaml:"Advance"`
	Iterate []Loop   `yaml:"Iterate"`
}

func (c *Loop) WriteGraph(g *gographviz.Graph, srcNode string) {
	currentNodeName := nodeName(uuid.NewV4().String())

	g.AddNode("G", currentNodeName, tableNode(c))
	if srcNode != "" {
		g.AddEdge(srcNode, currentNodeName, true, nil)
	}

	for _, loop := range c.Iterate {
		loop.WriteGraph(g, currentNodeName)
	}
}

func NewLoop(b []byte) *Loop {
	loop := &Loop{}

	err := json.Unmarshal(b, loop)
	if err != nil {
		log.Fatal(err)
	}

	return loop
}

func NewLoopFromYaml(b []byte) *Loop {
	loop := &Loop{}

	err := yaml.Unmarshal(b, loop)
	if err != nil {
		log.Fatal(err)
	}

	return loop
}

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	//log.SetFlags(log.Llongfile)

	fileName := flag.Arg(0)
	b, _ := ioutil.ReadFile(fileName)

	var loop *Loop
	if isYaml(fileName) {
		loop = NewLoopFromYaml(b)
	} else {
		loop = NewLoop(b)
	}

	g, _ := gographviz.Read([]byte(`digraph G {}`))
	if err := g.SetName("G"); err != nil {
		panic(err)
	}

	loop.WriteGraph(g, "")
	fmt.Printf(g.String())
}
