package pkg_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

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
			envMap, err := pkg.ExportedReadEnv(tt.filePath)
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

	envMap, err := pkg.ExportedReadEnv(envTestFilePath)
	if err != nil {
		t.Fatal(err)
	}
	cfg := pkg.ExportedNewConfig(envMap)
	assert.Equal(t, "event_name", cfg.IFTTTEventName)
	assert.Equal(t, "webhook_token", cfg.IFTTTWebhookToken)
	assert.Equal(t, "slack_token", cfg.SlackToken)
	assert.Equal(t, "channel_id", cfg.SlackChannelID)
}
