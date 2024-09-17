package authentication

import (
	"context"
	"io"
	"log/slog"
	"sso/internal/domain"
	"sso/internal/lib/config"
	"sso/internal/services/authentication/mocks"
	serviceModels "sso/internal/services/authentication/models"
	storeModels "sso/internal/store/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
		Return(false, storeModels.ErrUserNotFound)

	service := New(
		domain.NewContext(
			slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})),
			config.MustLoadByPath("../../../config/local.yaml"),
		),
		NewAuthenticationStoreProvider(
			userSaver,
			userProvider,
			appProvider,
		),
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
