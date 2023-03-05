package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/api/models"
	"main/database"
	"os"
	"time"

	"github.com/rgamba/evtwebsocket"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
)

// sended way
// c.Locals("userId", createdUser.ID.Hex())
// c.Locals("user_email", createdUser.Email)

// res := CreateAndSendVeryFicationCodeToMail(c)

//receved  way
// UidFromAuth := c.Locals("userId").(string)
// user_email := c.Locals("user_email").(string)

func CreateSendActivateAccountNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var ActivateAccountModel models.ActivateAccount

	err := NotifySchema.FindOne(ctx, bson.M{"userId": c.Locals("userId").(string)}).Decode(&NotificationModel)

	if err != nil {
		// if not found | created

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Account Veryfied Successfully"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		ActivateAccountModel.ActivateAccount = true
		NotifyDeatilsMOdel.NotifyTypeData = ActivateAccountModel

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		UserIDFromFunc := c.Locals("userId").(string)
		NotificationModel.UserId = UserIDFromFunc
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't Send Notyfication",
				"Error":   err,
			})
		}
	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Account Veryfied Successfully"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		ActivateAccountModel.ActivateAccount = true
		NotifyDeatilsMOdel.NotifyTypeData = ActivateAccountModel
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": c.Locals("userId").(string)}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": c.Locals("userId").(string)}).Decode(&NotificationModel)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}
	// SendNotificationToServer
	SendNotificationToServer(NotificationModel)
	return nil
}

func CreateSendUnReeadedMessagesPrivateRoomNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var UnReeadedMessagesModel models.UnReeadedMessages

	err := NotifySchema.FindOne(ctx, bson.M{"userId": c.Locals("userId").(string)}).Decode(&NotificationModel)

	if err != nil {
		// if not found | created
		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Un Readed Message Conversation ID"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		UnReeadedMessagesModel.UnReeadedMessages = true
		UnReeadedMessagesModel.ConversationID = c.Locals("ConversationID").(string)
		NotifyDeatilsMOdel.NotifyTypeData = UnReeadedMessagesModel

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = c.Locals("userId").(string)
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't Send Notyfication",
				"Error":   err,
			})
		}
	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Un Readed Message Conversation ID"

		GettimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GettimeNow

		UnReeadedMessagesModel.UnReeadedMessages = true
		UnReeadedMessagesModel.ConversationID = c.Locals("ConversationID").(string)

		NotifyDeatilsMOdel.NotifyTypeData = UnReeadedMessagesModel
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)
		NotificationModel.UserId = c.Locals("userId").(string)
		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": c.Locals("userId").(string)}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": c.Locals("userId").(string)}).Decode(&NotificationModel)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}

	// SendNotificationToServer
	SendNotificationToServer(NotificationModel)
	return nil
}

func CreateSendUnReeadedMessagesGroupRoomNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var MessageGroupRoomSchema = database.DB.Collection("group_messages_rooms")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	// var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var UnReeadedMessagesModel models.UnReeadedMessages
	var GroupRoom models.MessageGroupModel

	// get conversetion users
	UidFunc := c.Locals("userId").(string)
	SconvId := c.Locals("ConversationID").(string)
	roomID, _ := primitive.ObjectIDFromHex(c.Locals("ConversationID").(string))

	err := MessageGroupRoomSchema.FindOne(ctx, bson.M{"_id": roomID}).Decode(&GroupRoom)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	for _, el := range GroupRoom.UsersJoinedList {
		if UidFunc != el.JoinedUserID {

			// make job here
			var NastedNotificationModel models.NotificationModel
			err := NotifySchema.FindOne(ctx, bson.M{"userId": el.JoinedUserID}).Decode(&NastedNotificationModel)

			if err != nil {
				NotifyDeatilsMOdel.NotifyNumber = 1
				NotifyDeatilsMOdel.ReadedOrNot = false

				NotifyDeatilsMOdel.NotifyMessage = "Un Readed Message Conversation ID"

				GetTimeNow := time.Now()
				NotifyDeatilsMOdel.SendedAt = GetTimeNow

				UnReeadedMessagesModel.UnReeadedMessages = true
				UnReeadedMessagesModel.ConversationID = SconvId
				NotifyDeatilsMOdel.NotifyTypeData = UnReeadedMessagesModel

				NastedNotificationModel.NotifyList = append(NastedNotificationModel.NotifyList, NotifyDeatilsMOdel)

				NastedNotificationModel.UserId = el.JoinedUserID
				NastedNotificationModel.NumOfUnReadedNotify = 1

				_, err := NotifySchema.InsertOne(ctx, &NastedNotificationModel)

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "can't Send Notyfication",
						"Error":   err,
					})
				}

			} else { // update
				NotifyDeatilsMOdel.NotifyNumber = len(NastedNotificationModel.NotifyList) + 1
				NotifyDeatilsMOdel.ReadedOrNot = false

				NotifyDeatilsMOdel.NotifyMessage = "Un Readed Message Conversation ID"

				GetTimeNow := time.Now()
				NotifyDeatilsMOdel.SendedAt = GetTimeNow

				UnReeadedMessagesModel.UnReeadedMessages = true
				UnReeadedMessagesModel.ConversationID = SconvId

				NotifyDeatilsMOdel.NotifyTypeData = UnReeadedMessagesModel
				NastedNotificationModel.NotifyList = append(NastedNotificationModel.NotifyList, NotifyDeatilsMOdel)
				NastedNotificationModel.UserId = el.JoinedUserID
				NastedNotificationModel.NumOfUnReadedNotify = NastedNotificationModel.NumOfUnReadedNotify + 1

				NotifySchema.UpdateOne(ctx, bson.M{"userId": el.JoinedUserID}, bson.M{"$set": NastedNotificationModel})

			}
			SendNotificationToServer(NastedNotificationModel)

			// -------------
		}
	}
	//

	// err := NotifySchema.FindOne(ctx, bson.M{"userId": c.Locals("userId").(string)}).Decode(&NotificationModel)

	// if err != nil {

	// } else {

	// }
	// SendNotificationToServer
	return nil
}

func CreateSendLoveSendedNotification(c *fiber.Ctx) error {
	// fmt.Println("called")
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var LoveSended models.LoveSended

	RequestID := c.Locals("RequestID").(string)
	SuserId := c.Locals("SuserId").(string)
	RuserId := c.Locals("RuserId").(string)

	err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

	if err != nil {

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Love Sended From User"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		LoveSended.LoveSended = true
		LoveSended.RequestID = RequestID
		LoveSended.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = LoveSended

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = RuserId
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't send otyfication",
				"Error":   err,
			})
		}

	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Love Sended From User"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		LoveSended.LoveSended = true
		LoveSended.RequestID = RequestID
		LoveSended.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = LoveSended
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": RuserId}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}
	SendNotificationToServer(NotificationModel)

	return nil
}

func CreateSendRequestToChatNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var RequestToChat models.RequestToChat

	RequestID := c.Locals("RequestID").(string)
	SuserId := c.Locals("SuserId").(string)
	RuserId := c.Locals("RuserId").(string)

	err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

	if err != nil {

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Request To Chat"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		RequestToChat.RequestToChat = true
		RequestToChat.RequestID = RequestID
		RequestToChat.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = RequestToChat

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = RuserId
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't send otyfication",
				"Error":   err,
			})
		}

	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Request to chat"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		RequestToChat.RequestToChat = true
		RequestToChat.RequestID = RequestID
		RequestToChat.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = RequestToChat
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": RuserId}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}
	SendNotificationToServer(NotificationModel)

	return nil
}

func CreateSendStarSendedNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var StarSended models.StarSended

	RequestID := c.Locals("RequestID").(string)
	SuserId := c.Locals("SuserId").(string)
	RuserId := c.Locals("RuserId").(string)

	err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

	if err != nil {

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "User Send Start to you"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		StarSended.StarSended = true
		StarSended.RequestID = RequestID
		StarSended.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = StarSended

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = RuserId
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't send otyfication",
				"Error":   err,
			})
		}

	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "User Send Start to you"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		StarSended.StarSended = true
		StarSended.RequestID = RequestID
		StarSended.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = StarSended
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": RuserId}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}
	SendNotificationToServer(NotificationModel)

	return nil
}

func CreateSendBuzzSendedNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var BuzzSended models.BuzzSended

	RequestID := c.Locals("RequestID").(string)
	SuserId := c.Locals("SuserId").(string)
	RuserId := c.Locals("RuserId").(string)

	err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

	if err != nil {

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "User Send Buzz to you"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		BuzzSended.BuzzSended = true
		BuzzSended.RequestID = RequestID
		BuzzSended.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = BuzzSended

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = RuserId
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't send otyfication",
				"Error":   err,
			})
		}

	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "User Send Buzz to you"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		BuzzSended.BuzzSended = true
		BuzzSended.RequestID = RequestID
		BuzzSended.SendedFromUserID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = BuzzSended
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": RuserId}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}

	SendNotificationToServer(NotificationModel)
	return nil
}

func CreateSendAddedtoChatRoomNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var AddedtoChatRoom models.AddedtoChatRoom

	RequestID := c.Locals("RequestID").(string)
	SuserId := c.Locals("SuserId").(string)
	RuserId := c.Locals("RuserId").(string)

	err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

	if err != nil {

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Added by user to Chat Room Group"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		AddedtoChatRoom.AddedtoChatRoom = true
		AddedtoChatRoom.ConversationID = RequestID
		AddedtoChatRoom.UserAddedID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = AddedtoChatRoom

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = RuserId
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't send otyfication",
				"Error":   err,
			})
		}

	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Added by user to Chat Room Group"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		AddedtoChatRoom.AddedtoChatRoom = true
		AddedtoChatRoom.ConversationID = RequestID
		AddedtoChatRoom.UserAddedID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = AddedtoChatRoom
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": RuserId}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}

	SendNotificationToServer(NotificationModel)
	return nil
}

func CreateSendInviteUserToEventNotification(c *fiber.Ctx) error {
	var NotifySchema = database.DB.Collection("Notification")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var NotificationModel models.NotificationModel
	var NotifyDeatilsMOdel models.NotifyDeatils
	var InviteUserToEvent models.InviteUserToEvent

	RequestID := c.Locals("RequestID").(string)
	SuserId := c.Locals("SuserId").(string)
	RuserId := c.Locals("RuserId").(string)

	err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

	if err != nil {

		NotifyDeatilsMOdel.NotifyNumber = 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Added by user to Chat Room Group"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		InviteUserToEvent.InviteUserToEvent = true
		InviteUserToEvent.EventID = RequestID
		InviteUserToEvent.UserAddedID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = InviteUserToEvent

		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.UserId = RuserId
		NotificationModel.NumOfUnReadedNotify = 1

		_, err := NotifySchema.InsertOne(ctx, &NotificationModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't send otyfication",
				"Error":   err,
			})
		}

	} else { // update
		NotifyDeatilsMOdel.NotifyNumber = len(NotificationModel.NotifyList) + 1
		NotifyDeatilsMOdel.ReadedOrNot = false

		NotifyDeatilsMOdel.NotifyMessage = "Added by user to Chat Room Group"

		GetTimeNow := time.Now()
		NotifyDeatilsMOdel.SendedAt = GetTimeNow

		InviteUserToEvent.InviteUserToEvent = true
		InviteUserToEvent.EventID = RequestID
		InviteUserToEvent.UserAddedID = SuserId

		NotifyDeatilsMOdel.NotifyTypeData = InviteUserToEvent
		NotificationModel.NotifyList = append(NotificationModel.NotifyList, NotifyDeatilsMOdel)

		NotificationModel.NumOfUnReadedNotify = NotificationModel.NumOfUnReadedNotify + 1

		//
		result, err := NotifySchema.UpdateOne(ctx, bson.M{"userId": RuserId}, bson.M{"$set": NotificationModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		if result.MatchedCount == 1 {
			err := NotifySchema.FindOne(ctx, bson.M{"userId": RuserId}).Decode(&NotificationModel)

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(&fiber.Map{"data": err.Error()})
			}
		}

	}

	SendNotificationToServer(NotificationModel)
	return nil
}

func SendNotificationToServer(data models.NotificationModel) {
	// // Call Room in Chat Service

	// _, err := http.Get("http://localhost:8090/listenToNotification?UserId=" + data.UserId)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // send data to the server
	// // fmt.Println("userid", data.UserId, "......", "restofdata", data)

	// // _, err = http.Post("http://localhost:8090/notification?UserId="+data.UserId, "", nil)
	// // if err != nil {
	// // 	log.Fatalln(err)
	// // }
	conn := evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected!")
		},

		// Ping message to send (optional)
	}
	err := conn.Dial("ws://"+string(os.Getenv("MAINHOST"))+":8090/notification/"+data.UserId, "")
	if err != nil {
		log.Fatal(err)
	}

	// msg := evtwebsocket.Msg{
	// 	Body: []byte("Message body"),
	// }
	x, _ := json.Marshal(data)
	msg := evtwebsocket.Msg{
		Body: []byte(x),
	}
	conn.Send(msg)

}
