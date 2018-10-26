package main

import (
	"bytes"
	"fmt"
	"html/template"
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

type Input struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Task struct {
	Cron      string           `json:"cron"`
	Assignee  string           `json:"assignee"`
	Assignees *[]string        `json:"assignees"`
	Files     []string         `json:"files"`
	Webhook   string           `json:"webhook"`
	Inputs    map[string]Input `json:"inputs"`

	cronRun struct {
		Title   string
		NextRun time.Time
	}
}

func (t *Task) RunCron(rc *RepoConfig, title string) {
	if t.Cron == "" {
		return
	}

	tableKey := fmt.Sprintf("%s-%s", rc.repoInfo.GetGitURL(), title)

	schedule, err := cron.Parse(t.Cron)
	if err != nil {
		log.Fatal("Error", err)
	}

	result, err := rc.dyno.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamoTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Task": {
				S: aws.String(tableKey),
			},
		},
	})

	log.Println(err)

	err = dynamodbattribute.UnmarshalMap(result.Item, &t.cronRun)

	log.Println(err)

	log.Println("Foobar", t.cronRun, result.Item)

	if time.Now().After(t.cronRun.NextRun) {
		t.createSOP(rc, title)

		t.cronRun.Title = tableKey
		t.cronRun.NextRun = schedule.Next(time.Now())

		av, err := dynamodbattribute.MarshalMap(t.cronRun)
		if err != nil {
			log.Printf("failed to DynamoDB marshal Record, %v", err)
		}

		_, err = rc.dyno.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(dynamoTable),
			Item:      av,
		})
		if err != nil {
			log.Printf("failed to put Record to DynamoDB, %v", err)
		}
	}
}

func (task *Task) createSOP(rc *RepoConfig, title string) {
	timeNow := time.Now().UTC().Format("Mon Jan _2")
	var issueText []string

	log.Println(task)

	for _, fileName := range task.Files {
		repoConfig, _, _, err := rc.githubClient.Repositories.GetContents(rc.ctx, rc.repoInfo.GetOwner().GetLogin(), rc.repoInfo.GetName(), fileName, nil)
		log.Println(fileName)
		if err != nil {
			log.Println("Error", err)
			continue
		}

		fileContent, err := repoConfig.GetContent()
		if err != nil {
			log.Fatal("Error", err)
		}

		var fileContentBytes bytes.Buffer

		log.Println(fileContent)

		t := template.Must(template.New("t1").Funcs(template.FuncMap{
			"weekday": time.Now().Weekday().String,
		}).Parse(fileContent))
		if err != nil {
			log.Println(err)
		}
		t.Execute(&fileContentBytes, task.Inputs)

		issueText = append(issueText, string(fileContentBytes.String()))
	}

	if task.Webhook == "" {
		newIssue := &github.IssueRequest{
			Title:     github.String(fmt.Sprintf("%s: %s", timeNow, title)),
			Body:      github.String(strings.Join(issueText, "\n")),
			Assignees: task.Assignees,
		}

		pr, _, err := rc.githubClient.Issues.Create(rc.ctx, rc.repoInfo.GetOwner().GetLogin(), rc.repoInfo.GetName(), newIssue)
		if err != nil {
			log.Println("Error", err)
		}

		log.Printf("PR created: %s\n", pr.GetHTMLURL())
	} else {
		resp, err := http.PostForm(task.Webhook, url.Values{
			"title":   {title},
			"body":    {strings.Join(issueText, "\n")},
			"timeNow": {timeNow},
		})
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()
	}
}

// func runTask(taskName string, taskInputs map[string]string, owner string, repo string) {
// 	log.Println(taskInputs)
// 	ctx, client := githubAuth(githubAccessToken)
// 	config, repoInfo := repoConfig(ctx, client, owner, repo)

// 	task, ok := config[taskName]
// 	if !ok {
// 		log.Fatal(taskName, "doesn't exist.")
// 	}

// 	createTask(ctx, config, client, repoInfo, taskName, task, taskInputs)
// }
