package authentication

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sso/internal/services/authentication/models"

	"github.com/go-chi/chi/v5"
)

//go:embed static
var staticFiles embed.FS

type Authentication interface {
	RegisterNewUser(
		ctx context.Context,
		request models.RegisterRequest,
	) (response models.RegisterResponse, err error)
}

type authenticationAPI struct {
	log  *slog.Logger
	auth Authentication
}

type server = *authenticationAPI

func NewAuthenticationServer(
	log *slog.Logger,
	auth Authentication,
) *authenticationAPI {

	return &authenticationAPI{
		log:  log,
		auth: auth,
	}

}

func (srv server) InitRouters(router *chi.Mux) {

	fs := http.FileServer(http.FS(staticFiles))
	router.Handle("/static/*", fs)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	router.Post("/register", srv.register())

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		(w).Header().Set("Content-Type", "text;charset-utf-8")

		templ, err := template.ParseFS(staticFiles, "static/*.html")

		if err != nil {
			fmt.Println(err)
			return
		}

		templ.ExecuteTemplate(w, "index.html", nil)
	})

}
