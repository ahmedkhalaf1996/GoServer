package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"main/api/models"
	"main/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// secret jwt key
const SecretKey = "secret"

func ManualRegister(c *fiber.Ctx) error {
	fmt.Println("host in api ", os.Getenv("MAINHOST"))

	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var body struct {
		FirstName   string
		LastName    string
		DateOfBirth string
		Email       string
		Password    string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	CheckUser := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&body)

	if CheckUser == nil {
		return c.Status(fiber.StatusAlreadyReported).JSON(
			fiber.Map{
				"message": "user with eamil " + body.Email + " Alrady Exist",
			})
	}

	// hashing password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	// new user data
	newUser := models.UserModel{
		FirstName:         body.FirstName,
		LastName:          body.LastName,
		Email:             body.Email,
		DateOfBirth:       body.DateOfBirth,
		Password:          string(hashPassword),
		IsAccountVerified: false,
	}

	result, err := UsersSchema.InsertOne(ctx, &newUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "can't create new user",
			"Error":   err,
		})
	}
	// get the new user
	var createdUser *models.UserModel

	query := bson.M{"_id": result.InsertedID}

	UsersSchema.FindOne(ctx, query).Decode(&createdUser)
	// create and send the token to client side
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    createdUser.ID.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
	})

	token, _ := claims.SignedString([]byte(SecretKey))

	// call veryfied func
	c.Locals("userId", createdUser.ID.Hex())
	c.Locals("user_email", createdUser.Email)

	res := CreateAndSendVeryFicationCodeToMail(c)
	if res != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Can't Send  veryfied Code. to the provided Email",
		})
	}

	// --------
	CallnotifiationRoom(createdUser.ID.Hex())
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": createdUser, "token": token})
}

// ManualLogin

func ManualLogin(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var user models.UserModel
	var compeardUser models.UserModel

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}
	if err := c.BodyParser(&compeardUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	// 	// check if mail already exeist
	CheckEmail := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}}).Decode(&user)

	if CheckEmail != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Invalid User With Email " + user.Email,
		})
	}

	// check if we have the same pass or not
	CheckPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(compeardUser.Password))

	if CheckPass != nil {
		return c.Status(fiber.StatusFound).JSON(fiber.Map{
			"message": "given Password is not correct !",
			"Error":   string(CheckPass.Error()),
		})
	}

	// create and send the token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}
	CallnotifiationRoom(user.ID.Hex())

	// fmt.Println(token)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": user,
		"token":  token,
	})

}

// Auth With PhoneNumber
func AuthWithPhoneNumber(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	var user models.UserModel
	var body struct {
		FirstName   string
		LastName    string
		PhoneNumber string
		Password    string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	CheckUser := UsersSchema.FindOne(ctx, bson.D{{Key: "phoneNumber", Value: body.PhoneNumber}}).Decode(&user)

	if CheckUser == nil {
		// create and send the token to client side
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
		})

		token, _ := claims.SignedString([]byte(SecretKey))
		CallnotifiationRoom(user.ID.Hex())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": user, "token": token})
	} else {
		// hashing password
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

		// new user data
		newUser := models.UserModel{
			FirstName:         body.FirstName,
			LastName:          body.LastName,
			PhoneNumber:       body.PhoneNumber,
			Password:          string(hashPassword),
			IsAccountVerified: false,
		}

		result, err := UsersSchema.InsertOne(ctx, &newUser)

		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "can't create new user",
				"Error":   err,
			})
		}
		// get the new user
		var createdUser *models.UserModel

		query := bson.M{"_id": result.InsertedID}

		UsersSchema.FindOne(ctx, query).Decode(&createdUser)
		// create and send the token to client side
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    createdUser.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
		})

		token, _ := claims.SignedString([]byte(SecretKey))
		// call veryfied func
		c.Locals("userId", createdUser.ID.Hex())
		c.Locals("user_phone", createdUser.PhoneNumber)
		res := CreateAndSendVeryFicationCodeToPhone(c)
		if res != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Can't Send  veryfied Code. to the provided phone Number",
			})
		}
		CallnotifiationRoom(createdUser.ID.Hex())

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": createdUser, "token": token})

	}

}

// ManualLogin social-Google-Apple-Facebook
// we Should Expect FirstName LastName Email Only And ignore the password
func SocialAuth(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var user models.UserModel

	var body struct {
		FirstName string
		LastName  string
		Email     string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	CheckUser := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	if CheckUser == nil {
		// login user
		// create and send the token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not login",
			})
		}

		// fmt.Println(token)
		CallnotifiationRoom(user.ID.Hex())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": user,
			"token":  token,
		})
	} else {
		// Register User
		// new user data
		newUser := models.UserModel{
			FirstName:         body.FirstName,
			LastName:          body.LastName,
			Email:             body.Email,
			IsAccountVerified: true,
		}

		result, err := UsersSchema.InsertOne(ctx, &newUser)

		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "can't Register new user",
				"Error":   err,
			})
		}
		// get the new user
		var createdUser *models.UserModel

		query := bson.M{"_id": result.InsertedID}

		UsersSchema.FindOne(ctx, query).Decode(&createdUser)
		// create and send the token to client side
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    createdUser.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
		})

		token, _ := claims.SignedString([]byte(SecretKey))
		CallnotifiationRoom(createdUser.ID.Hex())

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": createdUser, "token": token})

	}

}

// Google Auth
func AuthWithGoogle(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var user models.UserModel

	var body struct {
		FirstName string
		LastName  string
		Email     string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	CheckUser := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	if CheckUser == nil {
		// login user
		// create and send the token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not login",
			})
		}
		CallnotifiationRoom(user.ID.Hex())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": user,
			"token":  token,
		})

	} else {
		// register user
		// new data
		newUser := models.UserModel{
			FirstName:         body.FirstName,
			LastName:          body.LastName,
			Email:             body.Email,
			IsAccountVerified: true,
		}

		result, err := UsersSchema.InsertOne(ctx, &newUser)

		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "can't Register new user",
				"Error":   err,
			})
		}

		// get the new user
		var createdUser *models.UserModel

		query := bson.M{"_id": result.InsertedID}

		UsersSchema.FindOne(ctx, query).Decode(&createdUser)
		// create and send the token to client side
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    createdUser.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
		})

		token, _ := claims.SignedString([]byte(SecretKey))
		CallnotifiationRoom(createdUser.ID.Hex())

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": createdUser, "token": token})

	}

}

// FaceBook Auth
func AuthWithFaceBook(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var user models.UserModel

	var body struct {
		FirstName string
		LastName  string
		Email     string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	CheckUser := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	if CheckUser == nil {
		// login user
		// create and send the token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not login",
			})
		}
		CallnotifiationRoom(user.ID.Hex())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": user,
			"token":  token,
		})

	} else {
		// register user
		// new data
		newUser := models.UserModel{
			FirstName:         body.FirstName,
			LastName:          body.LastName,
			Email:             body.Email,
			IsAccountVerified: true,
		}

		result, err := UsersSchema.InsertOne(ctx, &newUser)

		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "can't Register new user",
				"Error":   err,
			})
		}

		// get the new user
		var createdUser *models.UserModel

		query := bson.M{"_id": result.InsertedID}

		UsersSchema.FindOne(ctx, query).Decode(&createdUser)
		// create and send the token to client side
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    createdUser.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
		})

		token, _ := claims.SignedString([]byte(SecretKey))
		CallnotifiationRoom(createdUser.ID.Hex())

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": createdUser, "token": token})

	}

}

// Apple Auth
func AuthWithApple(c *fiber.Ctx) error {
	var UsersSchema = database.DB.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	var user models.UserModel

	var body struct {
		FirstName string
		LastName  string
		Email     string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"Error": err,
			})
	}

	CheckUser := UsersSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	if CheckUser == nil {
		// login user
		// create and send the token
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not login",
			})
		}
		CallnotifiationRoom(user.ID.Hex())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": user,
			"token":  token,
		})

	} else {
		// register user
		// new data
		newUser := models.UserModel{
			FirstName:         body.FirstName,
			LastName:          body.LastName,
			Email:             body.Email,
			IsAccountVerified: true,
		}

		result, err := UsersSchema.InsertOne(ctx, &newUser)

		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "can't Register new user",
				"Error":   err,
			})
		}

		// get the new user
		var createdUser *models.UserModel

		query := bson.M{"_id": result.InsertedID}

		UsersSchema.FindOne(ctx, query).Decode(&createdUser)
		// create and send the token to client side
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    createdUser.ID.Hex(),
			ExpiresAt: time.Now().Add(time.Hour * 720).Unix(), // token will reamain for 30 day
		})

		token, _ := claims.SignedString([]byte(SecretKey))
		CallnotifiationRoom(createdUser.ID.Hex())

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"result": createdUser, "token": token})

	}

}

func CallnotifiationRoom(id string) {
	_, err := http.Get("http://" + string(os.Getenv("MAINHOST")) + ":8090/listenToNotification?UserId=" + id)
	if err != nil {
		log.Fatalln(err)
	}
}

// func VeryfiyAccount
// ID
// Username
// FirstName
// LastName
// FullName
// Email
// Avatar
// Raw
// Token
