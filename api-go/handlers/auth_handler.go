package handlers

import (
	"time"

	"api-go/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type TokenRequest struct {
	Username string `json:"username"`
}

// GenerateToken crea un token JWT firmado válido por 24 horas.
func (h *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	var req TokenRequest
	_ = c.BodyParser(&req)

	username := req.Username
	if username == "" {
		username = "interseguro-user"
	}

	claims := jwt.MapClaims{
		"sub":  username,
		"iss":  "interseguro-challenge",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"role": "engineer",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error al firmar el token JWT",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"token":   tokenString,
		"expires": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"user":    username,
	})
}
