package main

type Tasks []*Task

type Task struct {
	Name   string `yaml:"Name"`
	Github struct {
		Owner string `yaml:"Owner"`
		Repo  string `yaml:"Repo"`
		Issue string `yaml:"Issue"`
	} `yaml:"Github"`
}
