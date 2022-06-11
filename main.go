package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoai1107/chatserver/chat"
	"github.com/hoai1107/chatserver/logwrapper"
)

var hub chat.Hub = *chat.NewHub()

func main() {
	logwrapper.InitLogger()

	var addr = flag.String("addr", ":8080", "The address of the app")
	flag.Parse()

	room := chat.NewRoom("Chat room 1")
	go room.Run()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Handle("/", &templateHandler{filename: "index.html"})
	r.HandleFunc("/create", addRoomHandler)
	r.Handle("/chat", &templateHandler{filename: "room.html"})
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(room, w, r)
	})

	logwrapper.Info(fmt.Sprint("Server is running on ", *addr))

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		logwrapper.Error(err)
	}
}

func addRoomHandler(w http.ResponseWriter, r *http.Request) {
	payload := make(map[string]interface{})

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logwrapper.Error(err)
	}

	json.Unmarshal(body, &payload)
	room := chat.NewRoom(payload["name"].(string))
	logwrapper.Info("New Room: " + payload["name"].(string))

	hub.RegisterRoom(room)
	go room.Run()
}
