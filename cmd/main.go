package main

import (
	"fmt"
	"main/database"
	_ "main/docs"
	"net"
	"os"

	"main/api/routes"

	"main/Services/chat"
	// "main/Services/liveservice"
	// "main/Services/notification"

	// "github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Some error occured. Err: %s", err)
	// }

	fmt.Println("g", os.Getenv("TRY"))

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			os.Setenv("MAINHOST", fmt.Sprint(ipv4))
		}
	}

	// --- run services

	go chat.Run()
	// go notification.Run()
	// go liveservice.Run()
	// --- run services

	database.Connect()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
	}))

	routes.Setup(app)
	// run server
	app.Static("/Services", "./Services/socketDocs/index.html")

	app.Get("/*", swagger.HandlerDefault) // default

	app.Listen(":80")

}
