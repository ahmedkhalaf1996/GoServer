package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

// HexID struct to get id from DB x
type HexID struct {
	ID primitive.ObjectID `bson:"_id"`
}
