package authentication

import (
	"context"
	"io"
	"log/slog"
	"sso/internal/services/authentication/mocks"
	serviceModels "sso/internal/services/authentication/models"
	"sso/internal/store/sqlite"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ConfigTest struct {
	tokenTTL time.Duration
}

func NewConfigTest() *ConfigTest {
	return &ConfigTest{
		tokenTTL: time.Minute,
	}
}

func (conf ConfigTest) GetTokenTTL() time.Duration {
	return conf.tokenTTL
}

func TestMaxWidth(t *testing.T) {

	TestingTable := []struct {
		name string
		arg  serviceModels.IsAdminRequest // аргументы
		want struct {
			result serviceModels.IsAdminResponse
			err    bool
		} // ожидаемое значение
	}{
		{
			name: "test-1",
			arg:  serviceModels.IsAdminRequest{UserID: "0"},
			want: struct {
				result serviceModels.IsAdminResponse
				err    bool
			}{
				result: serviceModels.IsAdminResponse{IsAdmin: false},
				err:    false,
			},
		},
		{
			name: "test-2",
			arg:  serviceModels.IsAdminRequest{UserID: "2"},
			want: struct {
				result serviceModels.IsAdminResponse
				err    bool
			}{
				result: serviceModels.IsAdminResponse{IsAdmin: true},
				err:    false,
			},
		},
		{
			name: "test-3",
			arg:  serviceModels.IsAdminRequest{UserID: "-1"},
			want: struct {
				result serviceModels.IsAdminResponse
				err    bool
			}{
				result: serviceModels.IsAdminResponse{IsAdmin: false},
				err:    true,
			},
		},
	}

	userSaver := mocks.NewUserSaver(t)
	userProvider := mocks.NewUserProvider(t)
	appProvider := mocks.NewAppProvider(t)

	userProvider.
		On("IsAdmin", mock.Anything, "0").
		Return(false, nil).
		On("IsAdmin", mock.Anything, "2").
		Return(true, nil).
		On("IsAdmin", mock.Anything, "-1").
		Return(false, sqlite.ErrUserNotFound)

	service := NewService(
		slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})),
		NewConfigTest(),
		userSaver,
		userProvider,
		appProvider,
	)

	for _, tt := range TestingTable {

		t.Run(tt.name, func(t *testing.T) {
			isAdmin, err := service.IsAdmin(context.TODO(), tt.arg)

			if tt.want.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.result, isAdmin)
		})
	}

}
