package notification

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a chat message
type Message struct {
	NotifyID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId              string             `json:"userId" bson:"userId"`
	NumOfUnReadedNotify int                `json:"numOfUnReadedNotify" bson:"numOfUnReadedNotify"`
	NotifyList          []NotifyDeatils    `json:"notifyList" bson:"notifyList"` // base64

}
type NotifyDeatils struct {
	NotifyNumber   int         `json:"notifyNumber" bson:"notifyNumber"`
	ReadedOrNot    bool        `json:"readedOrNot" bson:"readedOrNot"`
	SendedAt       time.Time   `json:"sendedAt" bson:"sendedAt" `
	NotifyMessage  string      `json:"notifyMessage" bson:"notifyMessage" `
	NotifyTypeData interface{} `json:"notifyTypeData" bson:"notifyTypeData"` // base64
}

// FromJSON created a new Message struct from given JSON
func FromJSON(jsonInput []byte) (message *Message) {
	json.Unmarshal(jsonInput, &message)
	return
}
