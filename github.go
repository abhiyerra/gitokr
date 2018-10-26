package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func githubAuth(githubAccessToken string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return ctx, github.NewClient(tc)
}

func repoConfig(ctx context.Context, client *github.Client, owner string, repo string) (config Config, repoInfo *github.Repository) {
	repoConfig, _, _, err := client.Repositories.GetContents(ctx, owner, repo, configFile, nil)
	if err != nil {
		log.Fatal("Error", err)
	}

	repoInfo, _, err = client.Repositories.Get(ctx, owner, repo)
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
