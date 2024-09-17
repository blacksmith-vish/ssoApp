package auth

import (
	"context"
	"io"
	"log/slog"
	"sso/internal/domain"
	"sso/internal/lib/config"
	"sso/internal/services/auth/mocks"
	"sso/internal/store/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMaxWidth(t *testing.T) {

	TestingTable := []struct {
		name string
		arg  string // аргументы
		want struct {
			result bool
			err    bool
		} // ожидаемое значение
	}{
		{
			name: "test-1",
			arg:  "0",
			want: struct {
				result bool
				err    bool
			}{
				result: false,
				err:    false,
			},
		},
		{
			name: "test-2",
			arg:  "2",
			want: struct {
				result bool
				err    bool
			}{
				result: true,
				err:    false,
			},
		},
		{
			name: "test-3",
			arg:  "-1",
			want: struct {
				result bool
				err    bool
			}{
				result: false,
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
		Return(false, models.ErrUserNotFound)

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
