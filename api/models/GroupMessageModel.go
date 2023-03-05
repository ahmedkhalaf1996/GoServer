package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageGroupModel struct {
	ConversationID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Creator         string             `json:"creator" bson:"creator"`
	UsersJoinedList []JoinedUsers      `json:"usersJoinedList" bson:"usersJoinedList"`
	GroupMessages   []PayloadGroupList `json:"groupMessages" bson:"groupMessages"` // base64
}

type JoinedUsers struct {
	JoinedUserID string `json:"joinedUserID" bson:"joinedUserID"`
	IsAdmin      bool   `json:"isAdmin" bson:"isAdmin"`
}

type PayloadGroupList struct {
	Message     string    `json:"message" bson:"message"   validate:"required,min=1"`
	MessageType []string  `json:"messageType" bson:"messageType"  validate:"required"`
	Sender      string    `json:"sender" bson:"sender"  validate:"required"`
	SendedAt    time.Time `json:"sendedAt" bson:"sendedAt" `
}
