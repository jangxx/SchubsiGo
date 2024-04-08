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

type VersionedClient struct {
	Client  *poclient.Client
	Version int
}

func NewVersionedClient(version int) *VersionedClient {
	return &VersionedClient{
		Client:  poclient.New(),
		Version: version,
	}
}

func initPOClient(config Config, maxRetries int) *VersionedClient {
	po := NewVersionedClient(currentClientVersion)

	po.Client.SetAppInfo(APP_NAME, APP_VERSION)

	if config.Userid != "" && config.Usersecret != "" {
		po.Client.RestoreLogin(config.Usersecret, config.Userid)
	}

	if config.Deviceid != "" {
		po.Client.RestoreDevice(config.Deviceid)
	}

	if loggedIn, registered := po.Client.GetStatus(); loggedIn && registered {
		_, err := po.Client.GetMessages() // get messages to test login

		if err != nil {
			urlErr := new(url.Error)

			if errors.As(err, &urlErr) {
				if retryInitialization(po, *urlErr, maxRetries) {
					log.Printf("Successfully restored Pushover login & device registration")
					tray_icon_channel <- 0
				} else {
					po = NewVersionedClient(currentClientVersion) // start from scratch
				}
			} else {
				po = NewVersionedClient(currentClientVersion) // start from scratch
			}
		} else {
			log.Printf("Successfully restored Pushover login & device registration")
			tray_icon_channel <- 0
		}
	}

	return po
}

func retryInitialization(po *VersionedClient, err url.Error, retries int) bool {
	log.Printf("Error while restoring Pushover login: %s\n", err.Error())

	if retries == 0 {
		return false
	}

	log.Println("Error was network related, retry in 5 seconds")
	tray_icon_channel <- 2

	time.Sleep(5 * time.Second)

	_, nextErr := po.Client.GetMessages() // get messages to test login

	nextUrlErr := new(url.Error)

	if nextErr != nil && errors.As(nextErr, &nextUrlErr) {
		return retryInitialization(po, *nextUrlErr, retries-1)
	} else {
		return true
	}
}

func resetPOClient() {
	currentClientVersion++

	if pushover != nil {
		pushover.Client.CloseWebsocket()
		close(pushover.Client.Messages)
	}

	pushover = initPOClient(config, maxConnectRetries)

	go listenForMessages(pushover)
	go listenForNotifications(pushover)
}

func listenForNotifications(po *VersionedClient) {
	messages, err := po.Client.GetMessages()

	if err == nil {
		for _, msg := range messages {
			po.Client.Messages <- msg // shove the old messages through the channel
		}

		err := po.Client.DeleteOldMessages(messages)
		if err != nil {
			log.Printf("Error while deleting old messages: %s\n", err.Error())
		}
	} else {
		log.Printf("Error while fetching old messages: %s\n", err.Error())
	}

	for {
		tray_icon_channel <- 0

		err := po.Client.ListenForNotifications()
		log.Println(err.Error())

		if currentClientVersion != po.Version { // if the client has been destroyed we don't need to handle anything here
			return
		}

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
			tray_icon_channel <- 2
			time.Sleep(15 * time.Second)
		} else if _, registered := po.Client.GetStatus(); registered {
			log.Println("Error was not network related, but we are registered; retry in 15 seconds")
			tray_icon_channel <- 2
			time.Sleep(15 * time.Second)
		} else {
			log.Println("Error was not network related, retry when we get the instruction. Error type: " + reflect.TypeOf(err).String())
			tray_icon_channel <- 1
			// wait until a retry makes sense if the error wasn't network related
			<-pushover_retry
		}

		log.Println("Retrying connection")
	}
}
