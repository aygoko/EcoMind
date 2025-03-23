package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

	pkgHttp "github.com/aygoko/EcoMInd/backend/api/types/user"
	repository "github.com/aygoko/EcoMInd/backend/repository/ram_storage"
	"github.com/aygoko/EcoMInd/backend/usecases/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	userHandler := pkgHttp.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{"*"},
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	userHandler.WithObjectHandlers(r)

	log.Printf("Starting HTTP server on %s", *addr)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func showForm(w http.ResponseWriter, r *http.Request) {
	const formHTML = `
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
    </html>
    `
	tmpl, err := template.New("form").Parse(formHTML)
	if err != nil {
		http.Error(w, "Ошибка при парсинге шаблона", http.StatusInternalServerError)
		return
	}

	data := PageData{
		AgreementChecked: false,
		Message:          "",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка при рендеринге шаблона", http.StatusInternalServerError)
	}
}

func handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
		return
	}

	agreement := r.FormValue("agreement") == "agree"

	var message string
	if agreement {
		message = "Спасибо! Вы согласились с пользовательским соглашением."
	} else {
		message = "Вы не согласились с пользовательским соглашением. Пожалуйста, поставьте галочку."
	}

	http.Redirect(w, r, fmt.Sprintf("/?message=%s", message), http.StatusSeeOther)
}
