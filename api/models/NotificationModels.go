package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationModel struct {
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

// data --------------
type ActivateAccount struct {
	ActivateAccount bool `json:"activateAccount" bson:"activateAccount"`
}

type UnReeadedMessages struct {
	UnReeadedMessages bool   `json:"unReeadedMessages" bson:"unReeadedMessages"`
	ConversationID    string `json:"conversationID" bson:"conversationID"`
}

type LoveSended struct {
	LoveSended       bool   `json:"loveSended" bson:"loveSended"`
	SendedFromUserID string `json:"sendedFromUserID" bson:"sendedFromUserID"`
	RequestID        string `json:"requestID" bson:"requestID"`
}

type RequestToChat struct {
	RequestToChat    bool   `json:"requestToChat" bson:"requestToChat"`
	SendedFromUserID string `json:"sendedFromUserID" bson:"sendedFromUserID"`
	RequestID        string `json:"requestID" bson:"requestID"`
}

type StarSended struct {
	StarSended       bool   `json:"starSended" bson:"requestToChat"`
	SendedFromUserID string `json:"sendedFromUserID" bson:"sendedFromUserID"`
	RequestID        string `json:"requestID" bson:"requestID"`
}

type BuzzSended struct {
	BuzzSended       bool   `json:"buzzSended" bson:"buzzSended"`
	SendedFromUserID string `json:"sendedFromUserID" bson:"sendedFromUserID"`
	RequestID        string `json:"requestID" bson:"requestID"`
}

type AddedtoChatRoom struct {
	AddedtoChatRoom bool   `json:"addedtoChatRoom" bson:"addedtoChatRoom"`
	ConversationID  string `json:"conversationID" bson:"conversationID"`
	UserAddedID     string `json:"userAddedID" bson:"userAddedID"`
}

type InviteUserToEvent struct {
	InviteUserToEvent bool   `json:"inviteUserToEvent" bson:"inviteUserToEvent"`
	UserAddedID       string `json:"userAddedID" bson:"userAddedID"`
	EventID           string `json:"eventID" bson:"eventID"`
}
