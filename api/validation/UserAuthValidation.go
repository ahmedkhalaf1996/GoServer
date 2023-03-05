package validation

import (
	"main/api/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var ValidatorUser = validator.New()

func ValidateUserSocialAuth(c *fiber.Ctx) error {
	var errors []*models.IError
	var body models.SocialAuthModel

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := ValidatorUser.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el models.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}

func ValidateUserPhoneAuth(c *fiber.Ctx) error {
	var errors []*models.IError
	var body models.PhoneAuthModel

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := ValidatorUser.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el models.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}

func ValidateUserManualLogin(c *fiber.Ctx) error {
	var errors []*models.IError
	var body models.ManualLoginModel

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := ValidatorUser.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el models.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}

func ValidateUser(c *fiber.Ctx) error {
	var errors []*models.IError
	var body models.UserModel

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	err := ValidatorUser.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el models.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, &el)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	return c.Next()
}
