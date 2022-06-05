package chat

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	conn *websocket.Conn

	send chan Message

	room *Room

	name string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) SetName(newName string) {
	log.Println(c.name + " change username to " + newName)
	c.name = newName
}

func (c *Client) read() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)

		log.Println(msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "Send":
			c.room.broadcast <- msg
		case "ChangeName":
			c.SetName(msg.Username)
		}

	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteJSON(msg)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func RandomUsername() string {
	var name strings.Builder
	name.WriteString("User ")
	name.WriteString(strconv.Itoa(rand.Intn(1000)))

	return name.String()
}

func ServeWs(room *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	params := r.URL.Query()
	username := params.Get("username")

	client := &Client{
		room: room,
		conn: conn,
		send: make(chan Message, 10),
		name: username,
	}

	client.room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
