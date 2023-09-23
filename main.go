package main

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"log"
	"one-repo/internal/domain/entity"
	"one-repo/internal/handler"
	"one-repo/internal/infrastructure/database"
	infraRepo "one-repo/internal/infrastructure/repository"
	"one-repo/internal/usecase"
)

func main() {
	app := fiber.New()
	db, _ := database.NewDB()

	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	userRepo := infraRepo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey:  usecase.JwtSecret,
		TokenLookup: "header:Authorization",
		SuccessHandler: func(c *fiber.Ctx) error {
			tokenStr := c.Get("Authorization")
			claims, err := usecase.ParseToken(tokenStr)
			if err != nil {
				log.Printf("ErrorHandler called: %v", err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized?"})
			}

			c.Locals("user_id", claims.UserID)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("ErrorHandler called: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		},
	})

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Get("/users/me", jwtMiddleware, userHandler.GetMe)

	app.Listen(":3000")
}
