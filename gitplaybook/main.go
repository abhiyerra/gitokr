package main

import (
	"math/rand"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	githubAccessToken string
)

type PlaybookJobs []*PlaybookJobs

type PlaybookJob struct {
	Name    string `yaml:"Name"`
	Content string `yaml:"Content"`
}

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	// Create a Graph
	// Execute Each Event
}
