package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"log"
	"one-repo/internal/usecase"
)

func JwtMiddleware() fiber.Handler {
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
	return jwtMiddleware
}
