package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func nodeName(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(input, "")
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

func main() {
	project := &Project{}

	b, _ := ioutil.ReadFile(os.Args[1])

	err := json.Unmarshal(b, project)
	if err != nil {
		log.Fatal(err)
	}

	project.WriteGraph()
}
