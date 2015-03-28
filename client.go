package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client: A Chat Client
type client struct {
	// client websocket connection
	socket *websocket.Conn
	// client message channel
	send chan *message
	// client current chatroom
	room *room
	// client's userData
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.SentTime = time.Now()
			msg.Sender = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
