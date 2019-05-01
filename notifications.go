package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jangxx/go-poclient"
)

func listenForMessages(po *poclient.Client) {
	for {
		select {
		case message := <-po.Messages:
			messages[message.UniqueID] = message

			if !config.cache.Exists(message.IconID + ".png") {
				//download icon first
				err := downloadMessageIcon(message)
				if err != nil {
					log.Println("Error while downloading icon: " + err.Error())
				}
			}
			sendNotification(message)
		}
	}
}

func downloadMessageIcon(message poclient.Message) error {
	resp, err := http.Get(getMessageIconURL(message))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	config.cache.WriteFile(message.IconID+".png", body)
	return nil
}

func getMessageIconURL(message poclient.Message) string {
	return fmt.Sprintf("https://api.pushover.net/icons/%s.png", message.IconID)
}
