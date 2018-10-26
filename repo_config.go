package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type RepoConfig struct {
	repo   string
	config map[string]*Task

	ctx          context.Context
	githubClient *github.Client
	repoInfo     *github.Repository

	dyno *dynamodb.DynamoDB
}

func NewRepoConfig(repo string, svc *dynamodb.DynamoDB) (r *RepoConfig) {
	r = &RepoConfig{}

	r.ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(r.ctx, ts)
	r.githubClient = github.NewClient(tc)

	spl := strings.Split(repo, "/")

	repoConfig, _, _, err := r.githubClient.Repositories.GetContents(r.ctx, spl[0], spl[1], configFile, nil)
	if err != nil {
		log.Fatal("Error", err)
	}

	repoInfo, _, err := r.githubClient.Repositories.Get(r.ctx, spl[0], spl[1])
	if err != nil {
		log.Fatal("Error", err)
	}

	repoConfigContent, err := repoConfig.GetContent()
	if err != nil {
		log.Fatal("Error", err)
	}
	log.Println(repoConfigContent)

	err = json.Unmarshal([]byte(repoConfigContent), &r.config)
	if err != nil {
		log.Fatal("Error", err)
	}
	log.Println(r.config)

	r.repoInfo = repoInfo
	r.dyno = svc

	return r
}

func (c *RepoConfig) RunCrons() {
	for {
		for title, task := range c.config {
			go task.RunCron(c, title)
		}
		time.Sleep(time.Minute)
	}
}
