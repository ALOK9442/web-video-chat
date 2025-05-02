package models

import "github.com/gorilla/websocket"

type User struct {
	Id   string
	Conn *websocket.Conn
	Send chan []byte
}

type Room struct {
	Id    string
	User1 *User
	User2 *User
}

type BroadcastMessage struct {
	Client  *User
	Message []byte
}

type SignalMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewUser(c *websocket.Conn) *User {
	return &User{
		Conn: c,
		Send: make(chan []byte),
	}
}
