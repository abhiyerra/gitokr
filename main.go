package main

import (
	"math/rand"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	githubAccessToken string

	awsAccessKey       string
	awsSecretAccessKey string
)

const (
	configFile  = ".gitsop/config.json"
	dynamoTable = "GitSOP"
)

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	// TODO
	var repos = []string{
		"acksin/consulting",
		"acksin/gitlead",
		"acksin/SaleIron",
		"abhiyerra/dotfiles",
		"startupsonoma/community",
	}
	crons(repos)
}
