package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

const (
	dynamoTable = "GitSOP"
)

type Cron struct {
	Name     string `yaml:"Name"`
	Schedule string `yaml:"Schedule"`

	Files   []string          `yaml:"Files"`
	Webhook string            `yaml:"Webhook"`
	Inputs  map[string]string `yaml:"Inputs"`

	Github struct {
		Owner     string    `yaml:"Owner"`
		Repo      string    `yaml:"Repo"`
		Assignees *[]string `yaml:"Assignees"`
	} `yaml:"Github"`

	cronRun struct {
		Task    string
		NextRun time.Time
	}
}

func (c *Cron) joinFiles() string {
	var issueText []string

	for _, fileName := range c.Files {
		var (
			buf bytes.Buffer
		)

		b, _ := ioutil.ReadFile(fileName)

		t := template.Must(template.New("t1").Funcs(template.FuncMap{
			"weekday": time.Now().Weekday().String,
		}).Parse(string(b)))

		t.Execute(&buf, c.Inputs)

		issueText = append(issueText, buf.String())
	}

	return strings.Join(issueText, "\n")
}

func (t *Cron) RunCron(srcName string) {
	tableKey := nodeName(srcName, t.Name)

	schedule, err := cron.Parse(t.Schedule)
	if err != nil {
		log.Fatal("Error", t.Name, err)
	}

	result, err := dyno.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamoTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Task": {
				S: aws.String(tableKey),
			},
		},
	})

	err = dynamodbattribute.UnmarshalMap(result.Item, &t.cronRun)

	log.Println(t.cronRun.NextRun)
	if time.Now().After(t.cronRun.NextRun) {
		var (
			issueText = t.joinFiles()
		)

		if t.Webhook == "" {
			t.newGithubIssue(issueText)
		} else {
			t.newWebhook(issueText)
		}

		t.cronRun.Task = tableKey
		t.cronRun.NextRun = schedule.Next(time.Now())

		av, err := dynamodbattribute.MarshalMap(t.cronRun)
		if err != nil {
			log.Printf("failed to DynamoDB marshal Record, %v", err)
		}

		_, err = dyno.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(dynamoTable),
			Item:      av,
		})
		if err != nil {
			log.Printf("failed to put Record to DynamoDB, %v", err)
		}
	}
}

func (c *Cron) newGithubIssue(issueText string) {
	timeNow := time.Now().UTC().Format("Mon Jan _2")
	newIssue := &github.IssueRequest{
		Title:     github.String(fmt.Sprintf("%s: %s", timeNow, c.Name)),
		Body:      github.String(issueText),
		Assignees: c.Github.Assignees,
	}

	ctx, githubClient := githubClient()

	_, _, err := githubClient.Issues.Create(ctx, c.Github.Owner, c.Github.Repo, newIssue)
	if err != nil {
		log.Println("Error", err)
	}
}

func (c *Cron) newWebhook(issueText string) {
	timeNow := time.Now().UTC().Format("Mon Jan _2")
	resp, err := http.PostForm(c.Webhook, url.Values{
		"title":   {c.Name},
		"body":    {issueText},
		"timeNow": {timeNow},
	})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}

type Crons []*Cron

func (o Crons) Table() (text string) {
	text2 := `<table border="0" cellspacing="0" cellborder="1">`
	for _, t := range o {
		text2 += fmt.Sprintf(`<tr><td align="left">%s</td></tr>`, t.Name)
	}
	text2 += "</table>"
	text += fmt.Sprintf(`<tr><td>Cron:</td><td>%s</td></tr>`, text2)

	return text
}
