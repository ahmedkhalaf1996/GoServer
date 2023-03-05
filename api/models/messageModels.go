package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessagePrivateRoomModels struct {
	ConversationID            primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstUserID               string                      `json:"firstUserID" bson:"firstUserID"`
	SecondUserID              string                      `json:"secondUserID" bson:"secondUserID"`
	StartedAt                 time.Time                   `json:"startedAt" bson:"startedAt" `
	FirstUserAllowedOrBloced  bool                        `json:"firstUserAllowedOrBloced" bson:"firstUserAllowedOrBloced"`
	SecondUserAllowedOrBloced bool                        `json:"secondUserAllowedOrBloced" bson:"secondUserAllowedOrBloced"`
	MessagePayloadFirstCopy   []PayloadPrivateDeatilsList `json:"messagePayloadFirstCopy" bson:"messagePayloadFirstCopy"` // base64
	MessagePayloadSecodCopy   []PayloadPrivateDeatilsList `json:"messagePayloadSecodCopy" bson:"messagePayloadSecodCopy"` // base64
}

type PayloadPrivateDeatilsList struct {
	Message     string    `json:"message" bson:"message"   validate:"required,min=1"`
	MessageType []string  `json:"messageType" bson:"messageType"  validate:"required"`
	Sender      string    `json:"sender" bson:"sender"  validate:"required"`
	Received    string    `json:"received" bson:"received"  validate:"required"`
	SendedAt    time.Time `json:"sendedAt" bson:"sendedAt" `
}
