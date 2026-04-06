package chat

import "encoding/json"

type MessageType string

const (
	JoinRoom    MessageType = "join_room"
	LeaveRoom   MessageType = "leave_room"
	MessageText MessageType = "message_text"
)

type Message struct {
	Type    MessageType     `json:"type"`
	RoomID  string          `json:"room_id"`
	Payload json.RawMessage `json:"payload"`
}

type TextChatPayload struct {
	Text string `json:"text"`
}
