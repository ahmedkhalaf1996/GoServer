package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockingModel struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MainUid     string             `json:"mainUid" bson:"mainUid"`
	BlockedList []BlockedListModel `json:"blockedList" bson:"blockedList"` // base64
}

type BlockedListModel struct {
	BlockedUserID string `json:"blockedUserID" bson:"blockedUserID"`
}
