package chat

import "log"

type Room struct {
	register chan *Client

	unregister chan *Client

	clients map[*Client]bool

	broadcast chan Message
}

func NewRoom() *Room {
	return &Room{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
	}
}
func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			log.Println(client.name + " join.")
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)

				log.Println(client.name + " left.")
			}
		case msg := <-r.broadcast:
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
