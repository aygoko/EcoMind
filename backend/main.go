package main

import (
    "context"
    "database/sql"
    "flag"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/go-redis/redis/v8"
    _ "github.com/lib/pq" // PostgreSQL driver
)

// domain.User struct (replace with your actual model)
type User struct {
    ID          string  `json:"id"`
    Login       string  `json:"login"`
    Email       string  `json:"email"`
    PhoneNumber string  `json:"phone_number"`
    Password    string  `json:"-"`
    CO2         float64 `json:"co2"`
}

// UserService interface
type UserService interface {
    Get(login string) (*User, error)
}

// UserRepository struct
type UserRepository struct {
    DB          *sql.DB
    RedisClient *redis.Client
    Logger      *log.Logger
}

// NewUserRepository initializes the repository
func NewUserRepository(db *sql.DB, redisClient *redis.Client, logger *log.Logger) *UserRepository {
    return &UserRepository{
        DB:          db,
        RedisClient: redisClient,
        Logger:      logger,
    }
}

// Get retrieves a user with Redis cache
func (r *UserRepository) Get(login string) (*User, error) {
    key := "user:login:" + login
    userJSON, err := r.RedisClient.Get(context.Background(), key).Result()
    if err == nil {
        var user User
        if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
            r.Logger.Printf("Failed to unmarshal cached user: %v", err)
            return nil, err
        }
        return &user, nil
    } else if err != redis.Nil {
        return nil, err
    }

    // Fetch from PostgreSQL
    row := r.DB.QueryRow("SELECT id, login, email, phone_number, co2 FROM users WHERE login = $1", login)
    var user User
    err = row.Scan(&user.ID, &user.Login, &user.Email, &user.PhoneNumber, &user.CO2)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

    // Cache the result
    if err := r.cacheUser(&user); err != nil {
        r.Logger.Printf("Failed to cache user: %v", err)
    }

    return &user, nil
}

// cacheUser stores the user in Redis
func (r *UserRepository) cacheUser(user *User) error {
    key := "user:login:" + user.Login
    userJSON, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return r.RedisClient.Set(context.Background(), key, userJSON, 5*time.Minute).Err()
}

// UserService implementation
type UserServiceImpl struct {
    Repo *UserRepository
}

// NewUserServiceImpl creates a new user service
func NewUserServiceImpl(repo *UserRepository) UserService {
    return &UserServiceImpl{Repo: repo}
}

func (s *UserServiceImpl) Get(login string) (*User, error) {
    return s.Repo.Get(login)
}

// UserHandler struct
type UserHandler struct {
    Service UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service UserService) *UserHandler {
    return &UserHandler{Service: service}
}

// RegisterRoutes registers user-related routes
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
    api := app.Group("/api/users")
    api.Get("/:login", h.GetUser)
}

// GetUser retrieves a user by login
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    login := c.Params("login")
    user, err := h.Service.Get(login)
    if err != nil {
        if err.Error() == "user not found" {
            return c.Status(404).JSON(fiber.Map{"error": "User not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
    }
    return c.JSON(user)
}

// Form handlers
type PageData struct {
    AgreementChecked bool
    Message          string
}

func showForm(c *fiber.Ctx) error {
    data := PageData{
        AgreementChecked: false,
        Message:          c.Query("message", ""),
    }
    return c.RenderString(`
        <!DOCTYPE html>
        <html>
        <head>
            <title>Пользовательское соглашение</title>
        </head>
        <body>
            <h1>Пользовательское соглашение</h1>
            <form action="/submit" method="POST">
                <p>Прочитайте соглашение и поставьте галочку:</p>
                <label>
                    <input type="checkbox" name="agreement" value="agree">
                    Я согласен с пользовательским соглашением
                </label><br><br>
                <button type="submit">Отправить</button>
            </form>
            {{if .Message}}
                <p>{{.Message}}</p>
            {{end}}
        </body>
        </html>`, data, "form")
}

func handleFormSubmission(c *fiber.Ctx) error {
    type Form struct {
        Agreement string `form:"agreement"`
    }

    var form Form
    if err := c.BodyParser(&form); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Ошибка при разборе формы")
    }

    if form.Agreement == "agree" {
        return c.Redirect("/", http.StatusFound)
    }

    return c.Redirect("/?message=Вы%20не%20согласились%20с%20соглашением", http.StatusSeeOther)
}

func main() {
    addr := flag.String("addr", ":8080", "HTTP server address")
    flag.Parse()

    // Initialize PostgreSQL
    pgDB, err := sql.Open("postgres", "user=youruser password=yourpass dbname=yourdb sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }
    defer pgDB.Close()

    // Initialize Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    // Initialize repository and service
    repo := NewUserRepository(pgDB, redisClient, log.Default())
    userService := NewUserServiceImpl(repo)
    userHandler := NewUserHandler(userService)

    // Initialize Fiber app
    app := fiber.New(fiber.Config{
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    })

    // CORS configuration
    app.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     "GET,POST,PUT,PATCH,DELETE",
        AllowHeaders:     "Content-Type,Authorization",
        AllowCredentials: true,
    }))

    // Middlewares
    app.Use(recover.New())
    app.Use(func(c *fiber.Ctx) error {
        log.Printf("%s %s %s", c.IP(), c.Method(), c.Path())
        return c.Next()
    })

    // Register routes
    userHandler.RegisterRoutes(app)
    app.Get("/", showForm)
    app.Post("/submit", handleFormSubmission)

    log.Printf("Server listening on %s", *addr)
    if err := app.Listen(*addr); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}