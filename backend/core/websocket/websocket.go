package hub

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ALOK9442/web-video-chat/backend/core/helpers"
	"github.com/ALOK9442/web-video-chat/backend/core/models"
)

// var mu sync.Mutex

type Hub struct {
	// Rooms        []*models.Room
	WaitingQueue []*models.User
	UserToRoom   map[*models.User]*models.Room
	Register     chan *models.User
	UnRegister   chan *models.User
	Skip         chan *models.User
	Broadcast    chan *models.BroadcastMessage
	mu           sync.Mutex
}

var HubInstance = &Hub{
	// return &Hub{
	// Rooms:        make([]*models.Room, 0),
	WaitingQueue: make([]*models.User, 0),
	UserToRoom:   make(map[*models.User]*models.Room),
	Register:     make(chan *models.User),
	UnRegister:   make(chan *models.User),
	Skip:         make(chan *models.User),
	Broadcast:    make(chan *models.BroadcastMessage),
	// }
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.HandleRegister(client)
		case client := <-h.UnRegister:
			h.HandleUnRegister(client)
		case client := <-h.Skip:
			h.HandleSkip(client)
		case broadCastMessage := <-h.Broadcast:
			h.HandleBroadcast(broadCastMessage.Client, broadCastMessage.Message)
		}
	}
}

func (h *Hub) HandleRegister(C *models.User) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.WaitingQueue) > 0 {
		partner := h.WaitingQueue[0]
		h.WaitingQueue = h.WaitingQueue[1:]

		room := &models.Room{
			Id:    partner.Id + " " + C.Id,
			User1: partner,
			User2: C,
		}
		// h.Rooms = append(h.Rooms, room)
		innerMessageClient := map[string]interface{}{
			"roomId":  room.Id,
			"message": "You're now Connected",
			"role":"calleee",
		}
		innerMessagePartner := map[string]interface{}{
			"roomId":  room.Id,
			"message": "You're now Connected",
			"role":"caller",
		}
		// innerMessageJSON, err := json.Marshal(innerMessage)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// outerMessageClient := map[string]interface{}{
		// 	"type": "success",
		// 	"data": innerMessageClient,
		// }
		// outerMessagePeer := map[string]interface{}{
		// 	"type": "success",
		// 	"data": innerMessagePeer,
		// }
		// finalJSON, err := json.Marshal(outerMessage)
		// if err != nil {
		// 	fmt.Println("Error marshalling:", err)
		// 	return
		// }

		h.UserToRoom[partner] = room
		h.UserToRoom[C] = room

		partner.Send <- helpers.MarshalMessage("caller", innerMessagePartner)
		C.Send <- helpers.MarshalMessage("calleee", innerMessageClient)

	} else {
		h.WaitingQueue = append(h.WaitingQueue, C)
		innerMessage := map[string]interface{}{
			"type":    "system",
			"message": "Waiting for a partner...",
		}
		// innerMessageJSON, err := json.Marshal(innerMessage)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		outerMessage := map[string]interface{}{
			"type": "waiting",
			"data": innerMessage,
		}

		finalJSON, err := json.Marshal(outerMessage)
		if err != nil {
			fmt.Println("Error marshalling:", err)
			return
		}

		C.Send <- finalJSON
	}
}

func (h *Hub) HandleUnRegister(C *models.User) {
	for i, user := range h.WaitingQueue {
		if C == user {
			h.WaitingQueue = append(h.WaitingQueue[:i], h.WaitingQueue[i+1:]...)
			return
		}
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, exists := h.UserToRoom[C]; exists {
		var partner *models.User
		if room.User1 == C {
			partner = room.User2
		} else if room.User2 == C {
			partner = room.User1
		}
		delete(h.UserToRoom, partner)
		delete(h.UserToRoom, C)
		if partner != nil {
			partner.Send <- []byte(`{"type":"system","message":"Your Partner has been disconnected, waiting for another partner..."}`)
			h.HandleRegister(partner)
		}
		return
	}

}

func (h *Hub) HandleSkip(C *models.User) {

	if room, exists := h.UserToRoom[C]; exists {
		if room.User1 == C || room.User2 == C {
			var partner *models.User
			if room.User1 == C {
				partner = room.User2
			} else if room.User2 == C {
				partner = room.User1
			}
			h.mu.Lock()
			defer h.mu.Unlock()
			delete(h.UserToRoom, partner)
			delete(h.UserToRoom, C)
			h.WaitingQueue = append(h.WaitingQueue, C)
			C.Send <- []byte(`{"type":"system","message":"You skipped. Searching again..."}`)
			if partner != nil {
				partner.Send <- []byte(`{"type":"system","message":"Your partner skipped. Searching again..."}`)
				h.HandleRegister(partner)
			}
		}
	}
}

func (h *Hub) HandleBroadcast(C *models.User, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, exists := h.UserToRoom[C]; exists {
		if C == room.User1 {
			select {
			case room.User2.Send <- message:
			default:
				close(room.User2.Send)
				delete(h.UserToRoom, room.User1)
				delete(h.UserToRoom, room.User2)
			}
		} else if C == room.User2 {
			select {
			case room.User1.Send <- message:
			default:
				close(room.User1.Send)
				delete(h.UserToRoom, room.User1)
				delete(h.UserToRoom, room.User2)
			}
		}
	}
}
