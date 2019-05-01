package main

import (
	"html"
	"log"
	"path"

	notify "github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"github.com/jangxx/go-poclient"
	"github.com/skratchdot/open-golang/open"
)

type NotificationActions struct {
	url string
}

var notifier notify.Notifier
var notifications map[uint32]NotificationActions

func initNotifications() {
	notifications = make(map[uint32]NotificationActions)

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Panicln("Error while connecting to session dbus for notifications")
	}

	notifier, err = notify.New(conn)

	if err != nil {
		log.Fatalln(err.Error())
	}

	actions := notifier.ActionInvoked()
	go func() {
		for {
			action := <-actions
			switch action.ActionKey {
			case "default":
				//do nothing, notification was closed
			case "openurl":
				open.Run(notifications[action.ID].url)
			}
			delete(notifications, action.ID)
		}
	}()
}

func sendNotification(message poclient.Message) {
	title := message.AppName
	if message.Title != "" {
		title = message.Title
	}

	body := message.Text
	actions := []string{}
	actionData := NotificationActions{}

	if message.Url != "" {
		label := "Open URL"
		if message.UrlTitle != "" {
			label = message.UrlTitle
		}

		actions = append(actions, "openurl", label)
		actionData.url = message.Url
	}

	n := notify.Notification{
		AppName:       "SchubsiGo",
		ReplacesID:    0,
		Summary:       title,
		Body:          html.EscapeString(body),
		Actions:       actions,
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: int32(5000),
	}

	if config.cache.Exists(message.IconId + ".png") {
		iconPath := path.Join(config.cache.Path, message.IconId+".png")
		n.AppIcon = iconPath
	}

	id, err := notifier.SendNotification(n)

	if err != nil {
		log.Println(err.Error())
	} else {
		notifications[id] = actionData
	}
}
