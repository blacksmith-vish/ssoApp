package validator

import (
	"errors"

	grpcValidator "github.com/blacksmith-vish/sso/pkg/validators/grpc_exchange"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Validate(x grpcValidator.ProtoExchange) error {

	if err := grpcValidator.Validate(x); err != nil {

		if errors.Is(err, grpcValidator.ErrNilRequest) {
			return status.Error(codes.Internal, err.Error())
		}

		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
