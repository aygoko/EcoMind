package main

import (
	"flag"
	"log"

	"github.com/gofiber/cors/v2"  // Fiber's CORS middleware [[5]]
	"github.com/gofiber/fiber/v2" // Fiber framework [[2]][[5]]
	"github.com/gofiber/fiber/v2/middleware/logger"

	repository "github.com/aygoko/EcoMInd/backend/repository/database"
	"github.com/aygoko/EcoMInd/backend/usecases/service"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	redisPool := initRedisPool()

	runMigrations(db)

	userRepo := repository.NewUserRepository(db, redisPool)

	userService := service.NewUserService(userRepo, redisPool)

	userHandler := user.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		AppName: "IQJ",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	userHandler.RegisterRoutes(app)

	log.Printf("Запуск сервера на %s", *addr)
	if err := app.Listen(*addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func connectToDB() (*gorm.DB, error) {
	dsn := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisHost+":"+redisPort)
		},
	}
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&repository.User{})
}
