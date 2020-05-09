package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	rice "github.com/GeertJohan/go.rice"

	"github.com/cratonica/trayhost"
	"github.com/jangxx/go-poclient"
	"github.com/shibukawa/configdir"
)

const APP_NAME string = "SchubsiGo"
const APP_VERSION string = "1.2.0"

var config Config
var server *http.Server
var pushover *poclient.Client
var messages map[int]poclient.Message
var pushover_retry = make(chan bool)
var quit_channel = make(chan bool)

func main() {
	iconBox := rice.MustFindBox("./icon")
	icondata := iconBox.MustBytes("icon_64.png")

	disableTrayIcon := flag.Bool("no-tray", false, "Disables the tray icon")

	flag.Parse()

	runtime.LockOSThread()

	go func() {
		configDirs := configdir.New("literalchaos", "schubsigo")
		config, _ = loadConfig(configDirs, "settings.json")

		if !*disableTrayIcon {
			trayhost.SetUrl("http://" + config.Webserver.Addr + ":" + config.Webserver.Port)
		}

		messages = make(map[int]poclient.Message)

		pushover = initPOClient(config)

		initNotifications()
		go listenForMessages(pushover)

		go listenForNotifications(pushover)

		server = initWebserver(config.Webserver)
	}()

	go func() {
		select {
		case <-quit_channel:
			os.Exit(0)
		}
	}()

	if !*disableTrayIcon {
		trayhost.EnterLoop("SchubsiGo", icondata)
	} else {
		select {} // block forever
	}
}

func onExit() {
	if server != nil {
		server.Shutdown(nil)
		log.Printf("Server is shutdown")
	}
}
