package chat

import (
	"fmt"
	"github.com/hoai1107/chatserver/logwrapper"
)

type Room struct {
	name string

	register chan *Client

	unregister chan *Client

	clients map[*Client]bool

	broadcast chan Message

	history []Message
}

func NewRoom(roomName string) *Room {
	return &Room{
		name:       roomName,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message, 5),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			msg := fmt.Sprintf("%s joins room %s", client.name, r.name)
			logwrapper.Info(msg)
			r.broadcast <- NewMessage(NotifyType, client.name, msg)

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)

				msg := fmt.Sprintf("%s lefts room %s", client.name, r.name)
				logwrapper.Info(msg)
				r.broadcast <- NewMessage(NotifyType, client.name, msg)
			}

		case msg := <-r.broadcast:
			r.history = append(r.history, msg)
			for c := range r.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(r.clients, c)
				}
			}
		}

	}
}
