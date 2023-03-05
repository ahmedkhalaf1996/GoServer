package chat

import (
	"encoding/json"
)

// Message represents a chat message
type Message struct {
	Message     string   `json:"message" bson:"message" "`
	MessageType []string `json:"messageType" bson:"messageType""`
	Sender      string   `json:"sender" bson:"sender""`
	Received    string   `json:"received" bson:"received""`
}

// FromJSON created a new Message struct from given JSON
func FromJSON(jsonInput []byte) (message *Message) {
	json.Unmarshal(jsonInput, &message)
	return
}
