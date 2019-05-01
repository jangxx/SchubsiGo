package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/GeertJohan/go.rice"

	"github.com/cratonica/trayhost"
	"github.com/jangxx/go-poclient"
	"github.com/shibukawa/configdir"
)

var config Config
var server *http.Server
var pushover *poclient.Client
var messages map[int]poclient.Message
var pushover_retry = make(chan bool)

func main() {
	iconBox := rice.MustFindBox("./icon")
	icondata := iconBox.MustBytes("icon_64.png")

	runtime.LockOSThread()

	go func() {
		configDirs := configdir.New("literalchaos", "schubsigo")
		config, _ = loadConfig(configDirs, "settings.json")

		trayhost.SetUrl("http://" + config.Webserver.Addr + ":" + config.Webserver.Port)

		messages = make(map[int]poclient.Message)

		pushover = initPOClient(config)

		initNotifications()
		go listenForMessages(pushover)

		go listenForNotifications(pushover)

		server = initWebserver(config.Webserver)
	}()

	trayhost.EnterLoop("SchubsiGo", icondata)
}

func onExit() {
	if server != nil {
		server.Shutdown(nil)
		log.Printf("Server is shutdown")
	}
}
