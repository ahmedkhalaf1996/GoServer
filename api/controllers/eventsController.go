package controllers

import (
	"context"
	"main/database"
	"math"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"

	"main/api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// create event // api/CreateEvent/:userId
func CreateEvent(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")

	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var EventModel models.EventsModels

	if err := c.BodyParser(&EventModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	UID, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": UID}).Decode(&checkUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": UID}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found",
		})
	}

	Userid := c.Params("userId")
	EventModel.CreatorID = Userid

	if EventModel.EventTitle == "" || EventModel.EventDescription == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": "event Title & description can't be empty",
			})
	}
	// eventStartedOn
	date, error := time.Parse("2006-01-02", EventModel.EventStartedOn)

	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": error,
			})
	}

	EventModel.EventTimeing = date

	result, err := EventSchema.InsertOne(ctx, &EventModel)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't create new Event",
			"Error":   err,
		})
	}

	query := bson.M{"_id": result.InsertedID}

	EventSchema.FindOne(ctx, query).Decode(&EventModel)
	return c.JSON(&fiber.Map{"erro": err, "data": &EventModel})
	// return fmt.Errorf("any")
}

// get event by id api/GetEvent/:eventId
func GetEventByid(c *fiber.Ctx) error {
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var EventModel models.EventsModels

	// eventId := c.Params("eventId")
	// fmt.Println("eventId", eventId)
	eventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))

	err := EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)
	// fmt.Println("d", EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": EventModel,
	})
}

// get events  created by user api/GetEventsCreatedByUser/:userId
func GetEventsCreatedByUser(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var EventModel models.EventsModels
	var EventSModel []models.EventsModels

	userId := c.Params("userId")

	UID, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	err := UserSchema.FindOne(ctx, bson.M{"_id": UID}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found",
		})
	}

	findOptions := options.Find()
	filter := bson.M{"creatorID": userId}

	cursor, err := EventSchema.Find(ctx, filter, findOptions)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Events",
		})
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&EventModel)
		EventSModel = append(EventSModel, EventModel)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": EventSModel,
	})

}

// get events that user intersted api/GetInterstedEvents/:userId
func GetUserNearEvents(c *fiber.Ctx) error {
	// return fmt.Errorf("any")
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var MainUser models.UserModel

	var EventModel models.EventsModels
	var EventSModel []models.EventsModels

	UID, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	err := UserSchema.FindOne(ctx, bson.M{"_id": UID}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found",
		})
	}

	// userId := c.Params("")
	userid, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	err = UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&MainUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	City := MainUser.UserLocation
	SCity := strings.Join(City, " ")

	filterCityLocation := bson.M{}

	findOptionsCity := options.Find()

	filterCityLocation = bson.M{
		"$or": []bson.M{
			{
				"eventLocation": bson.M{
					"$regex": primitive.Regex{
						Pattern: SCity,
						Options: "i",
					},
				},
			},
		},
	}

	cursorEvents, _ := EventSchema.Find(ctx, filterCityLocation, findOptionsCity)

	defer cursorEvents.Close(ctx)

	if cursorEvents.RemainingBatchLength() <= 1 {
		cursorEvents, _ = EventSchema.Find(ctx, bson.M{}, options.Find())
	}

	// userIdStr := c.Params("userId")

	for cursorEvents.Next(ctx) {
		cursorEvents.Decode(&EventModel)
		EventSModel = append(EventSModel, EventModel)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"User": EventSModel,
	})

}

// invite User To an event api/inviteUserToEvent/:eventId/:userId | Notify User
func InviteUserToEvent(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var EventModel models.EventsModels

	eventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))

	err := EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)
	// fmt.Println("d", EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}

	SuserId := c.Params("SuserId")
	RuserId := c.Params("RuserId")

	//
	SUID, _ := primitive.ObjectIDFromHex(SuserId)

	err = UserSchema.FindOne(ctx, bson.M{"_id": SUID}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found : " + SuserId,
		})
	}

	RUID, _ := primitive.ObjectIDFromHex(RuserId)

	err = UserSchema.FindOne(ctx, bson.M{"_id": RUID}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found : " + RuserId,
		})
	}

	c.Locals("RequestID", EventModel.EventID.Hex())
	c.Locals("SuserId", c.Params("SuserId"))
	c.Locals("RuserId", c.Params("RuserId"))
	CreateSendInviteUserToEventNotification(c)

	// instate of returning That Data We need To Passed To ReceverUserAs Notification
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"EventDeatils":  EventModel,
		"SendedUserId":  SuserId,
		"ReceverUserId": RuserId,
	})
}

// select if  UserGoingOrNotToEvent api/UserGoingOrNotToEvent/:eventId/userId

func UserGoingOrNotToEvent(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var MainUser models.UserModel
	var EventModel models.EventsModels

	var Body struct {
		GoingUsers      bool `json:"goingUsers" bson:"goingUsers"`
		MayBeGoingUsers bool `json:"mayBeGoingUsers" bson:"mayBeGoingUsers"`
		NotGoingUsers   bool `json:"notGoingUsers" bson:"notGoingUsers"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	RUID, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": RUID}).Decode(&checkUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": RUID}).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found : " + c.Params("userId"),
		})
	}

	//
	eventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))

	err = EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)
	// fmt.Println("d", EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}

	// user
	StrUserid := c.Params("userId")
	userid, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	err = UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&MainUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found",
		})
	}

	// clear slice first
	// check GoingUsers
	if slices.Contains(EventModel.GoingUsers, StrUserid) {
		i := slices.Index(EventModel.GoingUsers, StrUserid)
		EventModel.GoingUsers = slices.Delete(EventModel.GoingUsers, i, i+1)
	}

	// check MayBeGoingUsers
	if slices.Contains(EventModel.MayBeGoingUsers, StrUserid) {
		i := slices.Index(EventModel.MayBeGoingUsers, StrUserid)
		EventModel.MayBeGoingUsers = slices.Delete(EventModel.MayBeGoingUsers, i, i+1)
	}
	// check NotGoingUsers
	if slices.Contains(EventModel.NotGoingUsers, StrUserid) {
		i := slices.Index(EventModel.NotGoingUsers, StrUserid)
		EventModel.NotGoingUsers = slices.Delete(EventModel.NotGoingUsers, i, i+1)
	}

	if Body.GoingUsers {
		EventModel.GoingUsers = append(EventModel.GoingUsers, StrUserid)
	} else if Body.MayBeGoingUsers {
		EventModel.MayBeGoingUsers = append(EventModel.MayBeGoingUsers, StrUserid)
	} else if Body.NotGoingUsers {
		EventModel.NotGoingUsers = append(EventModel.NotGoingUsers, StrUserid)
	}

	// update Event Model
	result, err := EventSchema.UpdateOne(ctx, bson.M{"_id": eventId}, bson.M{"$set": EventModel})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	// var updatedSug models.SuggestedModel
	if result.MatchedCount == 1 {
		err := EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}
	//
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data":         Body,
			"UpdatedEvent": EventModel,
		})
}

func CheckIfUserGoingOrNotToEvent(c *fiber.Ctx) error {
	// return fmt.Errorf("any")
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var MainUser models.UserModel
	var EventModel models.EventsModels

	var Body struct {
		GoingUsers      bool `json:"goingUsers" bson:"goingUsers"`
		MayBeGoingUsers bool `json:"mayBeGoingUsers" bson:"mayBeGoingUsers"`
		NotGoingUsers   bool `json:"notGoingUsers" bson:"notGoingUsers"`
	}

	//
	RUID, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	ErroR := UserSchema.FindOne(ctx, bson.M{"_id": RUID}).Err()

	if ErroR != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        ErroR.Error(),
			"ErrorMessage": "User with Given id not found : " + c.Params("userId"),
		})
	}
	//

	eventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))
	err := EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)
	// fmt.Println("d", EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}

	// user
	StrUserid := c.Params("userId")
	userid, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	err = UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&MainUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found",
		})
	}

	if slices.Contains(EventModel.GoingUsers, StrUserid) {
		Body.GoingUsers = true
	} else if !slices.Contains(EventModel.GoingUsers, StrUserid) {
		Body.GoingUsers = false
	}

	if slices.Contains(EventModel.MayBeGoingUsers, StrUserid) {
		Body.MayBeGoingUsers = true
	} else if !slices.Contains(EventModel.MayBeGoingUsers, StrUserid) {
		Body.MayBeGoingUsers = false
	}

	if slices.Contains(EventModel.NotGoingUsers, StrUserid) {
		Body.NotGoingUsers = true
	} else if !slices.Contains(EventModel.NotGoingUsers, StrUserid) {
		Body.NotGoingUsers = false
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": Body,
		})

}

// add Comment to Event Page api/AddCommentToEvent/:eventId/

func AddCommentToEventPage(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var MainUser models.UserModel

	var EventModel models.EventsModels
	// var EventSModel []models.EventsModels

	var CommentsModel models.EventComments

	if err := c.BodyParser(&CommentsModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	// user
	// StrUserid := CommentsModel.UserAddedCommentID
	userid, _ := primitive.ObjectIDFromHex(CommentsModel.UserAddedCommentID)

	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&checkUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}

	eventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))

	err = EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&MainUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        err.Error(),
			"ErrorMessage": "User with Given id not found",
		})
	}

	EventModel.EventCommentsbyUsers = append(EventModel.EventCommentsbyUsers, CommentsModel)
	// update model on db
	result, err := EventSchema.UpdateOne(ctx, bson.M{"_id": eventId}, bson.M{"$set": EventModel})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	// var updatedSug models.SuggestedModel
	if result.MatchedCount == 1 {
		err := EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}
	//
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"UpdatedEvent": EventModel,
		})
}

// Update Event Deatils  | Should Notify The Joined Users
func UpdateEventDeatils(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var EventUpdateModel models.EventsUpdateModels
	var EventModel models.EventsModels

	UidFromAuth := c.Locals("userId").(string)
	// eventId := c.Params("eventId")
	userId := c.Params("userId")

	//
	RUID, _ := primitive.ObjectIDFromHex(c.Params("userId"))

	ErroR := UserSchema.FindOne(ctx, bson.M{"_id": RUID}).Err()

	if ErroR != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"Error":        ErroR.Error(),
			"ErrorMessage": "User with Given id not found : " + c.Params("userId"),
		})
	}
	//

	if err := c.BodyParser(&EventUpdateModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	if EventUpdateModel.EventStartedOn != "" {
		date, error := time.Parse("2006-01-02", EventUpdateModel.EventStartedOn)
		if error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"Error": error,
				})
		}

		EventUpdateModel.EventTimeing = date

	}

	// get the event from the db
	PeventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))

	err := EventSchema.FindOne(ctx, bson.M{"_id": PeventId}).Decode(&EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}
	//
	EventUpdateModel.EventID = EventModel.EventID
	EventUpdateModel.CreatorID = EventModel.CreatorID
	// Update The EventUpdated
	if EventUpdateModel.EventTitle == "" {
		EventUpdateModel.EventTitle = EventModel.EventTitle
	}

	if EventUpdateModel.EventDescription == "" {
		EventUpdateModel.EventDescription = EventModel.EventDescription
	}

	if EventUpdateModel.EventCoverImage == "" {
		EventUpdateModel.EventCoverImage = EventModel.EventCoverImage
	}

	// fmt.Println(EventUpdateModel.EventLocation != nil)
	if EventUpdateModel.EventLocation == nil {
		EventUpdateModel.EventLocation = EventModel.EventLocation
	}

	if EventUpdateModel.EventStartedOn == "" {
		EventUpdateModel.EventStartedOn = EventModel.EventStartedOn
	}

	// check the user with given id
	if userId != UidFromAuth {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"Error": "This Event Created By Anothor User You Are Not Authorized to Changed !",
		})
	}

	// update
	result, err := EventSchema.UpdateOne(ctx, bson.M{"_id": PeventId}, bson.M{"$set": EventUpdateModel})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	// var updatedSug models.SuggestedModel
	if result.MatchedCount == 1 {
		err := EventSchema.FindOne(ctx, bson.M{"_id": PeventId}).Decode(&EventUpdateModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(&fiber.Map{
		"data": EventUpdateModel,
	})

}

// Get left Days to Start The Event & Send Notification if it soon
func GetLeftDaysToStartTheEvent(c *fiber.Ctx) error {
	// days := t2.Sub(t1).Hours() / 24
	// fmt.Println(days) // 366

	var EventSchema = database.DB.Collection("events")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var EventModel models.EventsModels

	eventId, _ := primitive.ObjectIDFromHex(c.Params("eventId"))

	err := EventSchema.FindOne(ctx, bson.M{"_id": eventId}).Decode(&EventModel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"serverError":  err.Error(),
			"ErrorMessage": "Event with Given id not found",
		})
	}

	//
	GetTimeNow := time.Now()
	NextTime := EventModel.EventTimeing
	Fdays := NextTime.Sub(GetTimeNow).Hours() / 24
	days := math.Ceil(Fdays)

	// fmt.Println("days", days)
	// Notify if days =< 1
	//
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"days": days,
	})

	// return fmt.Errorf("any")

}
