package main

import (
	"fmt"
	"log"
	"net/http"
)

type MainHandler struct {
	Settings Settings
}

func (mh *MainHandler) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[knative-labs] Hi from %s (rev: %s)", mh.Settings.Service, mh.Settings.Revision)
}

func main() {
	settings := NewSettings()
	mainHandler := &MainHandler{Settings: settings}

	http.HandleFunc("/", mainHandler.handler)
	log.Fatal(http.ListenAndServe(":"+settings.Port, nil))
}
