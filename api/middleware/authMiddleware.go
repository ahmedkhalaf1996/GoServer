package middleware

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secret"

// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjgzMDIxODYsImlzcyI6IjYzNzAyMjZmNDJjNDk5M2JjNjU0OGRhYiJ9.JAaR1mj_ax8AfU8mfVXjoICImyIT706RBpxAKiizjzQ
// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjU4NTEwMDYsImlzcyI6IjYzNDk0NTRiYmMwNDg2NmI2ZWMzYjFiNyJ9.cZ759tG4fgpWbqfc15FkNDvkYp6O8wBRvqp-qZCSnLc
func AuthMiddleware(c *fiber.Ctx) error {
	tok := c.Get("Authorization")
	// fmt.Println("token", tok)
	// splited := strings.Split(tok, "Bearer ")
	if tok == "" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"Error": "No Token Provided",
		})
	}

	Tokerr := strings.Split(tok, " ")

	// fmt.Println("tok lenght", len(Tokerr))
	if len(Tokerr) != 2 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"Error": "UnExepeted Provided Authorization Token",
		})
	}
	if Tokerr == nil {
		// fmt.Println("tokk", Tokerr[len(Tokerr)-1])
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"Error": "Can't extract token",
		})
	}
	//  else {
	// 	fmt.Println("tokk", Tokerr[len(Tokerr)-1])
	// }

	tok = Tokerr[len(Tokerr)-1]
	token, err := jwt.ParseWithClaims(tok, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	c.Locals("userId", claims.Issuer)
	return c.Next()

}
