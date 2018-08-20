package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitsopConfig map[string]struct {
	Assignee  string `json:"assignee"`
	FileName  string `json:"fileName"`
	OutputDir string `json:"outputDir"`
}

// gitsop

// - Look at all the files
// - Look for cronjobs in the file. https://godoc.org/github.com/robfig/cron#Schedule
// 	- If it is the next time to run then copy the file and create a run file on a branch.
// 	- Create a pull request?

// - gitsop filename/foobar
func main() {
	var (
		githubOwner       string
		githubRepo        string
		githubAccessToken string

		config GitsopConfig
	)

	flag.StringVar(&githubRepo, "github-repo", "", "Github Repo. Ex. gitsop")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Owner. Ex. abhiyerra")
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	ctx, client := githubAuth(githubAccessToken)

	repoConfig, _, _, err := client.Repositories.GetContents(ctx, githubOwner, githubRepo, "gitsop.json", nil)
	if err != nil {
		log.Fatal("Error", err)
	}

	repoInfo, _, err := client.Repositories.Get(ctx, githubOwner, githubRepo)
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

	for k, v := range config {
		log.Println(k)

		log.Println(repoInfo.GetMasterBranch())
		// Create Branch

		branch, _, err := client.Repositories.GetBranch(ctx, githubOwner, githubRepo, "master")
		if err != nil {
			log.Fatal("Error", err)
		}

		_, _, err = client.Git.CreateRef(ctx, githubOwner, githubRepo, &github.Reference{
			Ref: github.String("refs/heads/foobar"),
			Object: &github.GitObject{
				SHA: branch.Commit.SHA,
			},
		})
		if err != nil {
			log.Fatal("Error", err)
		}

		repoConfig, _, _, err := client.Repositories.GetContents(ctx, githubOwner, githubRepo, v.FileName, nil)
		if err != nil {
			log.Fatal("Error", err)
		}

		fileContent, err := repoConfig.GetContent()
		if err != nil {
			log.Fatal("Error", err)
		}
		log.Println(fileContent)

		opts := &github.RepositoryContentFileOptions{
			Message: github.String("This is my commit message"),
			Content: []byte(fileContent),
			Branch:  github.String("foobar"),
			Committer: &github.CommitAuthor{
				Name:  github.String("FirstName LastName"),
				Email: github.String("user@example.com")},
		}
		_, _, err = client.Repositories.CreateFile(ctx, githubOwner, githubRepo, "foobar.md", opts)
		if err != nil {
			fmt.Println(err)
			return
		}

		newPR := &github.NewPullRequest{
			Title:               github.String("My awesome pull request"),
			Head:                github.String("foobar"),
			Base:                github.String("master"),
			Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
			MaintainerCanModify: github.Bool(true),
		}

		pr, _, err := client.PullRequests.Create(ctx, githubOwner, githubRepo, newPR)
		if err != nil {
			log.Fatal("Error", err)
		}

		log.Printf("PR created: %s\n", pr.GetHTMLURL())
	}

}

func githubAuth(githubAccessToken string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return ctx, github.NewClient(tc)
}
