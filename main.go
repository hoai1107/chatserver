package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoai1107/chatserver/chat"
	"github.com/hoai1107/chatserver/logwrapper"
)

var hub *chat.Hub = chat.NewHub()

func main() {
	logwrapper.InitLogger()

	var addr = flag.String("addr", ":8080", "The address of the app")
	flag.Parse()

	r := mux.NewRouter()

	//Serve static file
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//Routes
	r.Handle("/", &templateHandler{filename: "index.html"})
	r.Handle("/chat", &templateHandler{filename: "room.html"})
	r.Handle("/room_list", &templateHandler{filename: "room_list.html"})

	r.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		chat.AddRoomHandler(w, r, hub)
	})
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})
	r.HandleFunc("/all_rooms", func(w http.ResponseWriter, r *http.Request) {
		chat.GetAllRoom(w, r, hub)
	})

	//Start server
	logwrapper.Info(fmt.Sprint("Server is running on ", *addr))
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		logwrapper.Error(err)
	}
}
