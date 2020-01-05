package main

import (
	"log"
	"net"
	"reflect"
	"time"

	"github.com/jangxx/go-poclient"
)

func initPOClient(config Config) *poclient.Client {
	po := poclient.New()

	if config.Userid != "" && config.Usersecret != "" {
		po.RestoreLogin(config.Usersecret, config.Userid)
	}

	if config.Deviceid != "" {
		po.RestoreDevice(config.Deviceid)
	}

	if loggedIn, registered := po.GetStatus(); loggedIn && registered {
		_, err := po.GetMessages() //get messages to test login

		if err != nil {
			log.Printf("Error while restoring Pushover login: %s\n", err.Error())
			po = poclient.New() //start from scratch
		} else {
			log.Printf("Successfully restored Pushover login & device registration")
		}
	}

	return po
}

func resetPOClient() {
	pushover.CloseWebsocket()
	pushover = poclient.New()
}

func listenForNotifications(po *poclient.Client) {
	messages, err := po.GetMessages()

	if err == nil {
		for _, msg := range messages {
			po.Messages <- msg // shove the old messages through the channel
		}

		err := po.DeleteOldMessages(messages)
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
			// reset user config
			config.Userid = ""
			config.Usersecret = ""
			config.Deviceid = ""
			config.Display_Devicename = ""
			config.Display_Username = ""

			sendStatusNotification("A permanent error occured. You need to re-login after the application is closed.")
		}

		if _, isneterror := err.(net.Error); isneterror {
			log.Println("Error was network related, retry in 15 seconds")
			// don't wait for pushover_retry if the error was network related, instead wait for 15 seconds
			time.Sleep(15 * time.Second)
		} else if _, registered := po.GetStatus(); registered {
			log.Println("Error was not network related, but we are registered; retry in 15 seconds")
			time.Sleep(15 * time.Second)
		} else {
			log.Println("Error was not network related, retry when we get the instruction. Error type: " + reflect.TypeOf(err).String())
			// wait until a retry makes sense if the error wasn't network related
			<-pushover_retry
		}

		log.Println("Retrying connection")
	}
}
