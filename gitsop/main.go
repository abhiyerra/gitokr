package main

import (
	"math/rand"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	githubAccessToken string
)

const (
	configFile  = "SOP.yaml"
	dynamoTable = "GitSOP"
)

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	rcs = append(rcs, NewRepoConfig(i, svc))
	go RunHTTP()
}
