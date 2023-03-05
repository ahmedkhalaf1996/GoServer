package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VeryfiyModel struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         string             `json:"userID" bson:"userID"`
	IsVeryfiyedYet bool               `json:"isVeryfiyedYet" bson:"isVeryfiyedYet"`
	IsSendedToMail bool               `json:"isSendedToMail" bson:"isSendedToMail"`
	TryNumber      int                `json:"tryNumber" bson:"tryNumber"`
	LastUpdated    time.Time          `json:"lastUpdated" bson:"lastUpdated"`
	VeryfiyCode    string             `json:"veryfiyCode" bson:"veryfiyCode"`
}

type AddVeryfiyEmailModel struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         string             `json:"userID" bson:"userID"`
	IsVeryfiyedYet bool               `json:"isVeryfiyedYet" bson:"isVeryfiyedYet"`
	LastUpdated    time.Time          `json:"lastUpdated" bson:"lastUpdated"`
	TryNumber      int                `json:"tryNumber" bson:"tryNumber"`
	ProvidedEmail  string             `json:"providedEmail" bson:"providedEmail"`
	VeryfiyCode    string             `json:"veryfiyCode" bson:"veryfiyCode"`
}

type AddVeryfiyPhoneNumberModel struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID              string             `json:"userID" bson:"userID"`
	IsVeryfiyedYet      bool               `json:"isVeryfiyedYet" bson:"isVeryfiyedYet"`
	LastUpdated         time.Time          `json:"lastUpdated" bson:"lastUpdated"`
	TryNumber           int                `json:"tryNumber" bson:"tryNumber"`
	ProvidedPhoneNumber string             `json:"providedPhoneNumber" bson:"providedPhoneNumber"`
	VeryfiyCode         string             `json:"veryfiyCode" bson:"veryfiyCode"`
}
