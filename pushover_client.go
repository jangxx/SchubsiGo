package main

import (
	"log"

	"github.com/jangxx/go-poclient"
)

func initPOClient(config Config) *poclient.POClient {
	po := poclient.NewPOClient()

	if config.Userid != "" && config.Usersecret != "" {
		po.RestoreLogin(config.Usersecret, config.Userid)
	}

	if config.Deviceid != "" {
		po.RestoreDevice(config.Deviceid)
	}

	if loggedIn, registered := po.GetStatus(); loggedIn && registered {
		err, _ := po.GetMessages() //get messages to test login

		if err != nil {
			log.Printf("Error while restoring Pushover login: %s\n", err.Error())
			po = poclient.NewPOClient() //start from scratch
		} else {
			log.Printf("Successfully restored Pushover login & device registration")
		}
	}

	return po
}

func resetPOClient() {
	pushover = poclient.NewPOClient()
}

func listenForNotifications(po *poclient.POClient) {
	err, messages := po.GetMessages()

	if err == nil {
		for _, msg := range messages {
			po.Messages <- msg //shove the old messages through the channel
		}

		err := po.DeleteOldMessages(&messages)
		if err != nil {
			log.Printf("Error while deleting old messages: %s\n", err.Error())
		}
	} else {
		log.Printf("Error while fetching old messages: %s\n", err.Error())
	}

	for {
		err := po.ListenForNotifications()
		log.Println(err.Error())

		if _, iserrorframe := err.(*poclient.ErrorFrameError); iserrorframe {
			//relogin in required

		}
	}
}
