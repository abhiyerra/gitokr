package main

import "time"

type KeyResult struct {
	Metric      string
	GithubIssue string
	ReviewedAt  *time.Time
	Status      bool
}
