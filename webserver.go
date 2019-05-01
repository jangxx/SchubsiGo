package main

import (
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"

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

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/login", LoginRoute).Methods("POST")
	apiRouter.HandleFunc("/register", RegisterRoute).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    config.Addr + ":" + config.Port,
	}

	log.Printf("Server is listening on %s:%s\n", config.Addr, config.Port)

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
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	email := req.PostFormValue("email")
	password := req.PostFormValue("password")

	err := pushover.Login(email, password)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	} else {
		//store user info in config and return successfully
		config.Userid, config.Usersecret = pushover.User()

		err := storeConfig(config, "settings.json")
		if err != nil {
			resp.Write([]byte(err.Error()))
		}
	}
}

func RegisterRoute(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	devicename := req.PostFormValue("devicename")

	err := pushover.RegisterDevice(devicename)

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	} else {
		//store user info in config and return successfully
		config.Deviceid = pushover.Device()

		err := storeConfig(config, "settings.json")
		if err != nil {
			resp.Write([]byte(err.Error()))
		}
	}
}
