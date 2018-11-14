package main

import (
	"fmt"
	"log"
	"regexp"
)

type OKRs map[string]*OKR

func (o OKRs) Trs() (trs []string) {

	for period, okr := range o {
		f := fmt.Sprintf(`<tr><td colspan="2"><font>%s</font></td></tr>`, period)
		f += fmt.Sprintf(`%s`, okr.Table())

		trs = append(trs, f)
	}

	return trs
}

type OKR struct {
	Objective  string
	KeyResults []KeyResult
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

	if o.KeyResults == nil {
		return ""
	}

	text += fmt.Sprintf(`<tr><td>Objective:</td><td align="left">%s</td></tr>`, o.Objective)

	var text2 string = `<table border="0" cellspacing="0" cellborder="1">`
	for t := range o.KeyResults {
		status := ""
		if o.KeyResults[t].Done {
			status = "green"
		} else {
			status = "red"
		}
		text2 += fmt.Sprintf(`<tr><td align="left" color="%s">%s</td></tr>`, status, o.KeyResults[t].Metric)
	}
	text2 += "</table>"
	text += fmt.Sprintf(`<tr><td>Key Results:</td><td>%s</td></tr>`, text2)

	return text
}
