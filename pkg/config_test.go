package pkg_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

const (
	envTestFilePath = "./.env.test"
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

func createConfig() *pkg.Config {
	createEnvFile()
	defer os.Remove(envTestFilePath)

	envMap, _ := pkg.ReadEnv(envTestFilePath)
	cfg := pkg.NewConfig(envMap)
	return cfg
}

func TestReadEnv(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)

	var fixtures = []struct {
		filePath        string
		expected        map[string]string
		expectedIsError bool
	}{
		{
			filePath: envTestFilePath,
			expected: map[string]string{
				"IFTTT_EVENT_NAME":    "event_name",
				"IFTTT_WEBHOOK_TOKEN": "webhook_token",
				"SLACK_TOKEN":         "slack_token",
				"SLACK_CHANNEL_ID":    "channel_id",
			},
			expectedIsError: false,
		},
		{
			filePath:        "wrong_file_path",
			expected:        map[string]string{},
			expectedIsError: true,
		},
	}
	for _, tt := range fixtures {
		t.Run(tt.filePath, func(t *testing.T) {
			envMap, err := pkg.ReadEnv(tt.filePath)
			if tt.expectedIsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, envMap)
		})
	}
}

func TestNewConfig(t *testing.T) {
	createEnvFile()
	defer os.Remove(envTestFilePath)

	envMap, err := pkg.ReadEnv(envTestFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := pkg.NewConfig(envMap)
	assert.Equal(t, "event_name", cfg.IFTTTEventName)
	assert.Equal(t, "webhook_token", cfg.IFTTTWebhookToken)
	assert.Equal(t, "slack_token", cfg.SlackToken)
	assert.Equal(t, "channel_id", cfg.SlackChannelID)
}
