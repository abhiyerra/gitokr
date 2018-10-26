package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/github"
)

type Config map[string]Task

func createTask(ctx context.Context, config Config, client *github.Client, repoInfo *github.Repository, title string, task Task, taskInputs map[string]string) {
	timeNow := time.Now().UTC().Format("Mon Jan _2")
	var issueText []string

	log.Println(task)

	for _, fileName := range task.Files {
		repoConfig, _, _, err := client.Repositories.GetContents(ctx, repoInfo.GetOwner().GetLogin(), repoInfo.GetName(), fileName, nil)
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

		for inputName, input := range taskInputs {
			t, ok := task.Inputs[inputName]
			if ok {
				t.Value = input
				task.Inputs[inputName] = t
			}
		}

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

		pr, _, err := client.Issues.Create(ctx, repoInfo.GetOwner().GetLogin(), repoInfo.GetName(), newIssue)
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

func runTask(taskName string, taskInputs map[string]string, owner string, repo string) {
	log.Println(taskInputs)
	ctx, client := githubAuth(githubAccessToken)
	config, repoInfo := repoConfig(ctx, client, owner, repo)

	task, ok := config[taskName]
	if !ok {
		log.Fatal(taskName, "doesn't exist.")
	}

	createTask(ctx, config, client, repoInfo, taskName, task, taskInputs)
}
