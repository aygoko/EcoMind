package main

import (
    "flag"
    "fmt"
    "log"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
    "github.com/gofiber/template/html/v3"

    "github.com/aygoko/EcoMInd/backend/api/types/user"
    repository "github.com/aygoko/EcoMInd/backend/repository/ram_storage"
    "github.com/aygoko/EcoMInd/backend/usecases/service"
)

type PageData struct {
    AgreementChecked bool
    Message          string
}

func main() {
    addr := flag.String("addr", ":8080", "HTTP server address")
    flag.Parse()

    userRepo := repository.NewUserRepository()
    userService := service.NewUserService(userRepo)
    userHandler := user.NewUserHandler(userService) // Adjust import path if needed

    // Initialize Fiber app
    app := fiber.New(fiber.Config{
        Views: html.New("./views", ".html"), // Configure template engine
    })

    // CORS Middleware
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT",
        AllowHeaders: "Content-Type,Authorization",
    }))

    // Logging & Recovery Middleware
    app.Use(func(c *fiber.Ctx) error {
        log.Printf("%s %s %s", c.IP(), c.Method(), c.Path())
        return c.Next()
    })
    app.Use(func(c *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered from panic: %v", r)
                c.Status(500).Send("Internal Server Error")
            }
        }()
        return c.Next()
    })

    
    userHandler.WithObjectHandlers(app) // Assuming your handler can use Fiber's router

    
    app.Get("/", showForm)
    app.Post("/submit", handleFormSubmission)

    log.Printf("Starting HTTP server on %s", *addr)
    err := app.Listen(*addr)
    if err != nil {
        log.Fatalf("Server failed: %v", err)
    }
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
    form := new(struct {
        Agreement string `form:"agreement"`
    })

    err := c.BodyParser(form)
    if err != nil {
        return c.Status(400).Send("Ошибка при разборе формы")
    }

    agreement := form.Agreement == "agree"

    var message string
    if agreement {
        message = "Спасибо! Вы согласились с пользовательским соглашением."
    } else {
        message = "Вы не согласились с пользовательским соглашением. Пожалуйста, поставьте галочку."
    }

    
    return c.Redirect(fmt.Sprintf("/?message=%s", message), 303)
}
