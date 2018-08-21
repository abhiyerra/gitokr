package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"text/template"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
	"golang.org/x/oauth2"
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

type GitSOPConfig []struct {
	Cron      string   `json:"cron"`
	Assignee  string   `json:"assignee"`
	Files     []string `json:"files"`
	OutputDir string   `json:"outputDir"`

	Title        string `json:"title"`
	Instructions string `json:"instructions"`

	Inputs map[string]Input `json:"inputs"`
}

func main() {
	var (
		githubOwner       string
		githubRepo        string
		githubAccessToken string

		config GitSOPConfig
	)

	flag.StringVar(&githubRepo, "github-repo", "", "Github Repo. Ex. gitsop")
	flag.StringVar(&githubOwner, "github-owner", "", "Github Owner. Ex. abhiyerra")
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

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

	for _, task := range config {
		var (
			timeNow    = time.Now().UTC().Format(time.RFC3339)
			branchName = namesgenerator.GetRandomName(1)
		)

		log.Println("Branch", branchName)

		if task.Cron == "" {
			continue
		}

		schedule, err := cron.Parse(task.Cron)
		if err != nil {
			log.Fatal("Error", err)
		}

		if schedule.Next(time.Now()).After(time.Now()) {

			log.Println(task)

			log.Println(repoInfo.GetMasterBranch())
			// Create Branch

			branch, _, err := client.Repositories.GetBranch(ctx, githubOwner, githubRepo, "master")
			if err != nil {
				log.Fatal("Error", err)
			}

			_, _, err = client.Git.CreateRef(ctx, githubOwner, githubRepo, &github.Reference{
				Ref: github.String(fmt.Sprintf("refs/heads/%s", branchName)),
				Object: &github.GitObject{
					SHA: branch.Commit.SHA,
				},
			})
			if err != nil {
				log.Fatal("Error", err)
			}

			for _, fileName := range task.Files {
				repoConfig, _, _, err := client.Repositories.GetContents(ctx, githubOwner, githubRepo, fileName, nil)
				if err != nil {
					log.Fatal("Error", err)
				}

				fileContent, err := repoConfig.GetContent()
				if err != nil {
					log.Fatal("Error", err)
				}
				log.Println(fileContent)

				var fileContentBytes bytes.Buffer

				t := template.Must(template.New("t1").Parse(fileContent))
				t.Execute(&fileContentBytes, task.Inputs)

				opts := &github.RepositoryContentFileOptions{
					Message: github.String(fmt.Sprintf("%s: %s", timeNow, task.Title)),
					Content: fileContentBytes.Bytes(),
					Branch:  github.String(branchName),
					Committer: &github.CommitAuthor{
						Name:  github.String("GitSOP"),
						Email: github.String("bot@gitsop.com"),
					},
				}
				_, _, err = client.Repositories.CreateFile(ctx, githubOwner, githubRepo, filepath.Join(task.OutputDir, timeNow, fileName), opts)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			newPR := &github.NewPullRequest{
				Title:               github.String(fmt.Sprintf("%s: %s", timeNow, task.Title)),
				Head:                github.String(branchName),
				Base:                github.String("master"),
				Body:                github.String(task.Instructions),
				MaintainerCanModify: github.Bool(true),
			}

			pr, _, err := client.PullRequests.Create(ctx, githubOwner, githubRepo, newPR)
			if err != nil {
				log.Fatal("Error", err)
			}

			log.Printf("PR created: %s\n", pr.GetHTMLURL())
		}
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
