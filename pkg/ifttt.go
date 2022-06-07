//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOFILE -package=mock
package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	baseIFTTTEndpoint = "https://maker.ifttt.com/trigger/"
)

type InterfaceIFTTT interface {
	Emit(string, ...string) error
	Post(string, string, io.Reader) error
}

// https://github.com/domnikl/ifttt-webhook
type IFTTTClient struct {
	APIKey string
}

func NewIFTTTClient(apiKey string) *IFTTTClient {
	return &IFTTTClient{APIKey: apiKey}
}

func (i *IFTTTClient) Emit(eventName string, v ...string) error {
	url := baseIFTTTEndpoint + eventName + "/with/key/" + i.APIKey
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
	err = i.Post(url, "application/json", bytes.NewReader(body))
	return err
}

func (i *IFTTTClient) Post(url, contentType string, body io.Reader) error {
	res, err := http.Post(url, contentType, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error: status code %s", res.Status)
	}
	return nil
}
