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
