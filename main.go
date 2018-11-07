package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func tableNode(title, text string, tr []string) map[string]string {
	f := fmt.Sprintf(`<table border="0" cellspacing="0" cellborder="1">
    <tr>

     <td>%s</td>
    </tr>
     <tr>
     <td>%s</td>
     </tr>%s</table>
    `, title, text, strings.Join(tr, ""))

	return map[string]string{"label": "<" + f + ">"}
}

func main() {
	company := &Company{}

	b, _ := ioutil.ReadFile(os.Args[1])

	err := json.Unmarshal(b, company)
	if err != nil {
		log.Fatal(err)
	}

	company.WriteGraph()
}
