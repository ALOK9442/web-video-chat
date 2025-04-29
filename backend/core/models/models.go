package models

import "github.com/gorilla/websocket"

type User struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Room struct {
	User1 *User
	User2 *User
}

type BroadcastMessage struct {
	Client *User
	Message []byte
}

func NewUser(c *websocket.Conn) *User {
	return &User{
		Conn: c,
		Send: make(chan []byte),
	}
}
