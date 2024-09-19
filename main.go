package main

import (
  "encoding/json"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

type Message struct {
  ID      string `json:"id"`
  Message string `json:"message"`
}

func pushMessage(w http.ResponseWriter, r *http.Request) {
  log.Println("received request to push")
  if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
  }
  
  var message Message

  err := json.NewDecoder(r.Body).Decode(&message)
  if err != nil {
    log.Println("error decoding json", err)
    return
  }

  channel, ok := getChannel(message.ID)

  if !ok {
    log.Println("no channel")
    return
  }

  channel.broadcast <-[]byte(message.Message)
}

func main() {
	flag.Parse()

  lobby := newPrivateChannel()
  go lobby.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    log.Println("serving websocket")
		serveWs(lobby, w, r)
	})
  http.HandleFunc("/push-message", pushMessage)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
