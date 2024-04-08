package main

import (
	"encoding/json"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/jangxx/go-poclient"

	"github.com/gorilla/mux"
)

func initWebserver(config WebserverConfig) *http.Server {
	r := mux.NewRouter()
	buildBox := rice.MustFindBox("./webinterface/dist/")

	r.HandleFunc("/", IndexRoute)

	r.HandleFunc("/login", serveSingleFile(buildBox, "index.html")).Methods("GET")
	r.HandleFunc("/register-device", serveSingleFile(buildBox, "index.html")).Methods("GET")
	r.HandleFunc("/info", serveSingleFile(buildBox, "index.html")).Methods("GET")

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/login", LoginRoute).Methods("POST")
	apiRouter.HandleFunc("/register", RegisterRoute).Methods("POST")
	apiRouter.HandleFunc("/logout", LogoutRoute).Methods("POST")
	apiRouter.HandleFunc("/userinfo", UserinfoRoute).Methods("GET")
	apiRouter.HandleFunc("/quit", QuitRoute).Methods("POST")
	apiRouter.HandleFunc("/reset-connection", ResetConnectionRoute).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(buildBox.HTTPBox()))

	srv := &http.Server{
		Handler: r,
		Addr:    config.Addr + ":" + config.Port,
	}

	log.Printf("Server is listening on http://%s:%s\n", config.Addr, config.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server exited with error: %s\n", err)
		}
	}()

	return srv
}

func IndexRoute(resp http.ResponseWriter, req *http.Request) {
	if loggedin, registered := pushover.Client.GetStatus(); !loggedin {
		http.Redirect(resp, req, "/login", http.StatusFound)
	} else if loggedin && !registered {
		http.Redirect(resp, req, "/registerdevice", http.StatusFound)
	} else {
		http.Redirect(resp, req, "/info", http.StatusFound)
	}
}

func LoginRoute(resp http.ResponseWriter, req *http.Request) {
	requestData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code2FA  string `json:"twofacode"`
	}{}

	err := json.NewDecoder(req.Body).Decode(&requestData)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if requestData.Code2FA == "" {
		err = pushover.Client.Login(requestData.Email, requestData.Password)
	} else {
		err = pushover.Client.Login2FA(requestData.Email, requestData.Password, requestData.Code2FA)
	}

	if _, is2faerror := err.(*poclient.Missing2FAError); is2faerror {
		resp.Write([]byte("2FA_MISSING"))
		return
	}

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	//store user info in config and return successfully
	config.Userid, config.Usersecret = pushover.Client.User()
	config.Display_Username = requestData.Email

	err = storeConfig(config, "settings.json")
	if err != nil {
		resp.Write([]byte(err.Error()))
	}

	UserinfoRoute(resp, req)
}

func RegisterRoute(resp http.ResponseWriter, req *http.Request) {
	requestData := struct {
		Devicename string `json:"devicename"`
	}{}

	err := json.NewDecoder(req.Body).Decode(&requestData)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pushover.Client.RegisterDevice(requestData.Devicename)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	pushover_retry <- true // allow the notification thread to try again

	// store user info in config and return successfully
	config.Deviceid = pushover.Client.Device()
	config.Display_Devicename = requestData.Devicename

	err = storeConfig(config, "settings.json")
	if err != nil {
		resp.Write([]byte(err.Error()))
	}

	UserinfoRoute(resp, req)
}

func LogoutRoute(resp http.ResponseWriter, req *http.Request) {
	// reset user config
	config.Userid = ""
	config.Usersecret = ""
	config.Deviceid = ""
	config.Display_Devicename = ""
	config.Display_Username = ""

	err := storeConfig(config, "settings.json")
	if err != nil {
		resp.Write([]byte(err.Error()))
	}

	resetPOClient()

	UserinfoRoute(resp, req)
}

func UserinfoRoute(resp http.ResponseWriter, req *http.Request) {
	display_username := config.Display_Username
	display_devicename := config.Display_Devicename

	var loggedin = false
	var registered = false

	if pushover != nil {
		loggedin, registered = pushover.Client.GetStatus()
	}

	response := struct {
		Username   string `json:"username"`
		Devicename string `json:"devicename"`
		LoggedIn   bool   `json:"loggedin"`
		Registered bool   `json:"registered"`
	}{
		display_username,
		display_devicename,
		loggedin,
		registered,
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(response)
}

func QuitRoute(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("ok"))

	quit_channel <- true
}

func ResetConnectionRoute(resp http.ResponseWriter, req *http.Request) {
	resetPOClient()

	resp.Write([]byte("ok"))
}
