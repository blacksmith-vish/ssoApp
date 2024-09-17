package tests

import (
	"sso/internal/tests/suite"
	"testing"
	"time"

	"github.com/blacksmith-vish/sso/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = "iota"
	appID      = "iotafull"

	appSecret = "secret"

	passDefautlLen = 10
)

func TestRegisterLogin_Login_HappyPass(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()

	password := randomPassword()

	responseRegister, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)

	assert.NotEmpty(t, responseRegister.GetUserId())

	responseLogin, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})

	require.NoError(t, err)

	loginTime := time.Now()

	token := responseLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})

	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, responseRegister.GetUserId(), claims["userID"].(string))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, claims["appID"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Conf.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
