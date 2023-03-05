// package database

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var Client *mongo.Client
// var DB *mongo.Database

// func Connect() {
// 	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

// 	Client, err := mongo.Connect(ctx, options.Client().ApplyURI("http://localhost:27017/"))

// 	if err == nil {
// 		fmt.Println("connection to DB..")
// 		DB = Client.Database("gotest")
// 	}

// }

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	Client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err == nil {
		fmt.Println("connection to DB..")
		DB = Client.Database("test")
	}

}
