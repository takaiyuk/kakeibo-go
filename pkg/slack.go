package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type slackMessage struct {
	Timestamp float64
	Text      string
}

type filterSlackMessagesArgs struct {
	Messages       []*slackMessage
	DtNow          time.Time
	ExcludeDays    int
	ExcludeMinutes int
	IsSort         bool
}

type slackClient struct {
	Token string
}

func newSlackClient(token string) *slackClient {
	return &slackClient{Token: token}
}

func (api *slackClient) getConversationHistory(channelID string) ([]*slackMessage, error) {
	client := new(http.Client)
	endpoint := "https://slack.com/api/conversations.history"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	values := req.URL.Query()
	values.Set("channel", channelID)
	req.URL.RawQuery = values.Encode()
	req.Header.Set("Authorization", "Bearer "+api.Token)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	byteBody, _ := ioutil.ReadAll(res.Body)
	var jsonMap map[string]interface{}
	err = json.Unmarshal(byteBody, &jsonMap)
	if err != nil {
		return nil, err
	}
	slackMessages := []*slackMessage{}
	if jsonMap["ok"] != true {
		return nil, fmt.Errorf("error: %s", jsonMap["error"])
	}
	for _, m := range jsonMap["messages"].([]interface{}) {
		ts, err := strconv.ParseFloat(m.(map[string]interface{})["ts"].(string), 64)
		if err != nil {
			return nil, err
		}
		sm := &slackMessage{
			Timestamp: ts,
			Text:      m.(map[string]interface{})["text"].(string),
		}
		slackMessages = append(slackMessages, sm)
	}
	return slackMessages, nil
}

func (api *slackClient) fetchMessages(channelID string) ([]*slackMessage, error) {
	messages, err := api.getConversationHistory(channelID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (api *slackClient) filterMessages(args filterSlackMessagesArgs) []*slackMessage {
	filteredMessages := []*slackMessage{}
	for _, m := range args.Messages {
		threshold := float64(args.DtNow.AddDate(0, 0, -args.ExcludeDays).Add(time.Minute * -time.Duration(args.ExcludeMinutes)).Unix())
		if m.Timestamp > threshold {
			filteredMessages = append(filteredMessages, m)
		}
	}
	if args.IsSort {
		api.sortMessages(filteredMessages)
	}
	return filteredMessages
}

func (api *slackClient) sortMessages(messages []*slackMessage) {
	sort.Slice(messages, func(i, j int) bool { return messages[i].Timestamp < messages[j].Timestamp })
}

func getSlackMessages(cfg *config) ([]*slackMessage, error) {
	api := newSlackClient(cfg.SlackToken)
	messages, err := api.fetchMessages(cfg.SlackChannelID)
	if err != nil {
		return nil, err
	}
	args := filterSlackMessagesArgs{
		Messages:       messages,
		DtNow:          time.Now(),
		ExcludeDays:    0,
		ExcludeMinutes: 10,
		IsSort:         true,
	}
	filteredMessages := api.filterMessages(args)
	return filteredMessages, nil
}
