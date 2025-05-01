package handlers

import (
	"fmt"
	"net/http"

	"github.com/ALOK9442/web-video-chat/backend/core/models"
	"github.com/ALOK9442/web-video-chat/backend/core/hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// func HandleWebsocket(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		fmt.Println("err occured:", err)
// 		return
// 	}
// 	user := 
// }

func readMessage(user *models.User){
	defer func()  {
		hub.HubInstance.HandleUnRegister(user)
		user.Conn.Close()
	}()
	msgType, msg, err := websocket
}
