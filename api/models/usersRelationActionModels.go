package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserChatRequestes struct {
	RequestID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderUserID  string             `json:"senderUserID" bson:"senderUserID"`
	ReceverUserID string             `json:"receverUserID" bson:"receverUserID"`
	IsAcceptedYet bool               `json:"isAcceptedYet" bson:"isAcceptedYet"`
	SendedAt      time.Time          `json:"sendedAt" bson:"sendedAt" `
}

type UserLoveRequestes struct {
	RequestID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderUserID  string             `json:"senderUserID" bson:"senderUserID"`
	ReceverUserID string             `json:"receverUserID" bson:"receverUserID"`
	IsAcceptedYet bool               `json:"isAcceptedYet" bson:"isAcceptedYet"`
	SendedAt      time.Time          `json:"sendedAt" bson:"sendedAt" `
}

type UserBuzzRequestes struct {
	RequestID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderUserID  string             `json:"senderUserID" bson:"senderUserID"`
	ReceverUserID string             `json:"receverUserID" bson:"receverUserID"`
	IsAcceptedYet bool               `json:"isAcceptedYet" bson:"isAcceptedYet"`
	SendedAt      time.Time          `json:"sendedAt" bson:"sendedAt" `
}

type UserStarRequestes struct {
	RequestID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderUserID  string             `json:"senderUserID" bson:"senderUserID"`
	ReceverUserID string             `json:"receverUserID" bson:"receverUserID"`
	IsAcceptedYet bool               `json:"isAcceptedYet" bson:"isAcceptedYet"`
	SendedAt      time.Time          `json:"sendedAt" bson:"sendedAt" `
}

// type UserChatRequestes struct {
// 	RequestID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	SenderUserID  string             `json:"senderUserID" bson:"senderUserID"`
// 	ReceverUserID string             `json:"receverUserID" bson:"receverUserID"`
// 	IsAcceptedYet bool               `json:"isAcceptedYet" bson:"isAcceptedYet"`
// 	IsChatRequest bool               `json:"isChatRequest" bson:"isChatRequest"`
// 	IsLoveRequest bool               `json:"isLoveRequest" bson:"isLoveRequest"`
// 	IsBuzzRequest bool               `json:"isBuzzRequest" bson:"isBuzzRequest"`
// 	IsStarRequest bool               `json:"isStarRequest" bson:"isStarRequest"`
// 	SendedAt      time.Time          `json:"sendedAt" bson:"sendedAt" `
// }
