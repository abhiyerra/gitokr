package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/robfig/cron"
)

type CronRun struct {
	Task    string
	NextRun time.Time
}

func crons(repos []string) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", "opszero"),
	})
	if err != nil {
		log.Fatal("Error", err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	for {

		for _, repo := range repos {

			spl := strings.Split(repo, "/")

			go func(owner, repo string) {
				ctx, client := githubAuth(githubAccessToken)
				config, repoInfo := repoConfig(ctx, client, owner, repo)

				for title, task := range config {
					if task.Cron == "" {
						continue
					}

					tableKey := fmt.Sprintf("%s-%s", repoInfo.GetGitURL(), title)

					schedule, err := cron.Parse(task.Cron)
					if err != nil {
						log.Fatal("Error", err)
					}

					result, err := svc.GetItem(&dynamodb.GetItemInput{
						TableName: aws.String(dynamoTable),
						Key: map[string]*dynamodb.AttributeValue{
							"Task": {
								S: aws.String(tableKey),
							},
						},
					})

					log.Println(err)

					nextRun := CronRun{}
					err = dynamodbattribute.UnmarshalMap(result.Item, &nextRun)

					log.Println(err)

					log.Println("Foobar", nextRun, result.Item)

					if time.Now().After(nextRun.NextRun) {
						createTask(ctx, config, client, repoInfo, title, task, nil)

						nextRun.Task = tableKey
						nextRun.NextRun = schedule.Next(time.Now())

						av, err := dynamodbattribute.MarshalMap(nextRun)
						if err != nil {
							log.Printf("failed to DynamoDB marshal Record, %v", err)
						}

						_, err = svc.PutItem(&dynamodb.PutItemInput{
							TableName: aws.String(dynamoTable),
							Item:      av,
						})
						if err != nil {
							log.Printf("failed to put Record to DynamoDB, %v", err)
						}
					}
				}
			}(spl[0], spl[1])

		}
		time.Sleep(time.Minute)
	}
}
