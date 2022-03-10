package pkg

import (
	"io/ioutil"
	"strings"
)

const (
	envFilePath = "./.env"
)

type config struct {
	SlackToken        string
	SlackChannelID    string
	IFTTTWebhookToken string
	IFTTTEventName    string
}

func newConfig(envMap map[string]string) *config {
	return &config{
		SlackToken:        envMap["SLACK_TOKEN"],
		SlackChannelID:    envMap["SLACK_CHANNEL_ID"],
		IFTTTWebhookToken: envMap["IFTTT_WEBHOOK_TOKEN"],
		IFTTTEventName:    envMap["IFTTT_EVENT_NAME"],
	}
}

func readEnv(filePath string) (map[string]string, error) {
	envMap := make(map[string]string)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return envMap, err
	}
	contentString := string(content)
	lines := strings.Split(contentString, "\n")
	for _, line := range lines {
		// skip empty line
		if len(line) == 0 {
			continue
		}
		pair := strings.Split(strings.TrimSpace(line), "=")
		envMap[pair[0]] = pair[1]
	}
	return envMap, nil
}
