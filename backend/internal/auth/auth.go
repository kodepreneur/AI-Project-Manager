package auth

import (
	"strings"
	"time"

	"ai-project-manager/internal/config"
	"ai-project-manager/internal/database"
	"ai-project-manager/internal/logs"
	"ai-project-manager/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

type SetupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (s *AuthService) SetupStatus(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_setup": count > 0,
		},
	})
}

func (s *AuthService) Setup(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "ALREADY_SETUP", "message": "Admin already configured"},
		})
	}

	var req SetupRequest
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "Username, email, and password are required"},
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "HASH_ERROR", "message": "Could not hash password"},
		})
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "admin",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "CREATE_USER_FAILED", "message": err.Error()},
		})
	}

	token, _ := s.generateToken(&user)
	logs.Log.Info(nil, "Authentication", "Admin user setup completed successfully: "+user.Username)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil || req.UsernameOrEmail == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_CREDENTIALS", "message": "Credentials required"},
		})
	}

	var user models.User
	err := database.DB.Where("username = ? OR email = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_CREDENTIALS", "message": "Invalid username or password"},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_CREDENTIALS", "message": "Invalid username or password"},
		})
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "TOKEN_ERROR", "message": "Failed to create access token"},
		})
	}

	logs.Log.Info(nil, "Authentication", "User logged in: "+user.Username)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "UNAUTHORIZED", "message": "Not authenticated"},
		})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "User not found"},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (s *AuthService) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "INVALID_INPUT", "message": "New password is required"},
		})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "NOT_FOUND", "message": "User not found"},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{"code": "WRONG_PASSWORD", "message": "Current password is incorrect"},
		})
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(newHash)
	database.DB.Save(&user)

	logs.Log.Info(nil, "Authentication", "Password changed for user: "+user.Username)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password updated successfully",
	})
}

func (s *AuthService) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenStr := ""

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else if qToken := c.Query("token"); qToken != "" {
			tokenStr = qToken // For WebSocket connections
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{"code": "UNAUTHORIZED", "message": "Missing authentication token"},
			})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{"code": "UNAUTHORIZED", "message": "Invalid or expired token"},
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{"code": "UNAUTHORIZED", "message": "Invalid token claims"},
			})
		}

		c.Locals("userID", uint(claims["user_id"].(float64)))
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
