package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/exp/slices"

	// "main/Services/chat"
	"main/api/models"
	"main/database"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// var WebSocketServerUrl = "http://44.203.164.223:8080"
// var WebSocketServerUrl =

// GetMessageByNumbers
func GetMessageByNumbers(c *fiber.Ctx) error {
	var MessagePrivateRoomSchema = database.DB.Collection("private_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var room models.MessagePrivateRoomModels
	// var CopyedData models.PayloadPrivateDeatilsList

	roomID, _ := primitive.ObjectIDFromHex(c.Params("roomId"))
	userID := c.Params("userId")

	var body struct {
		From int
	}

	if reflect.TypeOf(body.From) != reflect.TypeOf(1) {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": "body From Should be an init type",
			})
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	DefaultSkipNumber := 4

	err := MessagePrivateRoomSchema.FindOne(ctx, bson.M{"_id": roomID}).Decode(&room)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	//
	// CheckUser First Or Second
	if room.FirstUserID == userID {

		endNum := 0
		if body.From <= 0 {
			endNum = len(room.MessagePayloadFirstCopy)
		} else {
			endNum = len(room.MessagePayloadFirstCopy) - (body.From * DefaultSkipNumber)
		}

		start := endNum - DefaultSkipNumber

		cals := (len(room.MessagePayloadFirstCopy) - 1) - (DefaultSkipNumber * body.From)
		// fmt.Println("clas", cals)
		if endNum >= 0 && start >= 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ms": room.MessagePayloadFirstCopy[start:endNum],
			})
		} else if endNum >= DefaultSkipNumber {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ms": "no message",
			})
		} else if cals > 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				// "ms": room.MessagePayloadFirstCopy[start:endNum],
				"ms": room.MessagePayloadFirstCopy[0:endNum],
			})
		} else if (len(room.MessagePayloadFirstCopy) - 1) <= DefaultSkipNumber {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				// "ms": room.MessagePayloadFirstCopy[start:endNum],
				"ms": room.MessagePayloadFirstCopy,
			})
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ms": "out of range",
			})
		}

	} else {
		// User Second

		endNum := 0
		if body.From <= 0 {
			endNum = len(room.MessagePayloadSecodCopy)
		} else {
			endNum = len(room.MessagePayloadSecodCopy) - (body.From * DefaultSkipNumber)
		}
		start := endNum - DefaultSkipNumber
		// reminder := (len(room.MessagePayloadSecodCopy) - 1) % DefaultSkipNumber

		if endNum >= 0 && start >= 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ms": room.MessagePayloadSecodCopy[start:endNum],
			})
		} else if endNum >= DefaultSkipNumber {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ms": "no message",
			})
		} else if (len(room.MessagePayloadSecodCopy) - 1) <= DefaultSkipNumber {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				// "ms": room.MessagePayloadSecodCopy[start:endNum],
				"ms": room.MessagePayloadSecodCopy,
			})
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				// "ms": room.MessagePayloadSecodCopy[start:endNum],
				"ms": room.MessagePayloadSecodCopy[0:endNum],
			})
		}
	}

}

// SendMessageToPrivateRoom
func SendMessageToPrivateRoom(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var MessagePrivateRoomSchema = database.DB.Collection("private_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var room models.MessagePrivateRoomModels
	var messagelist models.PayloadPrivateDeatilsList

	if err := c.BodyParser(&messagelist); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}
	if messagelist.Message == "" {
		err := "message faild Can't be empty"
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}
	//messageType
	if messagelist.MessageType == nil {
		err := "MessageType faild Can't be nil"
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	if messagelist.Sender == "" {
		err := "Sender faild Can't be empty"
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	if messagelist.Received == "" {
		err := "Recevied faild Can't be empty"
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	// chek if sender is veryfied
	Puserid, _ := primitive.ObjectIDFromHex(messagelist.Sender)
	SPuserid, _ := primitive.ObjectIDFromHex(messagelist.Received)
	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": Puserid}).Decode(&checkUser)
	Serr := UserSchema.FindOne(ctx, bson.M{"_id": SPuserid}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	if Serr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}
	//
	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}

	// check if any of the users blocked the Other User

	var BlockedSchema = database.DB.Collection("blocked_list")
	var BlockedMod models.BlockingModel
	var sBlcokModList models.BlockedListModel

	// check the sender is blocked the recever
	err = BlockedSchema.FindOne(ctx, bson.M{"mainUid": messagelist.Sender}).Decode(&BlockedMod)

	sBlcokModList.BlockedUserID = messagelist.Received

	if slices.Contains(BlockedMod.BlockedList, sBlcokModList) {
		c.Status(fiber.StatusOK)
		return c.JSON(&fiber.Map{
			"IsBlocked": "You Already Blocked this User",
		})
	}
	// check the reveced is blocked the sender

	BlockedSchema.FindOne(ctx, bson.M{"mainUid": messagelist.Received}).Decode(&BlockedMod)

	sBlcokModList.BlockedUserID = messagelist.Sender

	if slices.Contains(BlockedMod.BlockedList, sBlcokModList) {
		c.Status(fiber.StatusOK)
		return c.JSON(&fiber.Map{
			"IsBlocked": "This user is Already Blocked You , you can't send message to this user",
		})
	}

	// --------------------------------------------------

	GetTimeNow := time.Now()
	messagelist.SendedAt = GetTimeNow
	FirstnumOfDocs, _ := MessagePrivateRoomSchema.CountDocuments(ctx, bson.M{"firstUserID": messagelist.Sender, "secondUserID": messagelist.Received})
	SecondnumOfDocs, _ := MessagePrivateRoomSchema.CountDocuments(ctx, bson.M{"firstUserID": messagelist.Received, "secondUserID": messagelist.Sender})

	if FirstnumOfDocs == 0 && SecondnumOfDocs == 0 {
		// room
		room.FirstUserID = messagelist.Sender
		room.SecondUserID = messagelist.Received
		room.FirstUserAllowedOrBloced = true
		room.SecondUserAllowedOrBloced = true

		GetTimeNow := time.Now()
		room.StartedAt = GetTimeNow

		result, err := MessagePrivateRoomSchema.InsertOne(ctx, &room)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		} else {

			query := bson.M{"_id": result.InsertedID}

			MessagePrivateRoomSchema.FindOne(ctx, query).Decode(&room)

		}
	} else {
		MessagePrivateRoomSchema.FindOne(ctx, bson.M{"firstUserID": messagelist.Sender, "secondUserID": messagelist.Received}).Decode(&room)
		MessagePrivateRoomSchema.FindOne(ctx, bson.M{"firstUserID": messagelist.Received, "secondUserID": messagelist.Sender}).Decode(&room)
	}

	// add message to the room
	room.MessagePayloadFirstCopy = append(room.MessagePayloadFirstCopy, messagelist)
	room.MessagePayloadSecodCopy = append(room.MessagePayloadSecodCopy, messagelist)

	result, err := MessagePrivateRoomSchema.UpdateOne(ctx, bson.M{"_id": room.ConversationID}, bson.M{"$set": room})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}
	// err = MessagePrivateRoomSchema.UpdateOne(ctx, bson.M{"mainUid": mainUsserid},{sug}).Decode(&sug)
	if result.MatchedCount == 1 {
		err := MessagePrivateRoomSchema.FindOne(ctx, bson.M{"_id": room.ConversationID}).Decode(&room)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}

	// Call Room in Chat Service
	_, err = http.Get("http://" + string(os.Getenv("MAINHOST")) + ":8080/AddChat?ConvID=" + room.ConversationID.Hex())
	if err != nil {
		log.Fatalln(err)
	}

	c.Locals("userId", messagelist.Received)
	c.Locals("ConversationID", room.ConversationID.Hex())
	CreateSendUnReeadedMessagesPrivateRoomNotification(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// "room": room,
		"ms": messagelist,
	})

}

// RemoveUserCopyOfChat/:roomId/:userId
func RemoveUserCopyOfChat(c *fiber.Ctx) error {

	var MessagePrivateRoomSchema = database.DB.Collection("private_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var room models.MessagePrivateRoomModels

	roomID, _ := primitive.ObjectIDFromHex(c.Params("roomId"))
	userID := c.Params("userId")

	// check if room exist
	err := MessagePrivateRoomSchema.FindOne(ctx, bson.M{"_id": roomID}).Decode(&room)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "room with given id Not Exist",
		})
	}

	// checkUser First Or secod Or not any of them
	if room.FirstUserID == userID {

		var emptyarr []models.PayloadPrivateDeatilsList
		room.MessagePayloadFirstCopy = emptyarr

		_, err := MessagePrivateRoomSchema.UpdateOne(ctx, bson.M{"_id": roomID}, bson.M{"$set": room})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": "Successfully Removed the User Copy Of the Chat",
		})

	} else if room.SecondUserID == userID {

		var emptyarr []models.PayloadPrivateDeatilsList
		room.MessagePayloadSecodCopy = emptyarr

		_, err := MessagePrivateRoomSchema.UpdateOne(ctx, bson.M{"_id": roomID}, bson.M{"$set": room})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": "Successfully Removed the User Copy Of the Chat",
		})

	} else {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "user with given id Not Part Of this Private Conversation !",
		})
	}

	return nil
}

// GetPrivateRoomeData
func GetPrivateRoomeID(c *fiber.Ctx) error {
	var MessagePrivateRoomSchema = database.DB.Collection("private_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var room models.MessagePrivateRoomModels
	// var secod_user_room models.MessagePrivateRoomModels
	// var messagelist []models.PayloadPrivateDeatilsList

	MainUid := c.Params("FuserID")
	SecodUid := c.Params("SuserID")

	FirstnumOfDocs, _ := MessagePrivateRoomSchema.CountDocuments(ctx, bson.M{"firstUserID": MainUid, "secondUserID": SecodUid})
	SecondnumOfDocs, _ := MessagePrivateRoomSchema.CountDocuments(ctx, bson.M{"firstUserID": SecodUid, "secondUserID": MainUid})

	if FirstnumOfDocs == 0 && SecondnumOfDocs == 0 {
		// room
		room.FirstUserID = MainUid
		room.SecondUserID = SecodUid
		room.FirstUserAllowedOrBloced = true
		room.SecondUserAllowedOrBloced = true

		GetTimeNow := time.Now()
		room.StartedAt = GetTimeNow

		result, err := MessagePrivateRoomSchema.InsertOne(ctx, &room)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		} else {
			var createdRomm *models.MessagePrivateRoomModels

			query := bson.M{"_id": result.InsertedID}

			MessagePrivateRoomSchema.FindOne(ctx, query).Decode(&createdRomm)
			// Create New Room
			_, err := http.Get("http://" + string(os.Getenv("MAINHOST")) + ":8080/AddChat?ConvID=" + createdRomm.ConversationID.Hex())
			if err != nil {
				log.Fatalln(err)
			}
			// chat.AddNewRoom(createdRomm.ConversationID.Hex())
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"room": createdRomm.ConversationID,
			})
		}

	} else {
		roomResult := MessagePrivateRoomSchema.FindOne(ctx, bson.M{"firstUserID": MainUid, "secondUserID": SecodUid})

		if roomResult.Err() != nil {
			roomResult = MessagePrivateRoomSchema.FindOne(ctx, bson.M{"firstUserID": SecodUid, "secondUserID": MainUid})
		}

		roomResult.Decode(&room)

		// Create New Room
		_, err := http.Get("http://" + string(os.Getenv("MAINHOST")) + ":8080/AddChat?ConvID=" + room.ConversationID.Hex())
		if err != nil {
			log.Fatalln(err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"room": room.ConversationID,
		})
	}

}

// PrivateFunGetRoomData
// GetPrivateRoomeData
func PrivateFunGetRoomData(c *fiber.Ctx) error {
	var MessagePrivateRoomSchema = database.DB.Collection("private_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var room models.MessagePrivateRoomModels

	RoomID, _ := primitive.ObjectIDFromHex(c.Params("RoomID"))

	roomResult := MessagePrivateRoomSchema.FindOne(ctx, bson.M{"_id": RoomID})

	if roomResult.Err() != nil {
		roomResult = MessagePrivateRoomSchema.FindOne(ctx, bson.M{"_id": RoomID})
	}

	roomResult.Decode(&room)

	DefaultSkipNumber := 4
	if room.MessagePayloadFirstCopy != nil && len(room.MessagePayloadFirstCopy)-DefaultSkipNumber >= 0 {
		endNum := len(room.MessagePayloadFirstCopy)
		start := endNum - DefaultSkipNumber

		finald := room.MessagePayloadFirstCopy[start:endNum]
		return c.Status(fiber.StatusOK).JSON(finald)

	} else if room.MessagePayloadSecodCopy != nil && len(room.MessagePayloadFirstCopy)-DefaultSkipNumber >= 0 {
		endNum := len(room.MessagePayloadSecodCopy)
		start := endNum - DefaultSkipNumber

		finald := room.MessagePayloadSecodCopy[start:endNum]
		return c.Status(fiber.StatusOK).JSON(finald)

	} else if len(room.MessagePayloadFirstCopy) > 0 && room.MessagePayloadFirstCopy != nil {

		finald := room.MessagePayloadSecodCopy[:len(room.MessagePayloadFirstCopy)]
		return c.Status(fiber.StatusOK).JSON(finald)

	} else if len(room.MessagePayloadSecodCopy) > 0 && room.MessagePayloadSecodCopy != nil {

		finald := room.MessagePayloadSecodCopy[:len(room.MessagePayloadSecodCopy)]
		return c.Status(fiber.StatusOK).JSON(finald)

	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{})

	}

}

// --------------------- Group ------------------------- //

// GetGroupRoomID
func GetGroupRoomID(c *fiber.Ctx) error {

	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var GroupRoom models.MessageGroupModel
	var JoinModel models.JoinedUsers

	groupName := c.Params("GroupName")
	creatorID := c.Params("CreatorID")

	// check user veryfied

	// chek if sender is veryfied
	var UserSchema = database.DB.Collection("users")

	Puserid, _ := primitive.ObjectIDFromHex(creatorID)
	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": Puserid}).Decode(&checkUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	//
	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}
	////

	numOfDocs, _ := MessageGroupRoomSchema.CountDocuments(ctx, bson.M{"name": groupName, "creator": creatorID})

	if numOfDocs == 0 {

		GroupRoom.Name = groupName
		GroupRoom.Creator = creatorID
		JoinModel.IsAdmin = true
		JoinModel.JoinedUserID = creatorID
		GroupRoom.UsersJoinedList = append(GroupRoom.UsersJoinedList, JoinModel)

		result, err := MessageGroupRoomSchema.InsertOne(ctx, &GroupRoom)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		} else {
			var createdRomm *models.MessageGroupModel

			query := bson.M{"_id": result.InsertedID}

			MessageGroupRoomSchema.FindOne(ctx, query).Decode(&createdRomm)
			_, err := http.Get("http://" + string(os.Getenv("MAINHOST")) + ":8080/AddGroup?ConvID=" + createdRomm.ConversationID.Hex())
			if err != nil {
				log.Fatalln(err)
			}
			// chat.AddNewRoom(createdRomm.ConversationID.Hex())
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				// "room": createdRomm.ConversationID.Hex() + creatorID,
				"room": createdRomm.ConversationID.Hex(),
			})
		}

	} else {
		roomResult := MessageGroupRoomSchema.FindOne(ctx, bson.M{"name": groupName, "creator": creatorID})

		roomResult.Decode(&GroupRoom)

		// Create New Room Connected To Socket Server
		// _, err := http.Get("http://44.203.164.223:8080/AddChat?ConvID=" + GroupRoom.ConversationID.Hex())
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// Create New Room
		_, err := http.Get("http://" + string(os.Getenv("MAINHOST")) + ":8080/AddGroup?ConvID=" + GroupRoom.ConversationID.Hex())
		if err != nil {
			log.Fatalln(err)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			// "room": GroupRoom.ConversationID.Hex() + creatorID,
			"room": GroupRoom.ConversationID.Hex(),
		})
	}

}

// SendMessageToGroupChat
func SendMessageToGroupChat(c *fiber.Ctx) error {
	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var GroupRoom models.MessageGroupModel
	var messagelist models.PayloadGroupList

	if err := c.BodyParser(&messagelist); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}
	if messagelist.Message == "" {
		err := "message Can't be empty"
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	groupID, _ := primitive.ObjectIDFromHex(c.Params("GroupID"))
	userID := c.Params("SenerId")

	MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": groupID}).Decode(&GroupRoom)

	for _, User := range GroupRoom.UsersJoinedList {
		// if User.JoinedUserID == userID && User.IsAdmin {
		// 	fmt.Println("Excute Now")
		// }
		if User.JoinedUserID == userID {
			fmt.Println("Excute Now")
			fmt.Println("Message", messagelist)

			GroupRoom.GroupMessages = append(GroupRoom.GroupMessages, messagelist)

			result, err := MessageGroupRoomSchema.UpdateOne(ctx, bson.M{"_id": GroupRoom.ConversationID}, bson.M{"$set": GroupRoom})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}

			if result.MatchedCount == 1 {
				err := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": GroupRoom.ConversationID}).Decode(&GroupRoom)

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
				}
			}

			c.Locals("userId", User.JoinedUserID)
			c.Locals("ConversationID", GroupRoom.ConversationID.Hex())
			CreateSendUnReeadedMessagesGroupRoomNotification(c)

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				// "room": room,
				"ms": messagelist,
			})

		}
	}
	// if slices.Contains(GroupRoom.UsersJoinedList, name) {

	// }

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"Error": "User Not Registerd Member of The Group",
	})
}

// GetGroupMessageByNumbers
func GetGroupMessageByNumbers(c *fiber.Ctx) error {
	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var GroupRoom models.MessageGroupModel
	// var messagelist models.PayloadGroupList

	roomID, _ := primitive.ObjectIDFromHex(c.Params("roomId"))
	userID := c.Params("userId")

	var body struct {
		From int
	}

	if reflect.TypeOf(body.From) != reflect.TypeOf(1) {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": "body From Should be an init type",
			})
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	DefaultSkipNumber := 4

	err := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": roomID}).Decode(&GroupRoom)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	} else {

		for _, User := range GroupRoom.UsersJoinedList {
			// if User.JoinedUserID == userID && User.IsAdmin {
			// 	fmt.Println("Excute Now")
			// }
			if User.JoinedUserID == userID {
				endNum := 0
				if body.From <= 0 {
					endNum = len(GroupRoom.GroupMessages)
				} else {
					endNum = len(GroupRoom.GroupMessages) - (body.From * DefaultSkipNumber)
				}

				start := endNum - DefaultSkipNumber

				cals := (len(GroupRoom.GroupMessages) - 1) - (DefaultSkipNumber * body.From)
				// fmt.Println("clas", cals)
				if endNum >= 0 && start >= 0 {
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"ms": GroupRoom.GroupMessages[start:endNum],
					})
				} else if endNum >= DefaultSkipNumber {
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"ms": "no message",
					})
				} else if cals > 0 {
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						// "ms": GroupRoom.GroupMessages[start:endNum],
						"ms": GroupRoom.GroupMessages[0:endNum],
					})
				} else if (len(GroupRoom.GroupMessages) - 1) <= DefaultSkipNumber {
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						// "ms": GroupRoom.GroupMessages[start:endNum],
						"ms": GroupRoom.GroupMessages,
					})
				} else {
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"ms": "out of range",
					})
				}

				//

			} else {
				return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
					"Error": "Not Authorized",
				})
			}

		}

	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"Error": "Not founded",
	})

}

// AddNewUserToGroup
func AddNewUserToGroup(c *fiber.Ctx) error {
	// roomId
	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var jUser models.JoinedUsers
	var GroupRoom models.MessageGroupModel

	roomId, _ := primitive.ObjectIDFromHex(c.Params("roomId"))
	// JoinUserId, _ := primitive.ObjectIDFromHex(c.Params("JoinUserId"))
	JoinUserId := c.Params("JoinUserId")

	// chek if sender is veryfied
	var UserSchema = database.DB.Collection("users")

	Puserid, _ := primitive.ObjectIDFromHex(JoinUserId)
	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": Puserid}).Decode(&checkUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}
	////

	// UserIdFromAuth, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	UserIdFromAuth := c.Locals("userId").(string)

	numOfDocs, _ := MessageGroupRoomSchema.CountDocuments(ctx, bson.M{"_id": roomId})

	roomResult := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": roomId})
	roomResult.Decode(&GroupRoom)
	// if numOfDocs != 0 && UserIdFromAuth == GroupRoom.Creator {
	isAuthorized := false

	if numOfDocs != 0 {
		for _, Uid := range GroupRoom.UsersJoinedList {
			if Uid.JoinedUserID == UserIdFromAuth {
				isAuthorized = true
				fmt.Println("auth in loop", isAuthorized)

			}
		}
	}

	// fmt.Println("auth", isAuthorized)
	if numOfDocs != 0 && isAuthorized {
		jUser.IsAdmin = false
		jUser.JoinedUserID = JoinUserId

		for _, u := range GroupRoom.UsersJoinedList {
			if u.JoinedUserID == JoinUserId {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": "user Aleady Exist"})
			}
		}
		GroupRoom.UsersJoinedList = append(GroupRoom.UsersJoinedList, jUser)

		// update
		result, err := MessageGroupRoomSchema.UpdateOne(ctx, bson.M{"_id": GroupRoom.ConversationID}, bson.M{"$set": GroupRoom})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		if result.MatchedCount == 1 {
			err := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": GroupRoom.ConversationID}).Decode(&GroupRoom)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

		c.Locals("RequestID", GroupRoom.ConversationID.Hex())
		c.Locals("SuserId", c.Params("JoinUserId"))
		c.Locals("RuserId", c.Params("userId"))
		CreateSendAddedtoChatRoomNotification(c)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			// "room": room,
			"list": GroupRoom.UsersJoinedList,
		})

	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Error": "Group Not Found Or Not Authorized ",
		})
	}

}

// RemoveMemberFromchatGroup
func RemoveMemberFromchatGroup(c *fiber.Ctx) error {
	// roomId
	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var jUser models.JoinedUsers
	var GroupRoom models.MessageGroupModel

	roomId, _ := primitive.ObjectIDFromHex(c.Params("roomId"))
	JoinUserId := c.Params("JoinUserId")
	UserIdFromAuth := c.Locals("userId").(string)

	numOfDocs, _ := MessageGroupRoomSchema.CountDocuments(ctx, bson.M{"_id": roomId})
	if numOfDocs != 0 {
		roomResult := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": roomId})
		roomResult.Decode(&GroupRoom)
	} else {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{"data": "Not Found!"})
	}

	// if numOfDocs != 0 && UserIdFromAuth == GroupRoom.Creator {
	if numOfDocs != 0 && UserIdFromAuth == GroupRoom.Creator {
		jUser.IsAdmin = false
		jUser.JoinedUserID = JoinUserId

		var NewGropRoom models.MessageGroupModel
		for _, u := range GroupRoom.UsersJoinedList {
			// var model = models.
			if u.JoinedUserID != JoinUserId {
				NewGropRoom.UsersJoinedList = append(NewGropRoom.UsersJoinedList, u)
			}
		}

		GroupRoom.UsersJoinedList = NewGropRoom.UsersJoinedList

		// append(GroupRoom.UsersJoinedList, )
		// GroupRoom.UsersJoinedList = append(GroupRoom.UsersJoinedList, jUser)

		// update
		result, err := MessageGroupRoomSchema.UpdateOne(ctx, bson.M{"_id": GroupRoom.ConversationID}, bson.M{"$set": GroupRoom})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		if result.MatchedCount == 1 {
			err := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": GroupRoom.ConversationID}).Decode(&GroupRoom)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			// "room": room,
			"list": GroupRoom.UsersJoinedList,
		})

		// GroupRoom.Name = groupName
		// GroupRoom.Creator = creatorID
		// JoinModel.IsAdmin = true
		// JoinModel.JoinedUserID = creatorID
		// GroupRoom.UsersJoinedList = append(GroupRoom.UsersJoinedList, JoinModel)

	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Error": "Group Not Found Or Not Authorized ",
		})
	}

}

/// private fuctions for chat service ///

// GetAllPrivateRoomsID
func GetAllPrivateRoomsID(c *fiber.Ctx) error {
	var MessagePrivateRoomSchema = database.DB.Collection("private_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	cursorRooms, err := MessagePrivateRoomSchema.Find(ctx, bson.M{})

	if cursorRooms.RemainingBatchLength() == 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	if err != nil {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	defer cursorRooms.Close(ctx)

	var AllRooms []string

	for cursorRooms.Next(ctx) {
		var msg models.MessagePrivateRoomModels
		cursorRooms.Decode(&msg)
		//fmt.Println("d", &msg)
		AllRooms = append(AllRooms, string(msg.ConversationID.Hex()))
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"data": &AllRooms,
	})

	//return nil
}

// GetAllGroupRoomsID
func GetAllGroupRoomsID(c *fiber.Ctx) error {
	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	cursorRooms, err := MessageGroupRoomSchema.Find(ctx, bson.M{})

	if cursorRooms.RemainingBatchLength() == 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	if err != nil {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": nil,
		})
	}

	defer cursorRooms.Close(ctx)

	var AllRooms []string

	for cursorRooms.Next(ctx) {
		var msg models.MessagePrivateRoomModels
		cursorRooms.Decode(&msg)
		//fmt.Println("d", &msg)
		AllRooms = append(AllRooms, string(msg.ConversationID.Hex()))
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"data": &AllRooms,
	})

}
