package controllers

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"main/api/models"
	"main/api/sendmailsms"
	"main/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"main/api/countrycodes"

	"github.com/gofiber/fiber/v2"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPhoneCountriesCodesList(c *fiber.Ctx) error {

	d := countrycodes.GetCountriesListAlpha2()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": d,
		// "difftime": minute,
	})

	// return nil
}

// --- used by The code Send To Mail
func CreateAndSendVeryFicationCodeToMail(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("veryfied")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var body struct {
		UserID         string `json:"userID" bson:"userID"`
		IsVeryfiyedYet bool   `json:"isVeryfiyedYet" bson:"isVeryfiyedYet"`
		VeryfiyCode    string `json:"veryfiyCode" bson:"veryfiyCode"`
		IsSendedToMail bool   `json:"isSendedToMail" bson:"isSendedToMail"`
		TryNumber      int    `json:"tryNumber" bson:"tryNumber"`
	}
	UidFromAuth := c.Locals("userId").(string)
	body.UserID = UidFromAuth
	body.IsVeryfiyedYet = false
	body.IsSendedToMail = true
	body.TryNumber = 0

	randNum := strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
	body.VeryfiyCode = randNum
	//

	//
	_, err := VeryfiedSchema.InsertOne(ctx, &body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't Send Create veryfied Code.",
			"Error":   err,
		})
	}

	user_email := c.Locals("user_email").(string)
	// send verification code
	sendsms := sendmailsms.SendEmailWithVeryficationCodeToEmail(body.VeryfiyCode, user_email)
	// fmt.Println("sendsms", sendsms)
	if !sendsms {
		// user
		pid, _ := primitive.ObjectIDFromHex(body.UserID)
		UsersSchema.DeleteOne(ctx, bson.D{{Key: "_id", Value: pid}})

		// veryfied
		VeryfiedSchema.DeleteOne(ctx, bson.D{{Key: "userID", Value: body.UserID}})

		return fmt.Errorf("error")

	} else {
		return nil
	}
	// end

	// // fmt.Println("uid", body)
	// return nil
}

// --- used by The code Send To Mobile Phone
func CreateAndSendVeryFicationCodeToPhone(c *fiber.Ctx) error {
	var VeryfiedSchema = database.DB.Collection("veryfied")
	var UsersSchema = database.DB.Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var body struct {
		UserID         string `json:"userID" bson:"userID"`
		IsVeryfiyedYet bool   `json:"isVeryfiyedYet" bson:"isVeryfiyedYet"`
		VeryfiyCode    string `json:"veryfiyCode" bson:"veryfiyCode"`
		IsSendedToMail bool   `json:"isSendedToMail" bson:"isSendedToMail"`
		TryNumber      int    `json:"tryNumber" bson:"tryNumber"`
	}

	UidFromAuth := c.Locals("userId").(string)
	body.UserID = UidFromAuth
	body.IsVeryfiyedYet = false
	body.IsSendedToMail = false
	body.TryNumber = 0

	rand.Seed(time.Now().UnixNano())

	randNum := strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
	body.VeryfiyCode = randNum

	//

	_, err := VeryfiedSchema.InsertOne(ctx, &body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Create veryfied Code.",
			"Error":   err,
		})
	}

	phoneNumber := c.Locals("user_phone").(string)
	// send verification code
	sendsms := sendmailsms.SendSMSWithVeryficationCode(body.VeryfiyCode, phoneNumber)
	// fmt.Println("sendsms", sendsms)
	if !sendsms {
		// user
		pid, _ := primitive.ObjectIDFromHex(body.UserID)
		UsersSchema.DeleteOne(ctx, bson.D{{Key: "_id", Value: pid}})

		// veryfied
		VeryfiedSchema.DeleteOne(ctx, bson.D{{Key: "userID", Value: body.UserID}})

		return fmt.Errorf("error")

	} else {
		return nil
	}
	// end
}

// Use is As End Point to Mail
func ReSendVeryficationCodeToMail(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("veryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	var User models.UserModel

	var Veryfi models.VeryfiyModel

	var Body struct {
		Email    string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	// 	// check if mail already exeist
	CheckEmail := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: Body.Email}}).Decode(&User)

	if CheckEmail != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With Email " + Body.Email,
		})
	}

	// check if we have the same pass or not
	CheckPass := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(Body.Password))

	if CheckPass != nil {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"message": "given Password is not correct !",
		})
	}

	checkifveryfied := VeryfiedSchema.FindOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}).Decode(&Veryfi)

	if checkifveryfied == nil {
		// meaning exitst
		NewCode := strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
		Veryfi.VeryfiyCode = NewCode

		GetTimeNow := time.Now()

		// check diff time
		NextTime := Veryfi.LastUpdated
		Fminute := GetTimeNow.Sub(NextTime).Minutes()
		minute := math.Ceil(Fminute)
		// check if the time <= 2 minute
		if minute <= 3 && Veryfi.TryNumber >= 2 {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"Error": "You Asked to Send The Code many time In short beirut of time Plase Try Againe Later "})
		}

		if Veryfi.TryNumber >= 3 {
			Veryfi.TryNumber = 0
		} else {
			Veryfi.TryNumber = Veryfi.TryNumber + 1
		}

		Veryfi.LastUpdated = GetTimeNow //
		// update

		result, err := VeryfiedSchema.UpdateOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}, bson.M{"$set": Veryfi})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
		if result.MatchedCount == 1 {
			err := VeryfiedSchema.FindOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}).Decode(&Veryfi)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

		// send verification code
		sendsms := sendmailsms.SendEmailWithVeryficationCodeToEmail(Veryfi.VeryfiyCode, Body.Email)

		// sendsms := sendmailsms.SendSMSWithVeryficationCode(Veryfi.VeryfiyCode, Body.PhoneNumber)
		// fmt.Println("sendsms", sendsms)
		if sendsms {
			// user
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"userid":  User.ID,
				"message": "Code Sended Successfully",
			})

		}
		// end

	}

	return nil
}

// Use is As end Point To Phone
func ReSendVeryficationCodeToPhone(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("veryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	var User models.UserModel

	var Veryfi models.VeryfiyModel

	var Body struct {
		PhoneNumber string `json:"phoneNumber" bson:"phoneNumber" validate:"required"`
		Password    string `json:"password" bson:"password" validate:"required,min=5"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// check if Phone already exeist
	CheckEmail := UsersSchema.FindOne(ctx, bson.D{{Key: "phoneNumber", Value: Body.PhoneNumber}}).Decode(&User)

	if CheckEmail != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With Phone Number" + Body.PhoneNumber,
		})
	}

	// check i f we have the correct pass
	CheckPass := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(Body.Password))

	if CheckPass != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "given Password is Not Correct !",
		})
	}

	checkifveryfied := VeryfiedSchema.FindOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}).Decode(&Veryfi)

	if checkifveryfied == nil {
		// meaning exits
		NewPhoenCode := strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
		Veryfi.VeryfiyCode = NewPhoenCode

		GetTimeNow := time.Now()

		// check diff time
		NextTime := Veryfi.LastUpdated
		Fminute := GetTimeNow.Sub(NextTime).Minutes()
		minute := math.Ceil(Fminute)
		// check if the time <= 2 minute
		if minute <= 3 && Veryfi.TryNumber >= 2 {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"Error": "You Asked to Send The Code many time In short beirut of time Plase Try Againe Later "})
		}

		if Veryfi.TryNumber >= 3 {
			Veryfi.TryNumber = 0
		} else {
			Veryfi.TryNumber = Veryfi.TryNumber + 1
		}

		Veryfi.LastUpdated = GetTimeNow
		// update
		result, err := VeryfiedSchema.UpdateOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}, bson.M{"$set": Veryfi})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		if result.MatchedCount == 1 {
			err := VeryfiedSchema.FindOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}).Decode(&Veryfi)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}
		// return Updated One

		// send verification code
		sendsms := sendmailsms.SendSMSWithVeryficationCode(Veryfi.VeryfiyCode, Body.PhoneNumber)
		// fmt.Println("sendsms", sendsms)
		if sendsms {
			// user
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"userid":  User.ID,
				"message": "Code Sended Successfully",
			})

		}
		// end

	}

	return nil
}

// tasks => check number of tring and time between sending
// Use is As End Point To Phone

// veryfied the account only
func VeryfiedProfileAccount(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("veryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	var User models.UserModel

	var Veryfi models.VeryfiyModel

	var Body struct {
		VeryfiyCode string `json:"veryfiyCode" bson:"veryfiyCode"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// userid  // veryfiedCode
	userid := c.Params("userId")
	// veryfiedCode := c.Params("veryfiedCode")

	// get user from db
	PuserID, _ := primitive.ObjectIDFromHex(userid)

	GETUserErr := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PuserID}}).Decode(&User)

	if GETUserErr != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With id " + userid,
		})
	}

	// get veryfied data
	checkifveryfied := VeryfiedSchema.FindOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}).Decode(&Veryfi)

	if checkifveryfied == nil {
		// meaning exits
		// return Updated
		if Veryfi.VeryfiyCode == Body.VeryfiyCode {
			// update db
			User.IsAccountVerified = true
			if Veryfi.IsSendedToMail {
				User.IsPhoneNumberVerified = false
			} else {
				User.IsPhoneNumberVerified = true
			}
			result, err := UsersSchema.UpdateOne(ctx, bson.D{{Key: "_id", Value: User.ID}}, bson.M{"$set": User})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}

			if result.MatchedCount == 1 {
				err := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: User.ID}}).Err()

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
				}
			}

			// update veryfied
			Veryfi.IsVeryfiyedYet = true

			result, err = VeryfiedSchema.UpdateOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}, bson.M{"$set": Veryfi})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}

			if result.MatchedCount == 1 {
				err := VeryfiedSchema.FindOne(ctx, bson.D{{Key: "userID", Value: User.ID.Hex()}}).Err()

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
				}
			}

			c.Locals("userId", User.ID.Hex())
			CreateSendActivateAccountNotification(c)

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Account Veryfied Successfully.",
			})

		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": "Provided Code is Not Correct",
			})
		}
	}

	//
	return nil
}

// // veryfied the account and the phone
// func VeryfiedProfileAccountForPhoen(c *fiber.Ctx) error {
// 	// userid // veryfiedCode

// 	return nil
// }
