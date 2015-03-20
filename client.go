package main

import (
	"github.com/gorilla/websocket"
)

// client: A Chat Client
type client struct {
	// client websocket connection
	socket *websocket.Conn
	// client message channel
	send chan []byte
	// client current chatroom
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
