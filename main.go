package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"

	rice "github.com/GeertJohan/go.rice"
	"github.com/getlantern/systray"

	"github.com/jangxx/go-poclient"
	"github.com/shibukawa/configdir"
)

const APP_NAME string = "SchubsiGo"
const APP_VERSION string = "1.3.0"

var config Config
var server *http.Server
var pushover *poclient.Client
var messages map[int]poclient.Message
var pushover_retry = make(chan bool)
var quit_channel = make(chan bool)

func main() {
	disableTrayIcon := flag.Bool("no-tray", false, "Disables the tray icon")

	flag.Parse()

	if !*disableTrayIcon {
		systray.Run(onReadySystray, onExit)
	} else {
		onReady(true)
	}
}

func onReadySystray() {
	onReady(false)
}

func onReady(noTray bool) {
	iconBox := rice.MustFindBox("./icon")
	icondata := iconBox.MustBytes("icon_64.png")

	if !noTray {
		systray.SetTitle("SchubsiGo")
		systray.SetTooltip("SchubsiGo")
		systray.SetTemplateIcon(icondata, icondata)
	}

	go func() {
		configDirs := configdir.New("literalchaos", "schubsigo")
		config, _ = loadConfig(configDirs, "settings.json")

		if !noTray {
			mOpenWeb := systray.AddMenuItem("Open Web Interface", "Open the web interface")
			mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

			go func() {
				for {
					select {
					case <-mOpenWeb.ClickedCh:
						exec.Command("xdg-open", "http://"+config.Webserver.Addr+":"+config.Webserver.Port).Start()
					case <-mQuit.ClickedCh:
						quit_channel <- true
					}
				}
			}()
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
			if noTray {
				os.Exit(0)
			} else {
				systray.Quit()
			}
		}
	}()

	if noTray {
		select {} // block forever
	}
}

func onExit() {
	if server != nil {
		server.Shutdown(nil)
		log.Printf("Server has shut down")
	}
}
