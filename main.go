package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/hoai1107/chatserver/chat"
	"github.com/hoai1107/chatserver/logwrapper"
)

func main() {
	logwrapper.InitLogger()

	var addr = flag.String("addr", ":8080", "The address of the app")
	flag.Parse()

	// Serve static file
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	room := chat.NewRoom()
	go room.Run()

	http.Handle("/", &templateHandler{filename: "home.html"})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(room, w, r)
	})

	logwrapper.Info(fmt.Sprint("Server is running on ", *addr))

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		logwrapper.Error(err)
	}
}
