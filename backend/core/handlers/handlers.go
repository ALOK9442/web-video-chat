package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ALOK9442/web-video-chat/backend/core/models"
	"github.com/ALOK9442/web-video-chat/backend/core/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("err occured:", err)
		return
	}
	user := &models.User{
		Id:   uuid.NewString(),
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.HubInstance.Register <- user
	go readMessage(user)
	go writeMessage(user)
}

func readMessage(user *models.User) {
	defer func() {
		hub.HubInstance.UnRegister <- user
		user.Conn.Close()
	}()
	for {
		msgType, msg, err := user.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error occured: %v of type: %v \n", err, msgType)
			user.Send <- []byte(`{"type":"unknown","message":"Error occured send again and if the issue persists refresh page"}`)
			break
		}
		var signal models.SignalMessage
		if error := json.Unmarshal(msg, &signal); error != nil {
			fmt.Println("error occured", error)
			continue
		}
		switch signal.Type {
		case "offer", "answer", "ice-candidate", "chat":
			hub.HubInstance.Broadcast <- &models.BroadcastMessage{
				Client:  user,
				Message: msg,
			}
		}
	}
}

func writeMessage(user *models.User) {
	for msg := range user.Send {
		err := user.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("error occured:", err)
			break
		}
	}
}
