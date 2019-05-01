package main

import (
	"encoding/json"

	"github.com/shibukawa/configdir"
)

type Config struct {
	configDir configdir.ConfigDir
	cache     *configdir.Config

	Userid     string `json:"user_id"`
	Usersecret string `json:"user_secret"`
	Deviceid   string `json:"device_id"`

	Display_Username   string `json:"display_username"`
	Display_Devicename string `json:"display_devicename"`

	Webserver WebserverConfig `json:"webserver"`
}

type WebserverConfig struct {
	Addr string `json:"address"`
	Port string `json:"port"`
}

var DefaultConfig Config = Config{
	Webserver: WebserverConfig{
		Addr: "localhost",
		Port: "33322",
	},
}

func loadConfig(dir configdir.ConfigDir, filename string) (Config, error) {
	folder := dir.QueryFolderContainsFile(filename)
	cache := dir.QueryCacheFolder()

	config := DefaultConfig
	if folder != nil {
		data, err := folder.ReadFile(filename)
		if err != nil {
			return config, err
		}

		err = json.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}
	}

	config.configDir = dir
	config.cache = cache
	return config, nil
}

func storeConfig(config Config, filename string) error {
	folders := config.configDir.QueryFolders(configdir.Global)

	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = folders[0].WriteFile(filename, data)
	return err
}
