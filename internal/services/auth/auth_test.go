package auth

import (
	"context"
	"io"
	"log/slog"
	"sso/internal/domain"
	"sso/internal/lib/config"
	"sso/internal/services/auth/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMaxWidth(t *testing.T) {

	TestingTable := []struct {
		name string
		args []string // аргументы
		want int      // ожидаемое значение
	}{
		{
			name: "sd",
		},
	}

	for _, tt := range TestingTable {

		t.Run(tt.name, func(t *testing.T) {

			userSaver := mocks.NewUserSaver(t)
			userProvider := mocks.NewUserProvider(t)
			appProvider := mocks.NewAppProvider(t)

			userProvider.
				On("IsAdmin", mock.Anything, int64(0)).
				Return(false, nil)

			service := New(
				domain.NewContext(
					slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})),
					config.MustLoadByPath("../../../config/local.yaml"),
				),
				NewStoreProvider(
					userSaver,
					userProvider,
					appProvider,
				),
			)

			isAdmin, err := service.IsAdmin(context.TODO(), 0)

			assert.Nil(t, err)

			assert.False(t, isAdmin)

		})
	}

}
