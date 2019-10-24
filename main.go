package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	ArgoCDBaseURL     string `required:"true" split_words:"true"`
	ArgoCDProject     string `required:"true" split_words:"true"`
	ArgoCDApplication string `required:"true" split_words:"true"`
	ArgoCDEventType   string `required:"true" split_words:"true"`
	SlackURL          string `required:"true" split_words:"true"`
	SlackChannel      string `required:"true" split_words:"true"`
	SlackUsername     string `default:"Notifier"`
	SlackIconEmoji    string `default:":monkey:" split_words:"true"`
	SlackText         string `default:"Deployment notification" split_words:"true"`
	SlackFooterText   string `default:"" split_words:"true"`
}

func (c config) appLink() string {
	return fmt.Sprintf("%s/applications/%s", c.ArgoCDBaseURL, c.ArgoCDApplication)
}

func (c config) projectLink() string {
	return fmt.Sprintf("%s/settings/projects/%s", c.ArgoCDBaseURL, c.ArgoCDProject)
}

func (c config) rollbackLink() string {
	return fmt.Sprintf("%s/applications/%s?rollback=0", c.ArgoCDBaseURL, c.ArgoCDApplication)
}

func main() {

	var err error
	var c config

	err = envconfig.Process("notifier", &c)
	if err != nil {
		log.Println("Config Error:", err.Error())
		os.Exit(1)
	}

	payload := slack.Payload{
		Channel:     c.SlackChannel,
		Username:    c.SlackUsername,
		IconEmoji:   c.SlackIconEmoji,
		Attachments: []slack.Attachment{},
	}
	payload.Attachments = attachments(c)

	if c.ArgoCDEventType == "PreSync" {
		payload.Text = c.SlackText
	}

	errs := slack.Send(c.SlackURL, "", payload)
	if len(errs) > 0 {
		log.Println("Send Errors:")
		for _, e := range errs {
			log.Println(" - ", e.Error())
		}
	}
}

func attachments(c config) []slack.Attachment {

	t := time.Now()
	ts := t.Unix()

	o := []slack.Attachment{}

	a := slack.Attachment{
		Footer:    &c.SlackFooterText,
		Timestamp: &ts,
	}

	status := ""
	colour := ""
	switch c.ArgoCDEventType {
	case "PreSync":
		status = "Starting 💥"
		colour = "warning"
		break
	case "Sync":
		status = "Synchronising 🤖"
		colour = "warning"
		break
	case "PostSync":
		status = "Deployment Completed 🚀"
		colour = "good"
		break
	case "SyncFail":
		status = "Failed 💀"
		colour = "danger"
		break
	}

	a.Color = &colour

	a.AddField(slack.Field{
		Title: "",
		Value: status,
		Short: false,
	})

	a.AddField(slack.Field{
		Title: "Project",
		Value: fmt.Sprintf("<%s|%s>", c.projectLink(), c.ArgoCDProject),
		Short: true,
	})

	a.AddField(slack.Field{
		Title: "Application",
		Value: fmt.Sprintf("<%s|%s>", c.appLink(), c.ArgoCDApplication),
		Short: true,
	})

	if c.ArgoCDEventType == "PostSync" {
		a.AddAction(slack.Action{
			Type:  "button",
			Text:  "Initiate Rollback",
			Url:   c.rollbackLink(),
			Style: "danger",
		})
	}

	o = append(o, a)

	return o
}
