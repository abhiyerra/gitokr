package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	flag "github.com/spf13/pflag"
)

var (
	githubAccessToken string

	awsAccessKey       string
	awsSecretAccessKey string
)

const (
	configFile  = "SOP.yaml"
	dynamoTable = "GitSOP"
)

func main() {
	flag.StringVar(&githubAccessToken, "github-access-token", "", "Github Access Token")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "opszero"),
	})

	if err != nil {
		log.Fatal(err)
	}
	// Create DynamoDB githubClient
	svc := dynamodb.New(sess)

	// TODO
	var repos = []string{
		"acksin/consulting",
		"acksin/gitlead",
		"acksin/SaleIron",
		"abhiyerra/dotfiles",
		"startupsonoma/community",
	}

	var rcs []*RepoConfig
	waiter := make(chan struct{}, 1)

	for _, i := range repos {
		rcs = append(rcs, NewRepoConfig(i, svc))
	}

	go RunHTTP()

	for _, i := range rcs {
		go i.RunCrons()
	}

	<-waiter
}
