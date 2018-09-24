package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/github"
	"github.com/robfig/cron"
	flag "github.com/spf13/pflag"
	"golang.org/x/oauth2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	HumanInputType   = "human"
	CommandInputType = "command"
)

type Input struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Command string `json:"command"`
}

type Task struct {
	Cron      string    `json:"cron"`
	Assignee  string    `json:"assignee"`
	Assignees *[]string `json:"assignees"`
	Files     []string  `json:"files"`
	OutputDir string    `json:"outputDir"`

	Inputs map[string]Input `json:"inputs"`
}

type GitSOPConfig map[string]Task

var (
	githubOwner       string
	githubRepo        string
	githubAccessToken string

	awsAccessKey       string
	awsSecretAccessKey string
)

const (
	configFile  = ".gitsop/config.json"
	dynamoTable = "GitSOP"
)

func githubAuth(githubAccessToken string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return ctx, github.NewClient(tc)
}

func repoConfig(ctx context.Context, client *github.Client) (config GitSOPConfig, repoInfo *github.Repository) {
	repoConfig, _, _, err := client.Repositories.GetContents(ctx, githubOwner, githubRepo, configFile, nil)
	if err != nil {
		log.Fatal("Error", err)
	}

	repoInfo, _, err = client.Repositories.Get(ctx, githubOwner, githubRepo)
	if err != nil {
		log.Fatal("Error", err)
	}

	repoConfigContent, err := repoConfig.GetContent()
	if err != nil {
		log.Fatal("Error", err)
	}
	log.Println(repoConfigContent)

	err = json.Unmarshal([]byte(repoConfigContent), &config)
	if err != nil {
		log.Fatal("Error", err)
	}
	log.Println(config)

	return config, repoInfo
}

func createTask(ctx context.Context, config GitSOPConfig, client *github.Client, repoInfo *github.Repository, title string, task Task, taskInputs map[string]string) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	var issueText []string

	log.Println(task)

	for _, fileName := range task.Files {
		repoConfig, _, _, err := client.Repositories.GetContents(ctx, githubOwner, githubRepo, fileName, nil)
		if err != nil {
			log.Fatal("Error", err)
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

		t := template.Must(template.New("t1").Parse(fileContent))
		t.Execute(&fileContentBytes, task.Inputs)

		issueText = append(issueText, string(fileContentBytes.String()))
	}

	newIssue := &github.IssueRequest{
		Title:     github.String(fmt.Sprintf("%s: %s", timeNow, title)),
		Body:      github.String(strings.Join(issueText, "\n")),
		Assignees: task.Assignees,
	}

	pr, _, err := client.Issues.Create(ctx, githubOwner, githubRepo, newIssue)
	if err != nil {
		log.Fatal("Error", err)
	}

	log.Printf("PR created: %s\n", pr.GetHTMLURL())
}

type CronRun struct {
	Task    string
	NextRun time.Time
}

func crons() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "opszero"),
	})
	if err != nil {
		log.Fatal("Error", err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	for {
		ctx, client := githubAuth(githubAccessToken)
		config, repoInfo := repoConfig(ctx, client)

		for title, task := range config {
			if task.Cron == "" {
				continue
			}

			tableKey := fmt.Sprintf("%s-%s", repoInfo.GetGitURL(), title)

			schedule, err := cron.Parse(task.Cron)
			if err != nil {
				log.Fatal("Error", err)
			}

			result, err := svc.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(dynamoTable),
				Key: map[string]*dynamodb.AttributeValue{
					"Task": {
						S: aws.String(tableKey),
					},
				},
			})

			log.Println(err)

			nextRun := CronRun{}
			err = dynamodbattribute.UnmarshalMap(result.Item, &nextRun)

			log.Println(err)

			log.Println("Foobar", nextRun, result.Item)

			if time.Now().After(nextRun.NextRun) {
				createTask(ctx, config, client, repoInfo, title, task, nil)

				nextRun.Task = tableKey
				nextRun.NextRun = schedule.Next(time.Now())

				av, err := dynamodbattribute.MarshalMap(nextRun)
				if err != nil {
					log.Printf("failed to DynamoDB marshal Record, %v", err)
				}

				_, err = svc.PutItem(&dynamodb.PutItemInput{
					TableName: aws.String(dynamoTable),
					Item:      av,
				})
				if err != nil {
					log.Printf("failed to put Record to DynamoDB, %v", err)
				}
			}
		}

		time.Sleep(time.Minute)
	}
}

func runTask(taskName string, taskInputs map[string]string) {
	log.Println(taskInputs)
	ctx, client := githubAuth(githubAccessToken)
	config, repoInfo := repoConfig(ctx, client)

	task, ok := config[taskName]
	if !ok {
		log.Fatal(taskName, "doesn't exist.")
	}

	createTask(ctx, config, client, repoInfo, taskName, task, taskInputs)
}

func main() {
	var (
		taskName   string
		taskInputs map[string]string
	)

	flag.StringVar(&githubRepo, "github-repo", "", "Github Repo. Ex. gitsop")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Owner. Ex. abhiyerra")
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.StringVar(&taskName, "task", "", "The task to run.")
	flag.StringToStringVar(&taskInputs, "task-inputs", map[string]string{}, "Task Inputs")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	if taskName == "" {
		crons()
	} else {
		runTask(taskName, taskInputs)
	}
}
