package websocket

type Hub struct {
	Clients map[string]*Client
	Rooms   map[string]map[string]*Client

	Register   chan *Client
	Unregister chan *Client

	JoinRoom  chan JoinRequest
	LeaveRoom chan LeaveRequest

	Broadcast chan BroadcastRequest
}

type JoinRequest struct {
	Client *Client
	RoomID string
}

type LeaveRequest struct {
	Client *Client
	RoomID string
}

type BroadcastRequest struct {
	RoomID string
	Data   []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		JoinRoom:   make(chan JoinRequest),
		LeaveRoom:  make(chan LeaveRequest),
		Broadcast:  make(chan BroadcastRequest),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
		case client := <-h.Unregister:
			delete(h.Clients, client.ID)

			for roomID := range client.Rooms {
				h.removeFromRoom(client, roomID)
			}

			close(client.Send)

		case req := <-h.JoinRoom:
			h.addToRoom(req.Client, req.RoomID)

		case req := <-h.LeaveRoom:
			h.removeFromRoom(req.Client, req.RoomID)

		case msg := <-h.Broadcast:
			h.broadcastToRoom(msg.RoomID, msg.Data)
		}
	}
}

func (h *Hub) addToRoom(c *Client, roomID string) {
	if _, ok := h.Rooms[roomID]; !ok {
		h.Rooms[roomID] = make(map[string]*Client)
	}

	h.Rooms[roomID][c.ID] = c
	c.Rooms[roomID] = true
}

func (h *Hub) removeFromRoom(c *Client, roomID string) {
	if room, ok := h.Rooms[roomID]; ok {
		delete(room, c.ID)
	}

	delete(c.Rooms, roomID)
}

func (h *Hub) broadcastToRoom(roomID string, data []byte) {
	room, ok := h.Rooms[roomID]
	if !ok {
		return
	}

	for _, client := range room {
		select {
		case client.Send <- data:
		default:
			// slow client → drop
			close(client.Send)
			delete(h.Clients, client.ID)
		}
	}
}
