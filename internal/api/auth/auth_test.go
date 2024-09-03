package auth

import (
	"context"
	"sso/internal/api/auth/mocks"
	errs "sso/internal/domain/errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func TestIsAdmin(t *testing.T) {

	// mocks.NewAuthApi()

	ctx := context.Background()
	authApiMock := new(mocks.AuthApi)

	authApiMock.On("IsAdmin", ctx, gofakeit.Int64()).Return(false, errs.ErrUserNotFound)
	authApiMock.AssertExpectations(t)

	// if err != nil {
	// 	if errors.Is(err, errs.ErrUserNotFound) {
	// 		return nil, status.Error(codes.AlreadyExists, "login failed")
	// 	}
	// 	return nil, status.Error(codes.Internal, "login failed")
	// }
}

// func TestMaxWidth(t *testing.T) {

// 	// len(string 1) = 8
// 	// len(string 22) = 9
// 	// len(string 333) = 10

// 	var (
// 		testingTable = []CountNotNilElemsTest{
// 			{name: "41", args: slice41, want: 8},
// 			{name: "42", args: slice42, want: 8},
// 			{name: "43", args: slice43, want: 10},
// 			{name: "44", args: slice44, want: 8},
// 			{name: "45", args: slice45, want: 8},
// 			{name: "46", args: slice46, want: 10},
// 			{name: "47", args: slice47, want: 10},
// 			{name: "48", args: slice48, want: 10},
// 			{name: "49", args: slice49, want: 10},
// 			{name: "50", args: slice50, want: 8},
// 			{name: "51", args: slice51, want: 8},
// 			{name: "52", args: slice52, want: 10},
// 			{name: "53", args: slice53, want: 8},
// 			{name: "54", args: slice54, want: 0},
// 			{name: "55", args: slice55, want: 10},
// 			{name: "56", args: slice56, want: 10},
// 			{name: "57", args: slice57, want: 10},
// 			{name: "58", args: slice58, want: 10},
// 			{name: "59", args: slice59, want: 10},
// 			{name: "60", args: slice60, want: 10},
// 			{name: "61", args: slice61, want: 10},
// 			{name: "62", args: slice62, want: 10},
// 			{name: "63", args: slice63, want: 10},
// 			{name: "64", args: slice64, want: 10},
// 			{name: "65", args: slice65, want: 10},
// 			{name: "66", args: slice66, want: 10},
// 			{name: "67", args: slice67, want: 10},
// 			{name: "28", args: slice28, want: 0},
// 			{name: "29", args: slice29, want: 0},
// 			{name: "30", args: slice30, want: 8},
// 			{name: "31", args: slice31, want: 0},
// 			{name: "32", args: slice32, want: 0},
// 			{name: "33", args: slice33, want: 8},
// 			{name: "34", args: slice34, want: 8},
// 			{name: "35", args: slice35, want: 8},
// 			{name: "36", args: slice36, want: 9},
// 			{name: "37", args: slice37, want: 8},
// 			{name: "38", args: slice38, want: 0},
// 			{name: "39", args: slice39, want: 0},
// 			{name: "nil Slice", args: slice40, want: 0},
// 		}
// 	)

// 	for _, tt := range testingTable {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := maxWidth(tt.args)
// 			if got != tt.want {
// 				t.Errorf("MaxWidth() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
