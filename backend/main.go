package main

import (
	"flag"
	"log"

	"github.com/gofiber/cors/v2"  // Fiber's CORS middleware [[5]]
	"github.com/gofiber/fiber/v2" // Fiber framework [[2]][[5]]

	pkgHttp "github.com/aygoko/EcoMInd/backend/api/user"
	repository "github.com/aygoko/EcoMInd/backend/repository/ram_storage"
	"github.com/aygoko/EcoMInd/backend/usecases/service"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := pkgHttp.NewUserHandler(userService)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: []string{"GET", "POST", "PUT"},
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(fiber.Logger())
	app.Use(fiber.Recovery())

	log.Printf("Starting HTTP server on %s", *addr)
	err := app.Listen(*addr)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
