package logger

import (
	"fmt"

	"github.com/nlopes/slack"
)

func New(m *slack.Client, ch chan string, slackChannel string) *Logger {
	return &Logger{
		m:            m,
		ch:           ch,
		slackChannel: slackChannel,
	}
}

type Logger struct {
	m            *slack.Client
	ch           chan string
	slackChannel string
}

func (l *Logger) Start() {
	go l.handleLog()
}

func (l *Logger) Log(message string) {
	l.ch <- message
}

func (l *Logger) handleLog() {
	for message := range l.ch {
		fmt.Println(message)

		_, _, err := l.m.PostMessage(l.slackChannel, slack.MsgOptionText(message, false), slack.MsgOptionUsername("ceb.feeder"))
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	}
}
