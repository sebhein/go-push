package main

// PrivateChannel maintains the set of active clients and broadcasts messages to the
// clients.
type PrivateChannel struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newPrivateChannel() *PrivateChannel {
	return &PrivateChannel{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (pc *PrivateChannel) run() {
	for {
		select {
		case client := <-pc.register:
			pc.clients[client] = true
		case client := <-pc.unregister:
			if _, ok := pc.clients[client]; ok {
				delete(pc.clients, client)
				//close(client.send)
			}
		case message := <-pc.broadcast:
			for client := range pc.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(pc.clients, client)
				}
			}
		}
	}
}
