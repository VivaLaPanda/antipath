package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/VivaLaPanda/antipath/client"
	"github.com/VivaLaPanda/antipath/engine"
)

var apiPort = flag.String("apiPort", "localhost:9095", "Which port to serve the API on")

func main() {
	flag.Parse()

	// http.HandleFunc("/", serveHome)
	engine := engine.NewEngine(100)
	http.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWs(engine, w, r)
	})
	err := http.ListenAndServe(*apiPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
