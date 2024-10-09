package authentication

import (
	"context"
	"log/slog"
	"net/http"
	"sso/internal/services/authentication/models"
	auth_store "sso/internal/store/sql/authentication"

	"github.com/go-playground/validator/v10"

	"github.com/pkg/errors"
)

func (srv server) register() http.HandlerFunc {

	log := srv.log.With(
		slog.String("op", "authentication.register"),
	)

	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			log.Error(err.Error())
		}

		serviceRequest := models.RegisterRequest{
			Nickname: r.Form.Get("nickname"),
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("pswd"),
		}

		if err := validator.New().Struct(serviceRequest); err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
			return
		}

		_, err = srv.auth.RegisterNewUser(
			context.Background(),
			serviceRequest,
		)

		if err != nil {
			if errors.Is(err, auth_store.ErrUserExists) {

				log.Error("registration failed", slog.String("err", err.Error()))

				http.Error(w, errors.Wrap(err, "registration failed").Error(), http.StatusNotAcceptable)
				return
			}

			log.Error("registration failed", slog.String("err", err.Error()))

			http.Error(w, errors.Wrap(err, "registration failed").Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Hello BITCH!"))
	}
}
