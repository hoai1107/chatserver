package chat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hoai1107/chatserver/logwrapper"
)

func AddRoomHandler(w http.ResponseWriter, r *http.Request, hub *Hub) {
	payload := make(map[string]interface{})

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logwrapper.Error(err)
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		logwrapper.Error(err)
	}
	room := NewRoom(payload["name"].(string))
	logwrapper.Info("New Room: " + payload["name"].(string))

	hub.RegisterRoom(room)
	go room.Run()

	w.WriteHeader(http.StatusOK)
}

func GetAllRoom(w http.ResponseWriter, r *http.Request, hub *Hub) {
	keys := make([]string, len(hub.Rooms))

	i := 0
	for k := range hub.Rooms {
		keys[i] = k
		i++
	}

	resp := make(map[string]interface{})
	resp["rooms"] = keys
	json.NewEncoder(w).Encode(resp)
}
