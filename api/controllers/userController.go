package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/api/models"
	"main/api/sendmailsms"
	"main/database"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SuggestedPartner
func SuggestedPartner(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var SuggestedSchema = database.DB.Collection("suggested")

	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var MainUser models.UserModel
	var users []models.UserModel
	var Foundeduser models.UserModel

	// var Suggestedes []models.SuggestedModel
	var Suggested models.SuggestedModel

	userid, _ := primitive.ObjectIDFromHex(c.Params("id"))

	err := UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&MainUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	//
	if !MainUser.IsAccountVerified {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{"error": "Account Not Verified"})
	}

	// Searching
	// 1 Start By Searching of Parteners With in Same city and Location  Fo The User
	City := MainUser.UserLocation
	SCity := strings.Join(City, " ")

	filterCityLocation := bson.M{}

	findOptionsCity := options.Find()

	filterCityLocation = bson.M{
		"$or": []bson.M{
			{
				"userLocation": bson.M{
					"$regex": primitive.Regex{
						Pattern: SCity,
						Options: "i",
					},
				},
			},
		},
	}

	cursorUsers, _ := UserSchema.Find(ctx, filterCityLocation, findOptionsCity)

	defer cursorUsers.Close(ctx)

	if cursorUsers.RemainingBatchLength() <= 1 {
		cursorUsers, _ = UserSchema.Find(ctx, bson.M{}, options.Find())
		fmt.Println("Cursor", cursorUsers.RemainingBatchLength())
	}

	userIdStr := c.Params("id")

	for cursorUsers.Next(ctx) {
		// var sug models.SugListModel
		var user models.UserModel
		cursorUsers.Decode(&user)

		if userIdStr != user.ID.Hex() {
			users = append(users, user)

		}

		// sug.SugUserID = user.ID.Hex()
		// sug.LoveOrHate = true
		// sug.Score = 0

		// Suggested.MainUid = userIdStr
		// Suggested.SuggestedList = append(Suggested.SuggestedList, sug)
	}

	// Add many to Suggested Collection

	err = SuggestedSchema.FindOne(ctx, bson.M{"mainUid": userIdStr}).Decode(&Suggested)

	compeatedTime := time.Now().Sub(Suggested.CreatedAt).Hours()

	MaxNumber := 0
	var FinalSelectedUser string

	// -----------
	MSCORE := 0
	// var SelectedOne string
	for i := range Suggested.SuggestedList {
		sc := Suggested.SuggestedList[i].Score
		if sc > MSCORE && Suggested.SuggestedList[i].LoveOrHate {
			MSCORE = sc
			FinalSelectedUser = Suggested.SuggestedList[i].SugUserID
		}
	}

	if (err != nil && compeatedTime > 24) && MSCORE == 0 {
		for _, user := range users {
			var sug models.SugListModel
			sug.SugUserID = user.ID.Hex()
			sug.LoveOrHate = true
			sug.Score = 0
			//------------------------------
			for _, item := range user.PhyslcalAttraction {
				for _, i := range MainUser.UserBody {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserDrink {
				for _, i := range MainUser.UserDrink {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserDances {
				for _, i := range MainUser.UserDances {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserHobbyes {
				for _, i := range MainUser.UserHobbyes {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserLanguages {
				for _, i := range MainUser.UserLanguages {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserPets {
				for _, i := range MainUser.UserPets {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserZodlac {
				for _, i := range MainUser.UserZodlac {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			if user.PlaceOfJob == MainUser.PlaceOfJob {
				sug.Score = sug.Score + 1
			}

			if user.School == MainUser.School {
				sug.Score = sug.Score + 1
			}

			if user.UserHeight == MainUser.UserHeight {
				sug.Score = sug.Score + 1
			}

			if user.UserLookFor == MainUser.UserLookFor {
				sug.Score = sug.Score + 1
			}

			if user.UserRelationship == MainUser.UserRelationship {
				sug.Score = sug.Score + 1
			}

			if user.UserRole != MainUser.UserRole {
				sug.Score = sug.Score + 1
			}

			if user.IsUserSmoking == MainUser.IsUserSmoking {
				sug.Score = sug.Score + 1
			}

			//------------------------------

			if sug.Score > MaxNumber {
				FinalSelectedUser = user.ID.Hex()
			}

			MaxNumber = sug.Score
			//-----------------------------
			Suggested.MainUid = userIdStr
			Suggested.SuggestedList = append(Suggested.SuggestedList, sug)
		}

		GetTimeNow := time.Now()
		Suggested.CreatedAt = GetTimeNow
		_, err = SuggestedSchema.InsertOne(ctx, &Suggested)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

	} else if MSCORE == 0 {

		SuggestedSchema.DeleteOne(ctx, bson.M{"mainUid": userIdStr})
		Suggested.SuggestedList = nil
		for _, user := range users {
			var sug models.SugListModel
			sug.SugUserID = user.ID.Hex()
			sug.LoveOrHate = true
			sug.Score = 0
			//------------------------------
			for _, item := range user.PhyslcalAttraction {
				for _, i := range MainUser.UserBody {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserDrink {
				for _, i := range MainUser.UserDrink {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserDances {
				for _, i := range MainUser.UserDances {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserHobbyes {
				for _, i := range MainUser.UserHobbyes {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserLanguages {
				for _, i := range MainUser.UserLanguages {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserPets {
				for _, i := range MainUser.UserPets {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			for _, item := range user.UserZodlac {
				for _, i := range MainUser.UserZodlac {
					if item == i {
						sug.Score = sug.Score + 1
					}
				}
			}

			if user.PlaceOfJob == MainUser.PlaceOfJob {
				sug.Score = sug.Score + 1
			}

			if user.School == MainUser.School {
				sug.Score = sug.Score + 1
			}

			if user.UserHeight == MainUser.UserHeight {
				sug.Score = sug.Score + 1
			}

			if user.UserLookFor == MainUser.UserLookFor {
				sug.Score = sug.Score + 1
			}

			if user.UserRelationship == MainUser.UserRelationship {
				sug.Score = sug.Score + 1
			}

			if user.UserRole != MainUser.UserRole {
				sug.Score = sug.Score + 1
			}

			if user.IsUserSmoking == MainUser.IsUserSmoking {
				sug.Score = sug.Score + 1
			}

			//------------------------------

			if sug.Score > MaxNumber {
				FinalSelectedUser = user.ID.Hex()
			}

			MaxNumber = sug.Score
			//-----------------------------
			Suggested.MainUid = userIdStr
			Suggested.SuggestedList = append(Suggested.SuggestedList, sug)
		}

		GetTimeNow := time.Now()
		Suggested.CreatedAt = GetTimeNow
		SuggestedSchema.DeleteOne(ctx, bson.M{"mainUid": userIdStr})
		_, err = SuggestedSchema.InsertOne(ctx, Suggested)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"Erro": "No Available Suggested For Now Try Again"})

	}
	fmt.Println("F A S", FinalSelectedUser)
	if FinalSelectedUser != "" {
		fid, _ := primitive.ObjectIDFromHex(FinalSelectedUser)

		userResult := UserSchema.FindOne(ctx, bson.M{"_id": fid})

		if userResult.Err() != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User Or user Posts Not found",
			})
		}

		userResult.Decode(&Foundeduser)
		// fmt.Println("MSSSS", Foundeduser)

	}

	// ------

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"User": Foundeduser})

}

// Update user Sug LikeSugOrNot
func UpUserSug(c *fiber.Ctx) error {
	var SuggestedSchema = database.DB.Collection("suggested")
	var UserSchema = database.DB.Collection("users")

	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	var sug models.SuggestedModel

	mainUsserid := c.Params("mid")
	SecondUser := c.Params("nid")

	err := SuggestedSchema.FindOne(ctx, bson.M{"mainUid": mainUsserid}).Decode(&sug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}

	for i := range sug.SuggestedList {
		// fmt.Println(sug.SuggestedList[i])
		if sug.SuggestedList[i].SugUserID == SecondUser {
			sug.SuggestedList[i].LoveOrHate = !sug.SuggestedList[i].LoveOrHate
		}
	}

	result, err := SuggestedSchema.UpdateOne(ctx, bson.M{"mainUid": mainUsserid}, bson.M{"$set": sug})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}
	// err = SuggestedSchema.UpdateOne(ctx, bson.M{"mainUid": mainUsserid},{sug}).Decode(&sug)
	var updatedSug models.SuggestedModel
	if result.MatchedCount == 1 {
		err := SuggestedSchema.FindOne(ctx, bson.M{"mainUid": mainUsserid}).Decode(&updatedSug)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}

	// Get the Suggested
	// MaxNumber := 0
	var FinalSelectedUser string
	var Foundeduser models.UserModel

	// -----------
	MSCORE := 0
	// var SelectedOne string
	for i := range updatedSug.SuggestedList {
		sc := updatedSug.SuggestedList[i].Score
		if sc > MSCORE && updatedSug.SuggestedList[i].LoveOrHate {
			MSCORE = sc
			FinalSelectedUser = updatedSug.SuggestedList[i].SugUserID
		}
	}

	if FinalSelectedUser != "" {
		fid, _ := primitive.ObjectIDFromHex(FinalSelectedUser)

		userResult := UserSchema.FindOne(ctx, bson.M{"_id": fid})

		if userResult.Err() != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "User Or user Posts Not found",
			})
		}

		userResult.Decode(&Foundeduser)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": Foundeduser})
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"data": "No Matched Users Founded!"})
	}

}

// Update User Informatin

func UpdateUserInfo(c *fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	var ctx, _ = context.WithTimeout(context.Background(), 120*time.Second)

	var user models.UserModel
	c.BodyParser(&user)

	userid, _ := primitive.ObjectIDFromHex(c.Params("id"))

	var getingUser models.UserModel
	err := UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&getingUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	// if we have same user id
	UidFromAuth := c.Locals("userId").(string)
	idpath := c.Params("id")

	if UidFromAuth != idpath {
		return c.Status(fiber.StatusForbidden).JSON(&fiber.Map{
			"error": "Forbidden Provided token not For the user With the given id.",
		})
	}
	// if Account Not verified return
	if !getingUser.IsAccountVerified {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Account Not Verified !",
		})
	}

	b, err := json.Marshal(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
	}

	var jsonMap map[string]string
	json.Unmarshal([]byte(string(b)), &jsonMap)

	var Next map[string][]string
	json.Unmarshal([]byte(string(b)), &Next)

	update := bson.M{}

	for key, val := range jsonMap { // strings
		if val != "" && key != "_id" {
			if key != "phoneNumber" && key != "email" {
				update[key] = val
				fmt.Println("k", key, "val", val)
			}
		}
	}

	for key, val := range Next { // arrays
		if val != nil && key != "_id" {
			if key == "userMediaPhoto" {
				update[key] = append(getingUser.UserMediaPhoto, string(val[0]))
			}

			if key == "userMediaVideo" {
				update[key] = append(getingUser.UserMediaVideo, string(val[0]))
			}

			if key == "userLanguages" {
				update[key] = append(getingUser.UserLanguages, string(val[0]))
			}

			if key == "userDances" {
				update[key] = append(getingUser.UserDances, string(val[0]))
			}
			if key == "userLocation" {
				// update[key] = []string{}
				update[key] = append([]string{}, string(val[0]))
			}

			if key == "locationDetails" {
				// update[key] = []string{}

				update[key] = append([]string{}, string(val[0]))
			}

			if key == "userHobbyes" {
				update[key] = append(getingUser.UserHobbyes, string(val[0]))
			}

			if key == "userZodlac" {
				update[key] = append(getingUser.UserZodlac, string(val[0]))
			}

			if key == "userPets" {
				update[key] = append(getingUser.UserPets, string(val[0]))
			}

			if key == "physlcalAttraction" {
				update[key] = append(getingUser.PhyslcalAttraction, string(val[0]))
			}

			if key == "userTurnON" {
				update[key] = append(getingUser.UserTurnON, string(val[0]))
			}

			if key == "userStyle" {
				update[key] = append(getingUser.UserStyle, string(val[0]))
			}

			if key == "userBody" {
				update[key] = append(getingUser.UserBody, string(val[0]))
			}

			if key == "userSesson" {
				update[key] = append(getingUser.UserSesson, string(val[0]))
			}

			if key == "userMovies" {
				update[key] = append(getingUser.UserMovies, string(val[0]))
			}

			if key == "userDrink" {
				update[key] = append(getingUser.UserDrink, string(val[0]))
			}
		}
	}

	// fmt.Println(update)

	result, err := UserSchema.UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	}
	//get updated post details
	var updatedUser models.UserModel
	if result.MatchedCount == 1 {
		err := UserSchema.FindOne(ctx, bson.M{"_id": userid}).Decode(&updatedUser)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": updatedUser})

}

// Get Dances List
func GetDancesList(c *fiber.Ctx) error {

	// 1. Ballet
	// 2. Ballroom
	// 3. Contemporary
	// 4. Hip-Hop
	// 5. Jazz
	// 6. TapDance
	// 7. FolkDance
	// 8. IrishDance
	// 9. ModernDance
	// 10. SwingDance
	// 11. poledance
	// 12. stripdance

	list := struct {
		Ballet       string
		Ballroom     string
		Contemporary string
		HipHop       string
		Jazz         string
		TapDance     string
		FolkDance    string
		IrishDance   string
		ModernDance  string
		SwingDance   string
		poledance    string
		stripdance   string
	}{

		Ballet:       "Ballet",
		Ballroom:     "Ballroom",
		Contemporary: "Contemporary",
		HipHop:       "Hip Hop",
		Jazz:         "Jazz",
		TapDance:     "Tap Dance",
		FolkDance:    "Folk Dance",
		IrishDance:   "Irish Dance",
		ModernDance:  "Modern Dance",
		SwingDance:   "Swing Dance",
		poledance:    "Pole Dance",
		stripdance:   "Strip Dance",
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": list})

}

// AddNewEmailToUserProfile/userId
func AddNewEmailToUserProfile(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("emailveryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var User models.UserModel

	var Veryfi models.AddVeryfiyEmailModel

	var Body struct {
		ProvidedEmail string `json:"providedEmail" bson:"providedEmail"`
	}
	userid := c.Params("userId")

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	PuserID, _ := primitive.ObjectIDFromHex(userid)

	GETUserErr := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PuserID}}).Decode(&User)

	if GETUserErr != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With id " + userid,
		})
	}

	// create and send
	UidFromAuth := c.Params("userId")
	Veryfi.UserID = UidFromAuth
	Veryfi.IsVeryfiyedYet = false
	Veryfi.TryNumber = 0
	Veryfi.ProvidedEmail = Body.ProvidedEmail
	rand.Seed(time.Now().UnixNano())

	randNum := strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
	Veryfi.VeryfiyCode = randNum

	_, err := VeryfiedSchema.InsertOne(ctx, &Veryfi)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't Send Create veryfied Code.",
			"Error":   err,
		})
	}
	//
	sendsms := sendmailsms.SendEmailWithVeryficationCodeToEmail(Veryfi.VeryfiyCode, Body.ProvidedEmail)
	// fmt.Println("sendsms", sendsms)
	if sendsms {
		// user
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Code Sended Successfully",
		})

	}

	return nil
}

// VeryfiedAddedEmailToProfile/userId
func VeryfiedAddedEmailToProfile(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("emailveryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	var User models.UserModel

	var Veryfi models.AddVeryfiyEmailModel

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
		// fmt.Println("vvcode", Veryfi.VeryfiyCode, "***", Body.VeryfiyCode)
		if Veryfi.VeryfiyCode == Body.VeryfiyCode {
			// update db
			User.Email = Veryfi.ProvidedEmail

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

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"data": "Email " + Veryfi.ProvidedEmail + " Added To The User Profile Successfully",
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

// AddNewPhoneNumberToUserProfile/userId
func AddNewPhoneNumberToUserProfile(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("phoneveryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var User models.UserModel

	var Veryfi models.AddVeryfiyPhoneNumberModel

	userid := c.Params("userId")

	var Body struct {
		ProvidedPhoneNumber string `json:"providedPhoneNumber" bson:"providedPhoneNumber"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	PuserID, _ := primitive.ObjectIDFromHex(userid)

	GETUserErr := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PuserID}}).Decode(&User)

	if GETUserErr != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With id " + userid,
		})
	}

	// create and send
	UidFromAuth := c.Params("userId")
	Veryfi.UserID = UidFromAuth
	Veryfi.IsVeryfiyedYet = false
	Veryfi.TryNumber = 0
	Veryfi.ProvidedPhoneNumber = Body.ProvidedPhoneNumber

	rand.Seed(time.Now().UnixNano())

	randNum := strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
	Veryfi.VeryfiyCode = randNum

	//
	_, err := VeryfiedSchema.InsertOne(ctx, &Veryfi)
	sendsms := sendmailsms.SendSMSWithVeryficationCode(Veryfi.VeryfiyCode, Body.ProvidedPhoneNumber)
	// fmt.Println("sendsms", sendsms)
	if sendsms {
		// user
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Code Sended Successfully",
		})

	}

	if err != nil && !sendsms {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't Send Create veryfied Code.",
			"Error":   err,
		})
	}

	return nil
}

// VeryfiedAddedPhoneToProfile
func VeryfiedAddedPhoneToProfile(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var VeryfiedSchema = database.DB.Collection("phoneveryfied")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	var User models.UserModel

	var Veryfi models.AddVeryfiyPhoneNumberModel

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
			User.IsPhoneNumberVerified = true
			User.PhoneNumber = Veryfi.ProvidedPhoneNumber
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

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"data": "Phone Number " + Veryfi.ProvidedPhoneNumber + " Added To The User Profile Successfully",
			})

		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": "Provided Code is Not Correct",
			})
		}
	}

	return nil
}

// Add Request To chat with user |  RequestToChat/:SuserId/:RuserId
func RequestToChat(c *fiber.Ctx) error {

	var UsersSchema = database.DB.Collection("users")
	var RequestesSchema = database.DB.Collection("requests_chat")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	// model
	var Requestes models.UserChatRequestes

	// sederRequestUserid
	suserId := c.Params("SuserId")
	// ReceverRequestUserid
	ruserId := c.Params("RuserId")

	// check if suser is found
	PSuid, _ := primitive.ObjectIDFromHex(suserId)
	err := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PSuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid User With id " + suserId,
		})
	}

	// check if rsuer is found
	PRuid, _ := primitive.ObjectIDFromHex(ruserId)
	err = UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PRuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid User With id " + ruserId,
		})
	}

	// get request if exist
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requestes)

	if err == nil {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": Requestes,
		})
	}

	// set id's on models
	Requestes.SenderUserID = suserId
	Requestes.ReceverUserID = ruserId
	Requestes.IsAcceptedYet = false

	GetTimeNow := time.Now()
	Requestes.SendedAt = GetTimeNow

	//--------insert data on db-----------
	_, err = RequestesSchema.InsertOne(ctx, &Requestes)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// get created request
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requestes)

	c.Locals("RequestID", Requestes.RequestID.Hex())
	c.Locals("SuserId", c.Params("SuserId"))
	c.Locals("RuserId", c.Params("RuserId"))
	CreateSendRequestToChatNotification(c)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": Requestes,
	})

}

// send love to user || notify the next user
func SendLovetoUser(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var RequestesSchema = database.DB.Collection("requests_love")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	// model
	var Requests models.UserLoveRequestes

	// senderrequestUserid
	suserId := c.Params("SuserId")
	// ReceverRequestUserid
	ruserId := c.Params("RuserId")

	// check if suser is found
	PSuid, _ := primitive.ObjectIDFromHex(suserId)
	err := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PSuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid User With id " + suserId,
		})
	}

	// check i ruser is found
	PRuid, _ := primitive.ObjectIDFromHex(ruserId)
	err = UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PRuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid User With id " + ruserId,
		})
	}

	// get request if exist
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requests)

	if err == nil {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": Requests,
		})
	}

	// set id's on models
	Requests.SenderUserID = suserId
	Requests.ReceverUserID = ruserId
	Requests.IsAcceptedYet = false

	GetTimeNow := time.Now()
	Requests.SendedAt = GetTimeNow

	// --------- insert data on db ------------
	_, err = RequestesSchema.InsertOne(ctx, &Requests)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// get created request
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requests)

	c.Locals("RequestID", Requests.RequestID.Hex())
	c.Locals("SuserId", c.Params("SuserId"))
	c.Locals("RuserId", c.Params("RuserId"))
	CreateSendLoveSendedNotification(c)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": Requests,
	})

}

// send buzz to user || notify the next user
func SendBuzzToUser(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var RequestesSchema = database.DB.Collection("requests_buzz")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	// model
	var Requests models.UserBuzzRequestes

	// senderrequestUserid
	suserId := c.Params("SuserId")
	// ReceverRequestUserid
	ruserId := c.Params("RuserId")

	// check if suser is found
	PSuid, _ := primitive.ObjectIDFromHex(suserId)
	err := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PSuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid UserWith id " + suserId,
		})
	}

	// check i ruser is found
	PRuid, _ := primitive.ObjectIDFromHex(ruserId)
	err = UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PRuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid User With id " + ruserId,
		})
	}

	// get request if exist
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requests)

	if err == nil {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": Requests,
		})
	}

	// set id's on models
	Requests.SenderUserID = suserId
	Requests.ReceverUserID = ruserId
	Requests.IsAcceptedYet = false

	GettimeNow := time.Now()
	Requests.SendedAt = GettimeNow

	// ---------- insert data on db -----------------
	_, err = RequestesSchema.InsertOne(ctx, &Requests)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//  get created request
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requests)

	c.Locals("RequestID", Requests.RequestID.Hex())
	c.Locals("SuserId", c.Params("SuserId"))
	c.Locals("RuserId", c.Params("RuserId"))
	CreateSendBuzzSendedNotification(c)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": Requests,
	})

}

// send Star to user || notify the next user
func SendStarToUser(c *fiber.Ctx) error {

	var UsersSchema = database.DB.Collection("users")
	var RequestesSchema = database.DB.Collection("requests_star")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	// model
	var Requests models.UserStarRequestes

	// senderrequestUserid
	suserId := c.Params("SuserId")

	// ReceverRequestUserid
	ruserId := c.Params("RuserId")

	// check if suser is found
	PSuid, _ := primitive.ObjectIDFromHex(suserId)
	err := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PSuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid UserWith id " + suserId,
		})
	}

	// check i ruser is found
	PRuid, _ := primitive.ObjectIDFromHex(ruserId)
	err = UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PRuid}}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid User With id " + ruserId,
		})
	}

	// get request if exist
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requests)

	if err == nil {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"data": Requests,
		})
	}

	// set id's on models

	Requests.SenderUserID = suserId
	Requests.ReceverUserID = ruserId
	Requests.IsAcceptedYet = false

	GettimeNow := time.Now()
	Requests.SendedAt = GettimeNow

	// ---- insert data on Db ----------------
	_, err = RequestesSchema.InsertOne(ctx, &Requests)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// get created request
	err = RequestesSchema.FindOne(ctx, bson.M{"senderUserID": suserId, "receverUserID": ruserId}).Decode(&Requests)

	c.Locals("RequestID", Requests.RequestID.Hex())
	c.Locals("SuserId", c.Params("SuserId"))
	c.Locals("RuserId", c.Params("RuserId"))
	CreateSendStarSendedNotification(c)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": Requests,
	})
}

// accepet refuese requests || Notifiy the user Reaction ||
// api/AcceptRefuaseRequests/:ReqId/:userId
func AcceptRefuaseRequests(c *fiber.Ctx) error {

	var UsersSchema = database.DB.Collection("users")

	var StarSchema = database.DB.Collection("requests_star")
	var BuzzSchema = database.DB.Collection("requests_buzz")
	var LoveSchema = database.DB.Collection("requests_love")

	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var StarModel models.UserStarRequestes
	var BuzzModel models.UserBuzzRequestes
	var LoveModel models.UserLoveRequestes

	var Body struct {
		IsAcceptedYet bool `json:"isAcceptedYet" bson:"isAcceptedYet"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	// userid
	userId := c.Params("userId")

	// get user from db
	PuserID, _ := primitive.ObjectIDFromHex(userId)

	GETUserErr := UsersSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PuserID}}).Err()

	if GETUserErr != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With id " + userId,
		})
	}

	// ReqId
	ReqId := c.Params("ReqId")

	PreQuestID, _ := primitive.ObjectIDFromHex(ReqId)

	StarCount := StarSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PreQuestID}}).Decode(&StarModel)
	BuzzCount := BuzzSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PreQuestID}}).Decode(&BuzzModel)
	LoveCount := LoveSchema.FindOne(ctx, bson.D{{Key: "_id", Value: PreQuestID}}).Decode(&LoveModel)

	if StarCount == nil {
		StarModel.IsAcceptedYet = Body.IsAcceptedYet

		result, err := StarSchema.UpdateOne(ctx, bson.D{{Key: "_id", Value: PreQuestID}}, bson.M{"$set": StarModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
		//get updated data

		if result.MatchedCount == 1 {
			err := StarSchema.FindOne(ctx, bson.M{"_id": PreQuestID}).Decode(&StarModel)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

	} else if BuzzCount == nil {

		BuzzModel.IsAcceptedYet = Body.IsAcceptedYet

		result, err := BuzzSchema.UpdateOne(ctx, bson.D{{Key: "_id", Value: PreQuestID}}, bson.M{"$set": BuzzModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
		//get updated data

		if result.MatchedCount == 1 {
			err := BuzzSchema.FindOne(ctx, bson.M{"_id": PreQuestID}).Decode(&BuzzModel)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

	} else if LoveCount == nil {
		LoveModel.IsAcceptedYet = Body.IsAcceptedYet

		result, err := LoveSchema.UpdateOne(ctx, bson.D{{Key: "_id", Value: PreQuestID}}, bson.M{"$set": LoveModel})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}
		//get updated data

		if result.MatchedCount == 1 {
			err := LoveSchema.FindOne(ctx, bson.M{"_id": PreQuestID}).Decode(&LoveModel)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}

	} else {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Invalid Request With id " + ReqId,
		})
	}

	// return true
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"data": "Update Request Successfully",
	})
}

// -------------- block user sys --------------
// block user | addUserToBlocking User List For Anther User
// AddUserToBlockingList/:mainUID/secUID
func AddUserToBlockingList(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	var BlockedSchema = database.DB.Collection("blocked_list")

	var ctx, _ = context.WithTimeout(context.Background(), 120*time.Second)

	// var userMod models.UserModel

	var BlockedMod models.BlockingModel
	var sBlcokModList models.BlockedListModel
	// check mainuid
	mainUID := c.Params("mainUID")
	PmainUID, _ := primitive.ObjectIDFromHex(mainUID)

	err := UsersSchema.FindOne(ctx, bson.M{"_id": PmainUID}).Err()

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// check mainuid
	secUID := c.Params("secUID")
	PsecUID, _ := primitive.ObjectIDFromHex(secUID)

	err = UsersSchema.FindOne(ctx, bson.M{"_id": PsecUID}).Err()

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// get main user BlockedList
	err = BlockedSchema.FindOne(ctx, bson.M{"mainUid": mainUID}).Decode(&BlockedMod)
	//

	if err != nil {
		BlockedMod.MainUid = mainUID
		sBlcokModList.BlockedUserID = secUID
		BlockedMod.BlockedList = append(BlockedMod.BlockedList, sBlcokModList)
		_, err = BlockedSchema.InsertOne(ctx, &BlockedMod)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

	} else {
		sBlcokModList.BlockedUserID = secUID

		if slices.Contains(BlockedMod.BlockedList, sBlcokModList) {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"data": "user alredy exists",
			})
		}

		BlockedMod.MainUid = mainUID
		sBlcokModList.BlockedUserID = secUID
		BlockedMod.BlockedList = append(BlockedMod.BlockedList, sBlcokModList)

		result, err := BlockedSchema.UpdateOne(ctx, bson.M{"mainUid": mainUID}, bson.M{"$set": BlockedMod})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		if result.MatchedCount == 1 {
			err := BlockedSchema.FindOne(ctx, bson.M{"mainUid": mainUID}).Decode(&BlockedMod)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
			}
		}
	}

	// return blcokedmod

	c.Status(fiber.StatusOK)
	return c.JSON(&fiber.Map{
		"data": &BlockedMod,
	})

}

// ubblock user api/UnBlockUser/:mainUID/:secUID
func UnBlockUser(c *fiber.Ctx) error {

	var UsersSchema = database.DB.Collection("users")
	var BlockedSchema = database.DB.Collection("blocked_list")

	var ctx, _ = context.WithTimeout(context.Background(), 120*time.Second)

	var BlockedMod models.BlockingModel
	var sBlcokModList models.BlockedListModel

	mainUID := c.Params("mainUID")
	PmainUID, _ := primitive.ObjectIDFromHex(mainUID)

	err := UsersSchema.FindOne(ctx, bson.M{"_id": PmainUID}).Err()

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// check
	secUID := c.Params("secUID")
	PsecUID, _ := primitive.ObjectIDFromHex(secUID)

	err = UsersSchema.FindOne(ctx, bson.M{"_id": PsecUID}).Err()

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// get main user bockedlist
	err = BlockedSchema.FindOne(ctx, bson.M{"mainUid": mainUID}).Decode(&BlockedMod)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Blocked List is empty Or Not Created Yet",
		})
	}
	// start wroring
	sBlcokModList.BlockedUserID = secUID

	if slices.Contains(BlockedMod.BlockedList, sBlcokModList) {
		i := slices.Index(BlockedMod.BlockedList, sBlcokModList)
		BlockedMod.BlockedList = slices.Delete(BlockedMod.BlockedList, i, i+1)

		result, err := BlockedSchema.UpdateOne(ctx, bson.M{"mainUid": mainUID}, bson.M{"$set": BlockedMod})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
		}

		if result.MatchedCount == 1 {
			err := BlockedSchema.FindOne(ctx, bson.M{"mainUid": mainUID}).Decode(&BlockedMod)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})

			}
		}

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "user With Given id Not exists user Blocking List",
		})
	}

	// return blockedmod
	c.Status(fiber.StatusOK)
	return c.JSON(&fiber.Map{
		"data": &BlockedMod,
	})

}

// api/Getblockeduserslist/:mainUID
func Getblockeduserslist(c *fiber.Ctx) error {

	var UsersSchema = database.DB.Collection("users")
	var BlockedSchema = database.DB.Collection("blocked_list")

	var ctx, _ = context.WithTimeout(context.Background(), 120*time.Second)

	var BlockedMod models.BlockingModel

	mainUID := c.Params("mainUID")
	PmainUID, _ := primitive.ObjectIDFromHex(mainUID)

	err := UsersSchema.FindOne(ctx, bson.M{"_id": PmainUID}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// get main user blockedlist
	err = BlockedSchema.FindOne(ctx, bson.M{"mainUid": mainUID}).Decode(&BlockedMod)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(&fiber.Map{"error": "User Block List is Empty"})
	} else {
		c.Status(fiber.StatusOK)
		return c.JSON(&fiber.Map{
			"data": BlockedMod.BlockedList,
		})
	}

}

// check if user blocked or not by mainuid suid
// IsUserBlockedOrNot/:mainUID/:secUID
func IsUserBlockedOrNot(c *fiber.Ctx) error {

	var UsersSchema = database.DB.Collection("users")
	var BlockedSchema = database.DB.Collection("blocked_list")

	var ctx, _ = context.WithTimeout(context.Background(), 120*time.Second)

	var BlockedMod models.BlockingModel
	var sBlcokModList models.BlockedListModel

	mainUID := c.Params("mainUID")
	PmainUID, _ := primitive.ObjectIDFromHex(mainUID)

	err := UsersSchema.FindOne(ctx, bson.M{"_id": PmainUID}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// check
	secUID := c.Params("secUID")
	PsecUID, _ := primitive.ObjectIDFromHex(secUID)

	err = UsersSchema.FindOne(ctx, bson.M{"_id": PsecUID}).Err()

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(&fiber.Map{"error": err.Error()})
	}

	// get main user blcoedlist
	err = BlockedSchema.FindOne(ctx, bson.M{"mainUid": mainUID}).Decode(&BlockedMod)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Blocked List is empty Or Not Created Yet",
		})
	}

	sBlcokModList.BlockedUserID = secUID

	if slices.Contains(BlockedMod.BlockedList, sBlcokModList) {
		c.Status(fiber.StatusOK)
		return c.JSON(&fiber.Map{
			"IsBlocked": true,
		})
	} else {
		c.Status(fiber.StatusOK)
		return c.JSON(&fiber.Map{
			"IsBlocked": false,
		})
	}

}

func GetAllUsers(c *fiber.Ctx) error {
	var usersSchema = database.DB.Collection("users")
	var ctx, _ = context.WithTimeout(context.Background(), 300*time.Second)

	cursorRooms, err := usersSchema.Find(ctx, bson.M{})

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
		var msg models.UserModel
		cursorRooms.Decode(&msg)
		//fmt.Println("d", &msg)
		AllRooms = append(AllRooms, string(msg.ID.Hex()))
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"data": &AllRooms,
	})

	//return nil
}
