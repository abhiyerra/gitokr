package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func nodeName(srcNode, input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return srcNode + strings.Replace(reg.ReplaceAllString(input, ""), "_", "", -1)
}

func tableNode(title, text string, tr []string, labels map[string]string) (v map[string]string) {
	v = make(map[string]string)

	f := fmt.Sprintf(`<table border="0" cellspacing="0" cellborder="1">
    <tr>
     <td colspan="2" bgcolor="orange"><b>%s</b></td>
    </tr>
     <tr>
     <td colspan="2">%s</td>
     </tr>%s</table>
    `, title, text, strings.Join(tr, ""))

	v = map[string]string{
		"shape": "plaintext",
		"label": "<" + f + ">",
	}

	for k, v2 := range labels {
		v[k] = v2
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

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	//log.SetFlags(log.Llongfile)

	fileName := flag.Arg(0)
	b, _ := ioutil.ReadFile(fileName)

	var project *Project
	if isYaml(fileName) {
		project = NewProjectFromYaml(b)
	} else {
		project = NewProject(b)
	}

	g, _ := gographviz.Read([]byte(`digraph G {}`))
	if err := g.SetName("G"); err != nil {
		panic(err)
	}

	project.WriteGraph(g, "")
	project.WriteScore()
	fmt.Printf(g.String())
}
