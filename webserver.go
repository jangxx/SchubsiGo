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

	staticBox := rice.MustFindBox("./webinterface/static/")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox())))
	buildBox := rice.MustFindBox("./webinterface/build/")
	r.PathPrefix("/build/").Handler(http.StripPrefix("/build/", http.FileServer(buildBox.HTTPBox())))

	r.HandleFunc("/", IndexRoute)

	r.HandleFunc("/login", serveSingleFile(buildBox, "html/login.html")).Methods("GET")
	r.HandleFunc("/registerdevice", serveSingleFile(buildBox, "html/registerdevice.html")).Methods("GET")
	r.HandleFunc("/done", serveSingleFile(buildBox, "html/done.html")).Methods("GET")
	r.HandleFunc("/quit-app", serveSingleFile(buildBox, "html/quit.html")).Methods("GET")

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/login", LoginRoute).Methods("POST")
	apiRouter.HandleFunc("/register", RegisterRoute).Methods("POST")
	apiRouter.HandleFunc("/logout", LogoutRoute).Methods("POST")
	apiRouter.HandleFunc("/userinfo", UserinfoRoute).Methods("GET")
	apiRouter.HandleFunc("/quit", QuitRoute).Methods("POST")

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
	if loggedin, registered := pushover.GetStatus(); !loggedin {
		http.Redirect(resp, req, "/login", http.StatusFound)
	} else if loggedin && !registered {
		http.Redirect(resp, req, "/registerdevice", http.StatusFound)
	} else {
		http.Redirect(resp, req, "/done", http.StatusFound)
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
		err = pushover.Login(requestData.Email, requestData.Password)
	} else {
		err = pushover.Login2FA(requestData.Email, requestData.Password, requestData.Code2FA)
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
	config.Userid, config.Usersecret = pushover.User()
	config.Display_Username = requestData.Email

	err = storeConfig(config, "settings.json")
	if err != nil {
		resp.Write([]byte(err.Error()))
	}

	resp.Write([]byte("SUCCESS"))
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

	err = pushover.RegisterDevice(requestData.Devicename)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	pushover_retry <- true // allow the notification thread to try again

	// store user info in config and return successfully
	config.Deviceid = pushover.Device()
	config.Display_Devicename = requestData.Devicename

	err = storeConfig(config, "settings.json")
	if err != nil {
		resp.Write([]byte(err.Error()))
	}
}

func LogoutRoute(resp http.ResponseWriter, req *http.Request) {
	resetPOClient()

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
}

func UserinfoRoute(resp http.ResponseWriter, req *http.Request) {
	display_username := config.Display_Username
	display_devicename := config.Display_Devicename

	loggedin, registered := pushover.GetStatus()

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
