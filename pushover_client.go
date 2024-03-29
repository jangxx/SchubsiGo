package main

import (
	"errors"
	"log"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/jangxx/go-poclient"
)

func initPOClient(config Config, maxRetries int) *poclient.Client {
	po := poclient.New()

	po.SetAppInfo(APP_NAME, APP_VERSION)

	if config.Userid != "" && config.Usersecret != "" {
		po.RestoreLogin(config.Usersecret, config.Userid)
	}

	if config.Deviceid != "" {
		po.RestoreDevice(config.Deviceid)
	}

	if loggedIn, registered := po.GetStatus(); loggedIn && registered {
		_, err := po.GetMessages() // get messages to test login

		if err != nil {
			urlErr := new(url.Error)

			if errors.As(err, &urlErr) {
				if retryInitialization(po, *urlErr, maxRetries) {
					log.Printf("Successfully restored Pushover login & device registration")
				} else {
					po = poclient.New() // start from scratch
				}
			} else {
				po = poclient.New() // start from scratch
			}
		} else {
			log.Printf("Successfully restored Pushover login & device registration")
		}
	}

	return po
}

func retryInitialization(po *poclient.Client, err url.Error, retries int) bool {
	log.Printf("Error while restoring Pushover login: %s\n", err.Error())

	if retries == 0 {
		return false
	}

	log.Println("Error was network related, retry in 5 seconds")
	time.Sleep(5 * time.Second)

	_, nextErr := po.GetMessages() // get messages to test login

	nextUrlErr := new(url.Error)

	if nextErr != nil && errors.As(nextErr, &nextUrlErr) {
		return retryInitialization(po, *nextUrlErr, retries-1)
	} else {
		return true
	}
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

			err = storeConfig(config, "settings.json")
			if err != nil {
				log.Printf("Error while resetting user config: %s\n", err.Error())
			}

			sendStatusNotification("A permanent error occured. You need to re-login after the application is started again.")
			log.Println("The pushover server sent an error frame; qutting the application now")
			quit_channel <- true
			return
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
