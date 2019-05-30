package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nlopes/slack"
)

type Proxy struct {
	Url  string `json:"url"`
	Auth string `json:"auth"`
}

type Config struct {
	Port         int
	SlackToken   string
	SlackChannel string
	TestFlag     bool
	Limits       map[string]int
	List         []Proxy
}

var cfg Config

var slackClient *slack.Client

func init() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&cfg)
	cfg.Limits = map[string]int{
		"yandex":  20000,
		"sputnik": 20000,
		"google":  2000,
	}

	if err != nil {
		log.Fatal(err)
	}

	slackClient = slack.New(cfg.SlackToken)
}
