package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type InterfaceIFTTT interface {
	Post(string, ...string) error
}

// https://github.com/domnikl/ifttt-webhook
type IFTTT struct {
	APIKey string
}

func NewIFTTT(apiKey string) *IFTTT {
	return &IFTTT{APIKey: apiKey}
}

func (i *IFTTT) Post(eventName string, v ...string) error {
	url := "https://maker.ifttt.com/trigger/" + eventName + "/with/key/" + i.APIKey
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
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error: status code %s", res.Status)
	}
	return nil
}
