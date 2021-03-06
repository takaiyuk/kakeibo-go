//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOFILE -package=mock
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

var (
	baseSlackEndpoint = "https://slack.com/api/"
)

type SlackMessage struct {
	Timestamp float64
	Text      string
}

type FilterSlackMessagesOptions struct {
	DtNow          time.Time
	ExcludeDays    int
	ExcludeMinutes int
	IsSort         bool
}

type InterfaceSlackClient interface {
	GetConversationHistory(string) ([]*SlackMessage, error)
	FetchMessages(string) ([]*SlackMessage, error)
	FilterMessages([]*SlackMessage, *FilterSlackMessagesOptions) []*SlackMessage
}

type slackClient struct {
	Token string
}

func NewSlackClient(token string) *slackClient {
	return &slackClient{Token: token}
}

func (api *slackClient) GetConversationHistory(channelID string) ([]*SlackMessage, error) {
	client := new(http.Client)
	endpoint := baseSlackEndpoint + "conversations.history"
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
	slackMessages := []*SlackMessage{}
	if jsonMap["ok"] != true {
		return nil, fmt.Errorf("error: %s", jsonMap["error"])
	}
	for _, m := range jsonMap["messages"].([]interface{}) {
		ts, err := strconv.ParseFloat(m.(map[string]interface{})["ts"].(string), 64)
		if err != nil {
			return nil, err
		}
		sm := &SlackMessage{
			Timestamp: ts,
			Text:      m.(map[string]interface{})["text"].(string),
		}
		slackMessages = append(slackMessages, sm)
	}
	return slackMessages, nil
}

func (api *slackClient) FetchMessages(channelID string) ([]*SlackMessage, error) {
	messages, err := api.GetConversationHistory(channelID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (api *slackClient) FilterMessages(messages []*SlackMessage, options *FilterSlackMessagesOptions) []*SlackMessage {
	filteredMessages := []*SlackMessage{}
	for _, m := range messages {
		threshold := float64(options.DtNow.AddDate(0, 0, -options.ExcludeDays).Add(time.Minute * -time.Duration(options.ExcludeMinutes)).Unix())
		if m.Timestamp > threshold {
			filteredMessages = append(filteredMessages, m)
		}
	}
	if options.IsSort {
		SortMessages(filteredMessages)
	}
	return filteredMessages
}

func SortMessages(messages []*SlackMessage) {
	sort.Slice(messages, func(i, j int) bool { return messages[i].Timestamp < messages[j].Timestamp })
}
