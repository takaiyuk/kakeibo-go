//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOFILE -package=mock
package pkg

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	excludeDays    = 0
	excludeMinutes = 10
	isSort         = true
)

type InterfaceService interface {
	GetSlackMessages(*Config, *FilterSlackMessagesOptions) ([]*SlackMessage, error)
	PostIFTTTWebhook(*Config, []*SlackMessage) error
}

type service struct {
	API   InterfaceSlackClient
	IFTTT InterfaceIFTTT
}

func NewService(api InterfaceSlackClient, ifttt InterfaceIFTTT) *service {
	return &service{API: api, IFTTT: ifttt}
}

func (s *service) GetSlackMessages(cfg *Config, options *FilterSlackMessagesOptions) ([]*SlackMessage, error) {
	messages, err := s.API.FetchMessages(cfg.SlackChannelID)
	if err != nil {
		return nil, err
	}
	messages = s.API.FilterMessages(messages, options)
	SortMessages(messages)
	return messages, nil
}

func (s *service) PostIFTTTWebhook(cfg *Config, messages []*SlackMessage) error {
	for _, m := range messages {
		err := s.IFTTT.Emit(cfg.IFTTTEventName, strconv.FormatFloat(m.Timestamp, 'f', -1, 64), m.Text)
		if err != nil {
			return err
		}
		fmt.Printf("message to be posted: %f,%s\n", m.Timestamp, m.Text)
	}
	return nil
}

func Kakeibo() {
	envMap, err := ReadEnv(EnvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	cfg := NewConfig(envMap)
	api := NewSlackClient(cfg.SlackToken)
	ifttt := NewIFTTTClient(cfg.IFTTTWebhookToken)
	s := NewService(api, ifttt)
	filterSlackMessagesOptions := &FilterSlackMessagesOptions{
		DtNow:          time.Now(),
		ExcludeDays:    excludeDays,
		ExcludeMinutes: excludeMinutes,
		IsSort:         isSort,
	}
	filteredMessages, err := s.GetSlackMessages(cfg, filterSlackMessagesOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = s.PostIFTTTWebhook(cfg, filteredMessages)
	if err != nil {
		log.Fatal(err)
	}
}
