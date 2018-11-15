package main

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ExternalProject struct {
	Owner string `yaml:"Owner"`
	Repo  string `yaml:"Repo"`
	Path  string `yaml:"Path"`
}

func (c *ExternalProject) GetProject() *Project {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	repoConfig, _, _, err := githubClient.Repositories.GetContents(ctx, c.Owner, c.Repo, c.Path, nil)
	if err != nil {
		return nil
	}

	repoConfigContent, err := repoConfig.GetContent()
	if err != nil {
		return nil
	}

	log.Println(c.Path)
	if isYaml(c.Path) {
		log.Println("yaml")
		return NewProjectFromYaml([]byte(repoConfigContent))
	}

	return NewProject([]byte(repoConfigContent))
}
