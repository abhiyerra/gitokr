package main

type Input struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Task struct {
	Cron      string           `json:"cron"`
	Assignee  string           `json:"assignee"`
	Assignees *[]string        `json:"assignees"`
	Files     []string         `json:"files"`
	Webhook   string           `json:"webhook"`
	Inputs    map[string]Input `json:"inputs"`
}
