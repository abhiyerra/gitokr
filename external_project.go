package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ExternalProject struct {
	Owner string
	Repo  string
	Path  string
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

	return NewProject([]byte(repoConfigContent))
}
