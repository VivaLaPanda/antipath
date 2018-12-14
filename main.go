package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/VivaLaPanda/antipath/client"
	"github.com/VivaLaPanda/antipath/engine"
)

var apiPort = flag.String("apiPort", "localhost:9095", "Which port to serve the API on")

func main() {
	flag.Parse()

	// http.HandleFunc("/", serveHome)
	engine := engine.NewEngine(100, 40)
	for idx := 0; idx < 30; idx++ {
		engine.AddPlayer()
	}
	http.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWs(engine, w, r)
	})
	http.HandleFunc("/windowsize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(strconv.Itoa(engine.WindowSize)))
		return
	})
	err := http.ListenAndServe(*apiPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
