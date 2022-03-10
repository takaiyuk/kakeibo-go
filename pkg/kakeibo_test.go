package pkg_test

import (
	"io/ioutil"

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

func createSlackClient() (*pkg.ExportedSlackClient, error) {
	envMap, err := pkg.ExportedReadEnv(envTestFilePath)
	if err != nil {
		return nil, err
	}
	cfg := pkg.ExportedNewConfig(envMap)
	api := pkg.ExportedNewSlackClient(cfg.SlackToken)
	return api, nil
}
