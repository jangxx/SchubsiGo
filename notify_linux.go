package main

import (
	"html"
	"log"
	"path"
	"time"

	notify "github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
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

	notifier, err = notify.New(conn, notify.WithOnAction(func(action *notify.ActionInvokedSignal) {
		switch action.ActionKey {
		case "openurl":
			open.Run(notifications[action.ID].url)
		}
		delete(notifications, action.ID)
	}))

	if err != nil {
		log.Fatalln(err.Error())
	}
}

func sendNotification(message poclient.Message) {
	title := message.AppName
	if message.Title != "" {
		title = message.Title
	}

	body := message.Text
	actions := []notify.Action{}
	actionData := NotificationActions{}

	if message.URL != "" {
		label := "Open URL"
		if message.URLTitle != "" {
			label = message.URLTitle
		}

		actions = append(actions, notify.Action{
			Label: label,
			Key:   "openurl",
		})
		actionData.url = message.URL
	}

	n := notify.Notification{
		AppName:       "SchubsiGo",
		Actions:       actions,
		ReplacesID:    0,
		Summary:       title,
		Body:          html.EscapeString(body),
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: time.Second * 5,
	}

	if config.cache.Exists(message.IconID + ".png") {
		iconPath := path.Join(config.cache.Path, message.IconID+".png")
		n.AppIcon = iconPath
	}

	id, err := notifier.SendNotification(n)

	if err != nil {
		log.Println(err.Error())
	} else {
		notifications[id] = actionData
	}
}

func sendStatusNotification(status string) {
	n := notify.Notification{
		AppName:    "SchubsiGo",
		ReplacesID: 0,
		Body:       status,
	}

	_, err := notifier.SendNotification(n)

	if err != nil {
		log.Println(err.Error())
	}
}
