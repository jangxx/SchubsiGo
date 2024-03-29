package main

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

/*func serveSingleFile(path string) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, path)
	}
}*/

func serveSingleFile(box *rice.Box, path string) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		f, err := box.Open(path)

		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeContent(resp, req, d.Name(), d.ModTime(), f)
	}
}
