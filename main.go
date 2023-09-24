package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"one-repo/internal/domain/entity"
	"one-repo/internal/handler"
	"one-repo/internal/handler/middleware"
	"one-repo/internal/infrastructure/database"
	infraRepo "one-repo/internal/infrastructure/repository"
	"one-repo/internal/usecase"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	db, _ := database.NewDB()

	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	userRepo := infraRepo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Get("/users/me", middleware.JwtMiddleware(), userHandler.GetMe)

	app.Listen(":3000")
}
