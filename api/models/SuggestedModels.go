package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuggestedModel struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MainUid       string             `json:"mainUid" bson:"mainUid"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt" `
	SuggestedList []SugListModel     `json:"suggestedList" bson:"suggestedList"` // base64
}

type SugListModel struct {
	SugUserID  string `json:"sugUserID" bson:"sugUserID"`
	Score      int    `json:"score,omitempty" bson:"score,omitempty"`
	LoveOrHate bool   `json:"loveOrHate" bson:"loveOrHate"`
}
