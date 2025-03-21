package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aygoko/EcoMind/backend/domain"
	"github.com/aygoko/EcoMind/usecases/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// UserHandler handles user-related HTTP endpoints
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: s,
	}
}

// RegisterRoutes registers user routes with Fiber
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")

	// Existing user routes
	userGroup := apiGroup.Group("/users")
	userGroup.Post("/", h.CreateUser)
	userGroup.Get("/:login", h.GetUserByLogin)

	// Authentication routes
	authGroup := apiGroup.Group("/auth")
	authGroup.Get("/google", h.GoogleAuthInit)
	authGroup.Get("/google/callback", h.GoogleAuthCallback)
	authGroup.Get("/tiktok", h.TikTokAuthInit)
	authGroup.Get("/tiktok/callback", h.TikTokAuthCallback)
}

// GoogleAuthInit initiates Google authentication flow
func (h *UserHandler) GoogleAuthInit(c *fiber.Ctx) error {
	url := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
	return c.Redirect(url, 302)
}

// GoogleAuthCallback handles Google authentication callback
func (h *UserHandler) GoogleAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing authorization code"})
	}

	// Exchange code for token
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to exchange authorization code"})
	}

	// Get user info from Google
	client := googleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve user info from Google"})
	}
	defer resp.Body.Close()

	// Parse Google user info
	var googleUser struct {
		Email string `json:"email"`
		Sub   string `json:"sub"` // Google's user ID
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid user info format"})
	}

	// Find or create user
	user, err := h.UserService.FindOrCreateUserByProvider(
		"google",
		googleUser.Sub,
		googleUser.Email,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT token
	tokenString, err := generateJWT(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

// TikTokAuthInit initiates TikTok authentication flow
func (h *UserHandler) TikTokAuthInit(c *fiber.Ctx) error {
	url := tiktokConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
	return c.Redirect(url, 302)
}

// TikTokAuthCallback handles TikTok authentication callback
func (h *UserHandler) TikTokAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing authorization code"})
	}

	// Exchange code for token (TikTok requires custom handling)
	token, err := tiktokConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to exchange authorization code"})
	}

	// Get user info from TikTok
	resp, err := http.Get("https://open-api.tiktok.com/user/info/?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve user info from TikTok"})
	}
	defer resp.Body.Close()

	// Parse TikTok user info
	var tiktokUser struct {
		User struct {
			UserID   string `json:"user_id"`
			NickName string `json:"nick_name"`
		} `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tiktokUser); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid user info format"})
	}

	// Find or create user
	user, err := h.UserService.FindOrCreateUserByProvider(
		"tiktok",
		tiktokUser.User.UserID,
		"", // TikTok doesn't provide email by default
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT token
	tokenString, err := generateJWT(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

// Existing methods remain unchanged
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdUser, err := h.UserService.Create(&user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(createdUser)
}

func (h *UserHandler) GetUserByLogin(c *fiber.Ctx) error {
	login := c.Params("login")
	user, err := h.UserService.Get(login)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

// JWT Generation Helper
func generateJWT(userID string) (string, error) {
	// JWT Secret (replace with a secure value in production)
	jwtSecret := []byte("your-secure-jwt-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSecret)
}

// OAuth2 Configurations (should be in separate config file)
var googleConfig = &oauth2.Config{
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
	RedirectURL:  "http://localhost:3000/api/auth/google/callback",
	Endpoint:     google.Endpoint,
	Scopes:       []string{"openid", "email", "profile"},
}

var tiktokConfig = &oauth2.Config{
	ClientID:     "YOUR_TIKTOK_CLIENT_ID",
	ClientSecret: "YOUR_TIKTOK_CLIENT_SECRET",
	RedirectURL:  "http://localhost:3000/api/auth/tiktok/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.tiktok.com/v2/oauth/authorize/",
		TokenURL: "https://open-api.tiktok.com/oauth/access_token/",
	},
	Scopes: []string{"user.info.basic"},
}
