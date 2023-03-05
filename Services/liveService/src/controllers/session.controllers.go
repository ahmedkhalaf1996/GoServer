package controllers

import (
	"main/Services/liveService/src/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// x
// CreateSession - Creates user session
func CreateSession(ctx *gin.Context) {
	db := ctx.MustGet("db").(*mongo.Client)
	collection := db.Database("test").Collection("sessions")

	var session interfaces.Session
	if err := ctx.ShouldBindJSON(&session); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, _ := collection.InsertOne(ctx, session)
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	url := CreateSocket(session, ctx, insertedID)
	ctx.JSON(http.StatusOK, gin.H{"socket": url})
}
