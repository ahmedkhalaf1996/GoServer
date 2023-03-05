package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventsModels struct {
	EventID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatorID            string             `json:"creatorID" bson:"creatorID"`
	EventTitle           string             `json:"eventTitle" bson:"eventTitle" validate:"required"`
	EventDescription     string             `json:"eventDescription" bson:"eventDescription" validate:"required"`
	EventCoverImage      string             `json:"eventCoverImage" bson:"eventCoverImage"`
	EventLocation        []string           `json:"eventLocation" bson:"eventLocation"` // country city
	EventCommentsbyUsers []EventComments    `json:"eventCommentsbyUsers" bson:"eventCommentsbyUsers"`
	EventTimeing         time.Time          `json:"eventTimeing" bson:"eventTimeing"`
	EventStartedOn       string             `json:"eventStartedOn" bson:"eventStartedOn"`
	GoingUsers           []string           `json:"goingUsers" bson:"goingUsers"`
	MayBeGoingUsers      []string           `json:"mayBeGoingUsers" bson:"mayBeGoingUsers"`
	NotGoingUsers        []string           `json:"NotGoingUsers" bson:"NotGoingUsers"`
}

type EventComments struct {
	UserAddedCommentID string `json:"userAddedCommentID" bson:"userAddedCommentID"`
	CommentMessage     string `json:"commentMessage" bson:"commentMessage"`
}

// Model For update Only
type EventsUpdateModels struct {
	EventID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatorID            string             `json:"creatorID" bson:"creatorID"`
	EventTitle           string             `json:"eventTitle" bson:"eventTitle"`
	EventDescription     string             `json:"eventDescription" bson:"eventDescription"`
	EventCoverImage      string             `json:"eventCoverImage" bson:"eventCoverImage"`
	EventLocation        []string           `json:"eventLocation" bson:"eventLocation"` // country city
	EventCommentsbyUsers []EventComments    `json:"eventCommentsbyUsers" bson:"eventCommentsbyUsers"`
	EventTimeing         time.Time          `json:"eventTimeing" bson:"eventTimeing"`
	EventStartedOn       string             `json:"eventStartedOn" bson:"eventStartedOn"`
	GoingUsers           []string           `json:"goingUsers" bson:"goingUsers"`
	MayBeGoingUsers      []string           `json:"mayBeGoingUsers" bson:"mayBeGoingUsers"`
	NotGoingUsers        []string           `json:"NotGoingUsers" bson:"NotGoingUsers"`
}
