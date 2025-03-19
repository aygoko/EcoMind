package http

import (
	"net/http"

	repository "github.com/aygoko/EcoMInd/backend/domain"
	"github.com/aygoko/EcoMInd/usecases/service"
	"github.com/gofiber/fiber/v2" // [[2]][[3]][[7]]
)

// UserHandler handles user-related HTTP endpoints
type UserHandler struct {
	Service repository.UserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		Service: s, // Direct assignment if service implements the interface
	}
}

// RegisterRoutes registers user routes with Fiber
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	apiGroup := app.Group("/api/users") // Create route group [[6]][[9]]

	// POST /api/users
	apiGroup.Post("/", h.CreateUser)

	// GET /api/users/{login}
	apiGroup.Get("/:login", h.GetUserByLogin)
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user repository.User
	if err := c.BodyParser(&user); err != nil { // Fiber's built-in parser [[6]][[9]]
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdUser, err := h.Service.Save(&user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(createdUser) // Fiber's JSON response [[6]]
}

// GetUserByLogin retrieves user by login
func (h *UserHandler) GetUserByLogin(c *fiber.Ctx) error {
	login := c.Params("login") // Fiber's parameter extraction [[6]][[9]]
	user, err := h.Service.Get(login)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user) // Simplified response [[6]]
}
