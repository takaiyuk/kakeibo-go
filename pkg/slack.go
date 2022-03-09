package pkg

import (
	"sort"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

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