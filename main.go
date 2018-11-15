package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func nodeName(srcNode, input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return srcNode + strings.Replace(reg.ReplaceAllString(input, ""), "_", "", -1)
}

func tableNode(title, text string, tr []string) map[string]string {
	f := fmt.Sprintf(`<table border="0" cellspacing="0" cellborder="1">
    <tr>
     <td colspan="2" bgcolor="orange"><b>%s</b></td>
    </tr>
     <tr>
     <td colspan="2">%s</td>
     </tr>%s</table>
    `, title, text, strings.Join(tr, ""))

	return map[string]string{
		"shape": "plaintext",
		"label": "<" + f + ">",
	}
}

var (
	githubAccessToken string
)

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	b, _ := ioutil.ReadFile(flag.Arg(0))

	project := NewProject(b)

	g, _ := gographviz.Read([]byte(`digraph G {}`))
	if err := g.SetName("G"); err != nil {
		panic(err)
	}

	project.WriteGraph(g, "")

	fmt.Printf(g.String())
}
