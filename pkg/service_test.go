package pkg_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takaiyuk/kakeibo-go/mock"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func TestNewService(t *testing.T) {
	cfg := createConfig()
	api := pkg.NewSlackClient(cfg.SlackToken)
	ifttt := pkg.NewIFTTTClient(cfg.IFTTTWebhookToken)
	s := pkg.NewService(api, ifttt)
	assert.Equal(t, api, s.API)
	assert.Equal(t, ifttt, s.IFTTT)
}

func TestGetSlackMessages(t *testing.T) {
	var (
		slackMessages = []*pkg.SlackMessage{
			{Timestamp: 3.0, Text: "test3"},
			{Timestamp: 2.0, Text: "test2"},
			{Timestamp: 1.0, Text: "test1"},
		}
		filterdSlackMessages = []*pkg.SlackMessage{
			{Timestamp: 3.0, Text: "test3"},
			{Timestamp: 2.0, Text: "test2"},
		}
		sortrdSlackMessages = []*pkg.SlackMessage{
			{Timestamp: 2.0, Text: "test2"},
			{Timestamp: 3.0, Text: "test3"},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := mock.NewMockInterfaceSlackClient(ctrl)
	api.EXPECT().FetchMessages("channel_id").Return(slackMessages, nil)
	api.EXPECT().FetchMessages("").Return(nil, errors.New("error: channel_not_found"))
	api.EXPECT().FilterMessages(slackMessages, gomock.Any()).Return(filterdSlackMessages)

	s := pkg.NewService(api, nil)
	options := &pkg.FilterSlackMessagesOptions{}
	var fixtures = []struct {
		desc          string
		cfg           *pkg.Config
		expected      []*pkg.SlackMessage
		expectedError error
	}{
		{desc: "", cfg: createConfig(), expected: sortrdSlackMessages, expectedError: nil},
		{desc: "", cfg: &pkg.Config{}, expected: nil, expectedError: errors.New("error: channel_not_found")},
	}
	for _, tt := range fixtures {
		t.Run(tt.desc, func(t *testing.T) {
			msgs, err := s.GetSlackMessages(tt.cfg, options)
			assert.Equal(t, tt.expected, msgs)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}

func TestPostIFTTTWebhook(t *testing.T) {
	var (
		v1 = 1.0
		v2 = "test1"
	)
	cfg := createConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ifttt := mock.NewMockInterfaceIFTTT(ctrl)
	ifttt.EXPECT().Emit(cfg.IFTTTEventName, strconv.FormatFloat(v1, 'f', -1, 64), v2).Return(nil)
	ifttt.EXPECT().Emit("", strconv.FormatFloat(v1, 'f', -1, 64), v2).Return(errors.New("error: status code 401 Unauthorized"))

	s := pkg.NewService(nil, ifttt)
	msgs := []*pkg.SlackMessage{{Timestamp: v1, Text: v2}}
	var fixtures = []struct {
		desc     string
		cfg      *pkg.Config
		msgs     []*pkg.SlackMessage
		expected error
	}{
		{desc: "a", cfg: cfg, msgs: msgs, expected: nil},
		{desc: "b", cfg: &pkg.Config{}, msgs: msgs, expected: errors.New("error: status code 401 Unauthorized")},
	}
	for _, tt := range fixtures {
		t.Run(tt.desc, func(t *testing.T) {
			err := s.PostIFTTTWebhook(tt.cfg, tt.msgs)
			assert.Equal(t, tt.expected, err)
		})
	}
}

// func TestKakeibo(t *testing.T) {}
