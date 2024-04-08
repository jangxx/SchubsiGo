package main

import (
	"context"
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
const APP_VERSION string = "1.5.0"

var config Config
var server *http.Server
var pushover *VersionedClient
var messages map[int]poclient.Message
var pushover_retry = make(chan bool)
var quit_channel = make(chan bool)
var tray_icon_channel = make(chan int)
var maxConnectRetries int
var currentClientVersion int = 0

func main() {
	log.Printf("Starting %s %s\n", APP_NAME, APP_VERSION)

	disableTrayIcon := flag.Bool("no-tray", false, "Disables the tray icon")
	maxConnectRetriesFlag := flag.Int("max-connect-retries", 5, "Maximum number of retries to connect to the Pushover API")

	flag.Parse()

	maxConnectRetries = *maxConnectRetriesFlag
	log.Printf("Max connect retries: %d\n", maxConnectRetries)

	if maxConnectRetries < 0 {
		log.Fatalf("Invalid value for max-connect-retries: %d\n", maxConnectRetries)
	}

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
	iconBlueData := iconBox.MustBytes("icon_64.png")
	iconRedData := iconBox.MustBytes("icon_red_64.png")
	iconYellowData := iconBox.MustBytes("icon_yellow_64.png")

	if !noTray {
		systray.SetTitle("SchubsiGo")
		systray.SetTooltip("SchubsiGo")
		systray.SetTemplateIcon(iconRedData, iconRedData)
	}

	var mResetConnection *systray.MenuItem

	go func() {
		configDirs := configdir.New("literalchaos", "schubsigo")
		config, _ = loadConfig(configDirs, "settings.json")

		if !noTray {
			mVersion := systray.AddMenuItem("Version: "+APP_VERSION, "Version of the app")
			mVersion.Disable()

			mOpenWeb := systray.AddMenuItem("Open Web Interface", "Open the web interface")

			mResetConnection = systray.AddMenuItem("Retry Pushover connection", "Retry connecting to the Pushover API")
			mResetConnection.Hide()

			mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

			go func() {
				for {
					select {
					case <-mOpenWeb.ClickedCh:
						exec.Command("xdg-open", "http://"+config.Webserver.Addr+":"+config.Webserver.Port).Start()
					case <-mQuit.ClickedCh:
						quit_channel <- true
					case <-mResetConnection.ClickedCh:
						resetPOClient()
					}
				}
			}()
		}

		messages = make(map[int]poclient.Message)

		initNotifications()
		server = initWebserver(config.Webserver)

		resetPOClient()
	}()

	go func() {
		for {
			select {
			case <-quit_channel:
				if noTray {
					onExit()
					os.Exit(0)
				} else {
					systray.Quit()
				}
			case trayIcon := <-tray_icon_channel:
				if !noTray {
					switch trayIcon {
					case 0:
						systray.SetTemplateIcon(iconBlueData, iconBlueData)
						mResetConnection.Hide()
					case 1:
						systray.SetTemplateIcon(iconRedData, iconRedData)
						mResetConnection.Show()
					case 2:
						systray.SetTemplateIcon(iconYellowData, iconYellowData)
						mResetConnection.Hide()
					default:
						log.Printf("Unknown tray icon: %d\n", trayIcon)
					}
				}
			}
		}
	}()

	if noTray {
		select {} // block forever
	}
}

func onExit() {
	if server != nil {
		server.Shutdown(context.Background())
		log.Printf("Server has shut down")
	}
}
