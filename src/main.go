package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

const (
	envFilePath = "./.env"
)

type config struct {
	slackToken        string
	slackChannelID    string
	iftttWebhookToken string
	iftttEventName    string
}

func newConfig(envMap map[string]string) *config {
	return &config{
		slackToken:        envMap["SLACK_TOKEN"],
		slackChannelID:    envMap["SLACK_CHANNEL_ID"],
		iftttWebhookToken: envMap["IFTTT_WEBHOOK_TOKEN"],
		iftttEventName:    envMap["IFTTT_EVENT_NAME"],
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

type slackMessage struct {
	ts   float64
	text string
}

type filterSlackMessagesArgs struct {
	messages       []*slackMessage
	dtNow          time.Time
	excludeDays    int
	excludeMinutes int
	isSort         bool
}

type slackClient struct {
	token string
}

func newSlackClient(token string) *slackClient {
	return &slackClient{token: token}
}

func (c *slackClient) getConversationHistory(channelID string) ([]slack.Message, error) {
	api := slack.New(c.token)
	params := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
	}
	history, err := api.GetConversationHistory(&params)
	if err != nil {
		return nil, err
	}
	return history.Messages, nil
}

func (c *slackClient) fetchMessages(channelID string) ([]*slackMessage, error) {
	messages, err := c.getConversationHistory(channelID)
	if err != nil {
		return nil, err
	}
	slackMessages := []*slackMessage{}
	for _, m := range messages {
		ts, _ := strconv.ParseFloat(m.Timestamp, 64)
		sm := &slackMessage{
			ts:   ts,
			text: m.Text,
		}
		slackMessages = append(slackMessages, sm)
	}
	return slackMessages, nil
}

func (c *slackClient) filterMessages(args filterSlackMessagesArgs) []*slackMessage {
	filteredMessages := []*slackMessage{}
	for _, m := range args.messages {
		threshold := float64(args.dtNow.AddDate(0, 0, -args.excludeDays).Add(time.Minute * -time.Duration(args.excludeMinutes)).Unix())
		if m.ts > threshold {
			filteredMessages = append(filteredMessages, m)
		}
	}
	if args.isSort {
		c.sortMessages(filteredMessages)
	}
	return filteredMessages
}

func (c *slackClient) sortMessages(messages []*slackMessage) {
	sort.Slice(messages, func(i, j int) bool { return messages[i].ts < messages[j].ts })
}

func getSlackMessages(cfg *config) ([]*slackMessage, error) {
	c := newSlackClient(cfg.slackToken)
	messages, err := c.fetchMessages(cfg.slackChannelID)
	if err != nil {
		return nil, err
	}
	args := filterSlackMessagesArgs{
		messages:       messages,
		dtNow:          time.Now(),
		excludeDays:    0,
		excludeMinutes: 10,
		isSort:         true,
	}
	filteredMessages := c.filterMessages(args)
	return filteredMessages, nil
}

// https://github.com/domnikl/ifttt-webhook
type ifttt struct {
	apiKey string
}

func newIFTTT(apiKey string) *ifttt {
	return &ifttt{apiKey: apiKey}
}

func (i *ifttt) post(eventName string, v ...string) error {
	url := "https://maker.ifttt.com/trigger/" + eventName + "/with/key/" + i.apiKey
	values := map[string]string{}
	for x, value := range v {
		values["value"+strconv.Itoa(x+1)] = value
		// only include up to 3 values
		if x == 2 {
			log.Printf("only 3 values are allowed. argument %d (%s) and after that are ignored.", x+1, value)
			break
		}
	}
	body, err := json.Marshal(values)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func postIFTTTWebhook(cfg *config, messages []*slackMessage) error {
	i := newIFTTT(cfg.iftttWebhookToken)
	for _, m := range messages {
		err := i.post(cfg.iftttEventName, strconv.FormatFloat(m.ts, 'f', -1, 64), m.text)
		if err != nil {
			return err
		}
		fmt.Printf("message to be posted: %f,%s\n", m.ts, m.text)
	}
	return nil
}

func kakeibo() {
	envMap, err := readEnv(envFilePath)
	if err != nil {
		log.Fatal(err)
	}
	cfg := newConfig(envMap)
	filteredMessages, err := getSlackMessages(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = postIFTTTWebhook(cfg, filteredMessages)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		// Run locally
		kakeibo()
	} else {
		lambda.Start(kakeibo)
	}
}
