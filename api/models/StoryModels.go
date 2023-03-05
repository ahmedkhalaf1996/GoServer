package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryMainModel struct {
	StoryID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId       string             `json:"userId" bson:"userId"`
	NumOfStoryes int                `json:"numOfStoryes" bson:"numOfStoryes"`
	StoryPayload []StoryList        `json:"storyPayload" bson:"storyPayload"` // base64
}

type StoryList struct {
	StoryNumber  int       `json:"storyNumber" bson:"storyNumber"`
	StoryData    string    `json:"storyData" bson:"storyData" `
	StoryCaption string    `json:"storyCaption" bson:"storyCaption" `
	IsTypeVideo  bool      `json:"isTypeVideo" bson:"isTypeVideo"`
	SendedAt     time.Time `json:"sendedAt" bson:"sendedAt" `
}
