package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
	flag "github.com/spf13/pflag"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var (
	githubAccessToken string
	awsProfile        string
)

const (
	configFile  = "CRON.yml"
	dynamoTable = "GitSOP"
)

type Cron struct {
	Name     string `yaml:"Name"`
	Schedule string `yaml:"Schedule"`

	Files  []string               `yaml:"Files"`
	Inputs map[string]interface{} `yaml:"Inputs"`

	Type string `yaml:"Type"`

	Lambda struct {
		FunctionName string `yaml:"FunctionName"`
		Region       string `yaml:"Region"`
	} `yaml:"Lambda"`

	Webhook string `yaml:"Webhook"`
	Github  struct {
		Owner     string    `yaml:"Owner"`
		Repo      string    `yaml:"Repo"`
		Assignees *[]string `yaml:"Assignees"`
	} `yaml:"Github"`

	cronRun struct {
		Task    string
		NextRun time.Time
	}
}

type Crons []*Cron

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

func (t *Cron) RunCron() {
	tableKey := t.Name

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

		switch t.Type {
		case "Github":
			t.newGithubIssue(issueText)
		case "Webhook":
			t.newWebhook(issueText)
		case "Lambda":
			t.newLambda(issueText)
		default:
			t.newGithubIssue(issueText)
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

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

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

func (c *Cron) newLambda(issueText string) {
	sess := session.Must(session.NewSession())

	svc := lambda.New(sess, &aws.Config{
		Credentials: credentials.NewSharedCredentials("", awsProfile),
		Region:      aws.String(c.Lambda.Region),
	})

	payload := struct {
		IssueText string
	}{
		IssueText: issueText,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}

	result, err := svc.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(c.Lambda.FunctionName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        b,
	})

	if err != nil {
		log.Println(err, string(result.Payload))
		return
	}
}

var (
	dyno *dynamodb.DynamoDB
)

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.StringVar(&awsProfile, "aws-profile", "", "AWS Profile")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "opszero"),
	})

	if err != nil {
		log.Fatal(err)
	}
	// Create DynamoDB githubClient
	dyno = dynamodb.New(sess)

	log.Println(flag.Arg(0))
	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Println(err)
	}

	var crons Crons

	log.Println(string(b))

	err = yaml.Unmarshal(b, &crons)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(crons)

	for _, i := range crons {
		i.RunCron()
	}
}
