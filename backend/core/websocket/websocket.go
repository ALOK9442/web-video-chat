package websocket

import (
	"sync"

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

func HubInstance() *Hub {
	return &Hub{
		// Rooms:        make([]*models.Room, 0),
		WaitingQueue: make([]*models.User, 0),
		UserToRoom:   make(map[*models.User]*models.Room),
		Register:     make(chan *models.User),
		UnRegister:   make(chan *models.User),
		Skip:         make(chan *models.User),
		Broadcast:    make(chan *models.BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.handleRegister(client)
		case client := <-h.UnRegister:
			h.handleUnRegister(client)
		case client := <-h.Skip:
			h.handleSkip(client)
		case broadCastMessage := <-h.Broadcast:
			h.handleBroadcast(broadCastMessage.Client, broadCastMessage.Message)
		}
	}
}

func (h *Hub) handleRegister(C *models.User) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.WaitingQueue) > 0 {
		partner := h.WaitingQueue[0]
		h.WaitingQueue = h.WaitingQueue[1:]
		room := &models.Room{
			User1: partner,
			User2: C,
		}
		// h.Rooms = append(h.Rooms, room)

		h.UserToRoom[partner] = room
		h.UserToRoom[C] = room
		partner.Send <- []byte(`{"type":"system","message": "You're now Connected" }`)
		C.Send <- []byte(`{"type":"system","message": "You're now Connected" }`)

	} else {
		h.WaitingQueue = append(h.WaitingQueue, C)
		C.Send <- []byte(`{"type":"system", "message":"Waiting for a partner..."}`)
	}
}

func (h *Hub) handleUnRegister(C *models.User) {
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
			h.handleRegister(partner)
		}
		return
	}

}

func (h *Hub) handleSkip(C *models.User) {

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
				h.handleRegister(partner)
			}
		}
	}
}

func (h *Hub) handleBroadcast(C *models.User, message []byte) {
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
