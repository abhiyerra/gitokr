package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitsopConfig map[string]struct {
	Assignee    string `json:"assignee"`
	PullRequest string `json:"pullRequest"`
	FileName    string `json:"fileName"`
	OutputDir   string `json:"outputDir"`
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
	)

	flag.StringVar(&githubRepo, "github-repo", "", "Github Repo. Ex. gitsop")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Owner. Ex. abhiyerra")
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	ctx, client := githubAuth(githubAccessToken)

	// list all repositories for the authenticated user
	repo, _, err := client.Repositories.Get(ctx, githubOwner, githubRepo)
	if err != nil {
		log.Fatal("Error", err)
	}

	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/src-d/go-git",
		Progress: os.Stdout,
	})

	fileContent := []byte("This is the content of my file\nand the 2nd line of it")

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("This is my commit message"),
		Content:   fileContent,
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	_, _, err := client.Repositories.CreateFile(ctx, "myOrganization", "myRepository", "myNewFile.md", opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	newPR := &github.NewPullRequest{
		Title:               github.String("My awesome pull request"),
		Head:                github.String("branch_to_merge"),
		Base:                github.String("master"),
		Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(context.Background(), "myOrganization", "myRepository", newPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
}

func githubAuth(githubAccessToken string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return ctx, github.NewClient(tc)
}
