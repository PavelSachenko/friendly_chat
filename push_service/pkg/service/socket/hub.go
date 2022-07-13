package socket

import (
	"fmt"
	"github.com/pavel/push_service/pkg/utils"
)

//type Hub interface {
//	Run()
//}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan Broadcast

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	//InitHub
}

type Broadcast struct {
	Broadcast []byte
	Username  string
	UserIds   []uint64
}

func NewHub(broadcast chan Broadcast) *Hub {
	return &Hub{
		Broadcast:  broadcast,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			fmt.Println("Message.TEXT from socket: " + string(message.Broadcast))
			fmt.Println("Message.USERNAME from socket: " + string(message.Username))

			if message.UserIds == nil {
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			} else {
				for client := range h.clients {
					if utils.ContainsUint(message.UserIds, client.userID) {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}

			//if message.Username == "" {
			//	for client := range h.clients {
			//		select {
			//		case client.send <- message:
			//		default:
			//			close(client.send)
			//			delete(h.clients, client)
			//		}
			//	}
			//} else {
			//	for client := range h.clients {
			//		fmt.Println("h.broadcast.username: " + message.Username)
			//		fmt.Println("client name: " + client.username)
			//		if client.username != "" && message.Username == client.username {
			//			select {
			//			case client.send <- message:
			//			default:
			//				close(client.send)
			//				delete(h.clients, client)
			//			}
			//		}
			//	}
			//}
		}
	}
}
