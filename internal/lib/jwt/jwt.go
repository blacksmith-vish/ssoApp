package jwt

import (
	"sso/internal/store/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(
	user models.User,
	app models.App,
	duration time.Duration,
) (string, error) {
	// token := jwt.New(jwt.SigningMethodHS256)

	// claims := token.Claims.(jwt.MapClaims)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": user.ID,
			"email":  user.Email,
			"exp":    time.Now().Add(duration).Unix(),
			"appID":  app.ID,
		})

	// claims["userID"] = user.ID
	// claims["email"] = user.Email
	// claims["exp"] = time.Now().Add(duration).Unix()
	// claims["appID"] = app.ID

	tokeString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokeString, nil

}
