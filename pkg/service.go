package pkg

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type InterfaceService interface {
	GetSlackMessages(*Config, *FilterSlackMessagesOptions) ([]*SlackMessage, error)
	PostIFTTTWebhook(*Config, []*SlackMessage) error
}

type Service struct {
	API   InterfaceSlackClient
	IFTTT InterfaceIFTTT
}

func NewService(api InterfaceSlackClient, ifttt InterfaceIFTTT) *Service {
	return &Service{API: api, IFTTT: ifttt}
}

func (s *Service) GetSlackMessages(cfg *Config, options *FilterSlackMessagesOptions) ([]*SlackMessage, error) {
	messages, err := s.API.FetchMessages(cfg.SlackChannelID)
	if err != nil {
		return nil, err
	}
	messages = s.API.FilterMessages(messages, options)
	SortMessages(messages)
	return messages, nil
}

func (s *Service) PostIFTTTWebhook(cfg *Config, messages []*SlackMessage) error {
	for _, m := range messages {
		err := s.IFTTT.Post(cfg.IFTTTEventName, strconv.FormatFloat(m.Timestamp, 'f', -1, 64), m.Text)
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
		ExcludeDays:    0,
		ExcludeMinutes: 10,
		IsSort:         true,
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
