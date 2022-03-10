package pkg

import (
	"io/ioutil"
)

const (
	envTestFilePath = "./.env.test"
)

var (
	ExportedGetConversationHistory = (*slackClient).getConversationHistory
)

func createEnvFile() {
	env := `IFTTT_EVENT_NAME=event_name
IFTTT_WEBHOOK_TOKEN=webhook_token
SLACK_TOKEN=slack_token
SLACK_CHANNEL_ID=channel_id
`
	err := ioutil.WriteFile(envTestFilePath, []byte(env), 0644)
	if err != nil {
		panic(err)
	}
}

func createSlackClient() (*slackClient, error) {
	envMap, err := readEnv(envTestFilePath)
	if err != nil {
		return nil, err
	}
	cfg := newConfig(envMap)
	c := newSlackClient(cfg.slackToken)
	return c, nil
}
