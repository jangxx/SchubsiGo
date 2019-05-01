package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/cratonica/trayhost"
	"github.com/jangxx/go-poclient"
	"github.com/shibukawa/configdir"
)

var config Config
var server *http.Server
var pushover *poclient.Client
var messages map[int]poclient.Message

func main() {
	icondata, err := ioutil.ReadFile(filepath.FromSlash("./icon/icon.png"))
	if err != nil {
		log.Fatal(err)
	}

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
