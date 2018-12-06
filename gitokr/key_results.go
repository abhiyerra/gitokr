package main

import "time"

type KeyResult struct {
	Metric      string     `yaml:"Metric"`
	GithubIssue string     `yaml:"GithubIssue"`
	ReviewedAt  *time.Time `yaml:"ReviewedAt"`
	Done        bool       `yaml:"Done"`
}
