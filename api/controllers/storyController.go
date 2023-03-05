package controllers

import (
	"context"
	"fmt"
	"main/api/models"
	"main/database"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// AddStory
func AddStory(c *fiber.Ctx) error {
	var StorySchema = database.DB.Collection("story")
	// var UserSchema = database.DB.Collection("users")

	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var StoryModel models.StoryMainModel
	var StoryListModel models.StoryList

	if err := c.BodyParser(&StoryListModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	Userid := c.Params("userId")

	var UserSchema = database.DB.Collection("users")

	Puserid, _ := primitive.ObjectIDFromHex(Userid)
	var checkUser models.UserModel

	err := UserSchema.FindOne(ctx, bson.M{"_id": Puserid}).Decode(&checkUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	if !checkUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}

	// check expire Storyes
	RemoveExpireStory(&fiber.Ctx{}, Userid)
	// PrimitiveUserID ,_ := primitive.ObjectIDFromHex(Userid)

	err = StorySchema.FindOne(ctx, bson.M{"userId": Userid}).Decode(&StoryModel)

	if err != nil { // Create One
		newStoryContainer := models.StoryMainModel{
			UserId:       Userid,
			NumOfStoryes: 0,
		}

		result, err := StorySchema.InsertOne(ctx, &newStoryContainer)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "can't create new Story Container",
				"Error":   err,
			})
		}
		// get the new user

		query := bson.M{"_id": result.InsertedID}

		StorySchema.FindOne(ctx, query).Decode(&StoryModel)
		// return c.JSON(&fiber.Map{"data": StoryModel})
	}

	// list
	GetTimeNow := time.Now()
	StoryListModel.SendedAt = GetTimeNow
	StoryListModel.StoryNumber = StoryModel.NumOfStoryes + 1

	// main
	StoryModel.StoryPayload = append(StoryModel.StoryPayload, StoryListModel)
	StoryModel.NumOfStoryes = StoryModel.NumOfStoryes + 1

	// update

	result, err := StorySchema.UpdateOne(ctx, bson.M{"userId": Userid}, bson.M{"$set": StoryModel})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}
	if result.MatchedCount == 1 {
		err := StorySchema.FindOne(ctx, bson.M{"userId": Userid}).Decode(&StoryModel)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}

	return c.JSON(&fiber.Map{"erro": err, "data": &StoryModel})

}

// GetStory
func GetStory(c *fiber.Ctx) error {
	var StorySchema = database.DB.Collection("story")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)
	var StoryModel models.StoryMainModel
	// var StoryListModel models.StoryList

	Userid := c.Params("userId")
	// check expire Storyes
	RemoveExpireStory(&fiber.Ctx{}, Userid)

	err := StorySchema.FindOne(ctx, bson.M{"userId": Userid}).Decode(&StoryModel)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"data": "No Story",
			})
	} else {
		// return data
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"data": StoryModel,
			})
	}

}

// // RemoveStory
func RemoveStory(c *fiber.Ctx) error {
	var StorySchema = database.DB.Collection("story")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var StoryModel models.StoryMainModel
	// var StoryListModel models.StoryList

	Userid := c.Params("userId")
	StoryNum := c.Params("StoryNumber")

	// check expire Storyes
	RemoveExpireStory(&fiber.Ctx{}, Userid)

	err := StorySchema.FindOne(ctx, bson.M{"userId": Userid}).Decode(&StoryModel)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{"error": err})
	} else {
		newArrStory := StoryModel.StoryPayload
		newArrStory = nil
		var next models.StoryList
		for _, x := range StoryModel.StoryPayload {
			int1, err := strconv.Atoi(StoryNum)

			if err != nil {
				return c.Status(fiber.StatusOK).JSON(
					fiber.Map{
						"error": "Invalid Number",
					})
			}
			// ---
			if x.StoryNumber == int1 {
				if StoryModel.NumOfStoryes-1 >= 0 {
					StoryModel.NumOfStoryes = StoryModel.NumOfStoryes - 1
				}
				continue
			} else if x.StoryNumber > int1 {
				next.StoryNumber = x.StoryNumber - 1
				next.StoryData = x.StoryData
				next.StoryCaption = x.StoryCaption
				next.IsTypeVideo = x.IsTypeVideo
				next.SendedAt = x.SendedAt

				newArrStory = append(newArrStory, next)
			} else {
				next.StoryNumber = x.StoryNumber
				next.StoryData = x.StoryData
				next.StoryCaption = x.StoryCaption
				next.IsTypeVideo = x.IsTypeVideo
				next.SendedAt = x.SendedAt

				newArrStory = append(newArrStory, next)
			}
		}

		StoryModel.StoryPayload = newArrStory

		result, err := StorySchema.UpdateOne(ctx, bson.M{"userId": Userid}, bson.M{"$set": StoryModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
		if result.MatchedCount == 1 {
			err := StorySchema.FindOne(ctx, bson.M{"userId": Userid}).Decode(&StoryModel)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"data": StoryModel,
				// "next":   next,
				// "newArr": newArrStory,
			})
	}
}

// // RemoveExpireStory
func RemoveExpireStory(c *fiber.Ctx, Userid string) error {
	fmt.Println("Remove expire called!! ")

	var StorySchema = database.DB.Collection("story")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)
	var StoryModel models.StoryMainModel
	// var StoryListModel models.StoryList

	err := StorySchema.FindOne(ctx, bson.M{"userId": Userid}).Decode(&StoryModel)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"data": "No Story",
			})
	} else {

		for _, data := range StoryModel.StoryPayload {
			m := time.Now().Month()
			d := time.Now().Day()
			if (data.SendedAt.Month() != m || d-data.SendedAt.Day() > 1) || (d == data.SendedAt.Day() && data.SendedAt.Hour() > 23) {
				//remove from list
				newArrStory := StoryModel.StoryPayload
				newArrStory = nil
				var next models.StoryList

				for _, x := range StoryModel.StoryPayload {
					// ---
					if x.StoryNumber == data.StoryNumber {
						if StoryModel.NumOfStoryes-1 >= 0 {
							StoryModel.NumOfStoryes = StoryModel.NumOfStoryes - 1
						}
						continue
					} else if x.StoryNumber > data.StoryNumber {
						next.StoryNumber = x.StoryNumber - 1
						next.StoryData = x.StoryData
						next.StoryCaption = x.StoryCaption
						next.IsTypeVideo = x.IsTypeVideo
						next.SendedAt = x.SendedAt

						newArrStory = append(newArrStory, next)
					} else {
						next.StoryNumber = x.StoryNumber
						next.StoryData = x.StoryData
						next.StoryCaption = x.StoryCaption
						next.IsTypeVideo = x.IsTypeVideo
						next.SendedAt = x.SendedAt

						newArrStory = append(newArrStory, next)
					}
				}

				StoryModel.StoryPayload = newArrStory

				_, err := StorySchema.UpdateOne(ctx, bson.M{"userId": Userid}, bson.M{"$set": StoryModel})

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
				}

				continue
			}
		}

	}

	return nil
}
