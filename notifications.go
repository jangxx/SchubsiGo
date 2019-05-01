package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jangxx/go-poclient"
)

func listenForMessages(po *poclient.POClient) {
	for {
		select {
		case message := <-po.Messages:
			messages[message.UniqueId] = message

			if !config.cache.Exists(message.IconId + ".png") {
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
	resp, err := http.Get(message.GetIconUrl())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	config.cache.WriteFile(message.IconId+".png", body)
	return nil
}
