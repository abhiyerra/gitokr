package main

import (
	"fmt"
	"log"
	"regexp"
)

type OKRs map[string]*OKR

func (o OKRs) Trs() (trs []string) {

	for period, okr := range o {
		f := fmt.Sprintf(`<tr><td>%s</td></tr>`, period)
		f += fmt.Sprintf(`%s`, okr.Table())

		trs = append(trs, f)
	}

	return trs
}

type OKR struct {
	Objective  string
	KeyResults []KeyResult
	Status     string
}

func (o OKR) NodeName() string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(o.Objective, "")
}

func (o OKR) Table() string {
	var text string

	text += fmt.Sprintf(`<tr><td><b>Objective:</b> %s</td></tr>`, o.Objective)
	text += fmt.Sprintf(`<tr><td><b>Key Results:</b></td></tr>`)
	for t := range o.KeyResults {
		text += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", o.KeyResults[t].Metric, o.KeyResults[t].Status)
	}

	return text
}
