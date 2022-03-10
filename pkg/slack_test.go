package pkg_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

// TODO: replace with mock
func TestSlackClient_getConversationHistory(t *testing.T) {
	envMap, err := pkg.ExportedReadEnv("." + pkg.ExportedEnvFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := pkg.ExportedNewConfig(envMap)
	api := pkg.ExportedNewSlackClient(cfg.SlackToken)

	var fixtures = []struct {
		channelID      string
		expectedLength int
		expectedError  error
	}{
		{channelID: cfg.SlackChannelID, expectedLength: 100, expectedError: nil},
		{channelID: "wrong_channel_id", expectedLength: 0, expectedError: errors.New("error: channel_not_found")},
	}
	for _, tt := range fixtures {
		t.Run(tt.channelID, func(t *testing.T) {
			slackMessages, err := pkg.ExportedGetConversationHistory(api, tt.channelID)
			assert.Equal(t, tt.expectedLength, len(slackMessages))
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestSlackClient_fetchMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	envMap, err := pkg.ExportedReadEnv(envTestFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := pkg.ExportedNewConfig(envMap)
	api := pkg.ExportedNewSlackClient(cfg.SlackToken)

	var fixtures = []struct {
		channelID     string
		patchFunc     func(*pkg.ExportedSlackClient, string) ([]*pkg.ExportedSlackMessage, error)
		expected      []*pkg.ExportedSlackMessage
		expectedError error
	}{
		{
			channelID: cfg.SlackChannelID,
			patchFunc: func(*pkg.ExportedSlackClient, string) ([]*pkg.ExportedSlackMessage, error) {
				messages := []*pkg.ExportedSlackMessage{
					{Timestamp: 2.0, Text: "test2"},
					{Timestamp: 1.0, Text: "test1"},
				}
				return messages, nil
			},
			expected: []*pkg.ExportedSlackMessage{
				{Timestamp: 2.0, Text: "test2"},
				{Timestamp: 1.0, Text: "test1"},
			},
			expectedError: nil,
		},
		{
			channelID: "wrong_channel_id",
			patchFunc: func(*pkg.ExportedSlackClient, string) ([]*pkg.ExportedSlackMessage, error) {
				return nil, errors.New("error: channel_not_found")
			},
			expected:      nil,
			expectedError: errors.New("error: channel_not_found"),
		},
	}
	for _, tt := range fixtures {
		t.Run(tt.channelID, func(t *testing.T) {
			monkey.Patch(pkg.ExportedGetConversationHistory, tt.patchFunc)
			slackMessages, err := pkg.ExportedFetchMessages(api, tt.channelID)
			assert.Equal(t, tt.expected, slackMessages)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestSlackClient_filterMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	api, err := createSlackClient()
	if err != nil {
		t.Fatal(err)
	}

	inputs := []time.Time{
		time.Date(2020, 1, 1, 11, 59, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 11, 51, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 11, 49, 0, 0, time.UTC),
	}
	slackMessages := []*pkg.ExportedSlackMessage{
		{Timestamp: float64(inputs[0].Unix()), Text: "test1"},
		{Timestamp: float64(inputs[1].Unix()), Text: "test2"},
		{Timestamp: float64(inputs[2].Unix()), Text: "test3"},
	}
	args := pkg.ExportedFilterSlackMessagesArgs{
		Messages:       slackMessages,
		DtNow:          time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		ExcludeDays:    0,
		ExcludeMinutes: 10,
		IsSort:         true,
	}
	filteredMessages := pkg.ExportedFilterMessages(api, args)
	// ソートで新しいメッセージが先頭になる
	expected := []*pkg.ExportedSlackMessage{
		{Timestamp: float64(inputs[1].Unix()), Text: "test2"},
		{Timestamp: float64(inputs[0].Unix()), Text: "test1"},
	}
	assert.Equal(t, expected, filteredMessages)
}

func TestSlackClient_sortSlackMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	api, err := createSlackClient()
	if err != nil {
		t.Fatal(err)
	}

	messages := []*pkg.ExportedSlackMessage{
		{Timestamp: 3.0, Text: "test3"},
		{Timestamp: 2.0, Text: "test2"},
		{Timestamp: 1.0, Text: "test1"},
	}
	expected := []*pkg.ExportedSlackMessage{
		{Timestamp: 1.0, Text: "test1"},
		{Timestamp: 2.0, Text: "test2"},
		{Timestamp: 3.0, Text: "test3"},
	}
	pkg.ExportedSortMessages(api, messages)
	assert.Equal(t, expected, messages)
}
