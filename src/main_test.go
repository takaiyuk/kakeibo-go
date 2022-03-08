package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

const (
	envFilePath = "./.env.test"
)

func createEnvFile() {
	env := `IFTTT_EVENT_NAME=event_name
IFTTT_WEBHOOK_TOKEN=webhook_token
SLACK_TOKEN=slack_token
SLACK_CHANNEL_ID=channel_id
`
	err := ioutil.WriteFile(envFilePath, []byte(env), 0644)
	if err != nil {
		panic(err)
	}
}

func TestReadEnv(t *testing.T) {
	createEnvFile()
	defer os.Remove(envFilePath)

	envMap, err := readEnv(envFilePath)
	assert.NoError(t, err)
	assert.Equal(t, "event_name", envMap["IFTTT_EVENT_NAME"])
	assert.Equal(t, "webhook_token", envMap["IFTTT_WEBHOOK_TOKEN"])
	assert.Equal(t, "slack_token", envMap["SLACK_TOKEN"])
	assert.Equal(t, "channel_id", envMap["SLACK_CHANNEL_ID"])
}

func TestNewConfig(t *testing.T) {
	envMap := make(map[string]string)
	envMap["IFTTT_EVENT_NAME"] = "event_name"
	envMap["IFTTT_WEBHOOK_TOKEN"] = "webhook_token"
	envMap["SLACK_TOKEN"] = "slack_token"
	envMap["SLACK_CHANNEL_ID"] = "channel_id"
	cfg := newConfig(envMap)
	assert.Equal(t, "event_name", cfg.iftttEventName)
	assert.Equal(t, "webhook_token", cfg.iftttWebhookToken)
	assert.Equal(t, "slack_token", cfg.slackToken)
	assert.Equal(t, "channel_id", cfg.slackChannelID)
}

func TestGetSlackMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envFilePath)

	envMap, err := readEnv(envFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := newConfig(envMap)
	monkey.Patch(getSlackConversationHistory, func(*config) []slack.Message {
		return []slack.Message{
			{
				Msg: slack.Msg{
					Timestamp: "2.0",
					Text:      "test2",
				},
			},
			{
				Msg: slack.Msg{
					Timestamp: "1.0",
					Text:      "test1",
				},
			},
		}
	})
	slackMessages := getSlackMessages(cfg)
	expected := []*slackMessage{
		{
			ts:   2.0,
			text: "test2",
		},
		{
			ts:   1.0,
			text: "test1",
		},
	}
	assert.Equal(t, expected, slackMessages)
}

func TestFilterSlackMessages(t *testing.T) {
	inputs := []time.Time{
		time.Date(2020, 1, 1, 11, 59, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 11, 51, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 11, 49, 0, 0, time.UTC),
	}
	slackMessages := []*slackMessage{
		{ts: float64(inputs[0].Unix()), text: "test1"},
		{ts: float64(inputs[1].Unix()), text: "test2"},
		{ts: float64(inputs[2].Unix()), text: "test3"},
	}
	dtNow := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	args := filterSlackMessagesArgs{
		messages:       slackMessages,
		dtNow:          dtNow,
		excludeDays:    0,
		excludeMinutes: 10,
		isSort:         true,
	}
	filteredMessages := filterSlackMessages(args)
	// ソートで新しいメッセージが先頭になる
	expected := []*slackMessage{
		{ts: float64(inputs[1].Unix()), text: "test2"},
		{ts: float64(inputs[0].Unix()), text: "test1"},
	}
	assert.Equal(t, expected, filteredMessages)
}

func TestSortSlackMessages(t *testing.T) {
	messages := []*slackMessage{
		{ts: 3.0, text: "test3"},
		{ts: 2.0, text: "test2"},
		{ts: 1.0, text: "test1"},
	}
	expected := []*slackMessage{
		{ts: 1.0, text: "test1"},
		{ts: 2.0, text: "test2"},
		{ts: 3.0, text: "test3"},
	}
	sortSlackMessages(messages)
	assert.Equal(t, expected, messages)
}

func TestNewIFTTT(t *testing.T) {
	apiKey := "key"
	i := newIFTTT(apiKey)
	expected := &ifttt{apiKey: apiKey}
	assert.Equal(t, expected, i)
}
