package pkg

import (
	"log"
)

func Kakeibo() {
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
