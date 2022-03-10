package pkg

import (
	"errors"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestSlackClient_fetchMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	envMap, err := readEnv(envTestFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := newConfig(envMap)
	c := newSlackClient(cfg.slackToken)

	var fixtures = []struct {
		channelID     string
		patchFunc     func(*slackClient, string) ([]*slackMessage, error)
		expected      []*slackMessage
		expectedError error
	}{
		{
			channelID: cfg.slackChannelID,
			patchFunc: func(*slackClient, string) ([]*slackMessage, error) {
				messages := []*slackMessage{
					{ts: 2.0, text: "test2"},
					{ts: 1.0, text: "test1"},
				}
				return messages, nil
			},
			expected: []*slackMessage{
				{
					ts:   2.0,
					text: "test2",
				},
				{
					ts:   1.0,
					text: "test1",
				},
			},
			expectedError: nil,
		},
		{
			channelID: "wrong_channel_id",
			patchFunc: func(*slackClient, string) ([]*slackMessage, error) {
				return nil, errors.New("channel_not_found")
			},
			expected:      nil,
			expectedError: errors.New("channel_not_found"),
		},
	}
	for _, tt := range fixtures {
		t.Run(tt.channelID, func(t *testing.T) {
			monkey.Patch(ExportedGetConversationHistory, tt.patchFunc)
			slackMessages, err := c.fetchMessages(tt.channelID)
			assert.Equal(t, tt.expected, slackMessages)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestSlackClient_filterMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	c, err := createSlackClient()
	if err != nil {
		t.Fatal(err)
	}

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
	args := filterSlackMessagesArgs{
		messages:       slackMessages,
		dtNow:          time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		excludeDays:    0,
		excludeMinutes: 10,
		isSort:         true,
	}
	filteredMessages := c.filterMessages(args)
	// ソートで新しいメッセージが先頭になる
	expected := []*slackMessage{
		{ts: float64(inputs[1].Unix()), text: "test2"},
		{ts: float64(inputs[0].Unix()), text: "test1"},
	}
	assert.Equal(t, expected, filteredMessages)
}

func TestSlackClient_sortSlackMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	c, err := createSlackClient()
	if err != nil {
		t.Fatal(err)
	}

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
	c.sortMessages(messages)
	assert.Equal(t, expected, messages)
}
