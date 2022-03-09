package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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
