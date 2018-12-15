package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/google/go-github/github"
	flag "github.com/spf13/pflag"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
)

var (
	githubAccessToken string
	githubOwner       string
	githubRepo        string
)

type PlaybookJobs []PlaybookJob

type PlaybookJob struct {
	Name    string `yaml:"Name"`
	Content string `yaml:"Content"`
}

func (p PlaybookJob) CreateGithubIssue() {
	newIssue := &github.IssueRequest{
		Title: github.String(p.Name),
		Body:  github.String(p.Content),
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	_, _, err := githubClient.Issues.Create(ctx, githubOwner, githubRepo, newIssue)
	if err != nil {
		log.Println("Error", err)
	}
}

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Owner")
	flag.StringVar(&githubRepo, "github-repo", "", "Github Repo")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	var playbooks PlaybookJobs

	b, _ := ioutil.ReadFile(flag.Arg(0))

	err := yaml.Unmarshal(b, &playbooks)
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range playbooks {
		log.Println(i.Content)
		i.CreateGithubIssue()
	}
}
