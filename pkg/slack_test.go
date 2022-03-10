package pkg_test

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func createSlackClient() (*pkg.SlackClient, error) {
	envMap, err := pkg.ReadEnv(envTestFilePath)
	if err != nil {
		return nil, err
	}
	cfg := pkg.NewConfig(envMap)
	api := pkg.NewSlackClient(cfg.SlackToken)
	return api, nil
}
func TestSlackClient_FetchMessages(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)
	envMap, err := pkg.ReadEnv(envTestFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := pkg.NewConfig(envMap)
	api := pkg.NewSlackClient(cfg.SlackToken)

	var fixtures = []struct {
		channelID     string
		patchFunc     func(*pkg.SlackClient, string) ([]*pkg.SlackMessage, error)
		expected      []*pkg.SlackMessage
		expectedError error
	}{
		{
			channelID: "correct_channel_id",
			patchFunc: func(*pkg.SlackClient, string) ([]*pkg.SlackMessage, error) {
				messages := []*pkg.SlackMessage{
					{Timestamp: 2.0, Text: "test2"},
					{Timestamp: 1.0, Text: "test1"},
				}
				return messages, nil
			},
			expected: []*pkg.SlackMessage{
				{Timestamp: 2.0, Text: "test2"},
				{Timestamp: 1.0, Text: "test1"},
			},
			expectedError: nil,
		},
		{
			channelID: "wrong_channel_id",
			patchFunc: func(*pkg.SlackClient, string) ([]*pkg.SlackMessage, error) {
				return nil, errors.New("error: channel_not_found")
			},
			expected:      nil,
			expectedError: errors.New("error: channel_not_found"),
		},
	}
	for _, tt := range fixtures {
		t.Run(tt.channelID, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(api), "GetConversationHistory", tt.patchFunc)
			slackMessages, err := api.FetchMessages(tt.channelID)
			assert.Equal(t, tt.expected, slackMessages)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestSlackClient_FilterMessages(t *testing.T) {
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
	slackMessages := []*pkg.SlackMessage{
		{Timestamp: float64(inputs[0].Unix()), Text: "test1"},
		{Timestamp: float64(inputs[1].Unix()), Text: "test2"},
		{Timestamp: float64(inputs[2].Unix()), Text: "test3"},
	}
	options := &pkg.FilterSlackMessagesOptions{
		DtNow:          time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		ExcludeDays:    0,
		ExcludeMinutes: 10,
		IsSort:         true,
	}
	filteredMessages := api.FilterMessages(slackMessages, options)
	// ソートで新しいメッセージが先頭になる
	expected := []*pkg.SlackMessage{
		{Timestamp: float64(inputs[1].Unix()), Text: "test2"},
		{Timestamp: float64(inputs[0].Unix()), Text: "test1"},
	}
	assert.Equal(t, expected, filteredMessages)
}

func TestSortSlackMessages(t *testing.T) {
	messages := []*pkg.SlackMessage{
		{Timestamp: 3.0, Text: "test3"},
		{Timestamp: 2.0, Text: "test2"},
		{Timestamp: 1.0, Text: "test1"},
	}
	expected := []*pkg.SlackMessage{
		{Timestamp: 1.0, Text: "test1"},
		{Timestamp: 2.0, Text: "test2"},
		{Timestamp: 3.0, Text: "test3"},
	}
	pkg.SortMessages(messages)
	assert.Equal(t, expected, messages)
}
