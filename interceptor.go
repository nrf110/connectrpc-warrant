package connectrpc_permit

import (
	"connectrpc.com/connect"
	"context"
	"errors"
)

func NewPermitInterceptor(client CheckClient, tokenExtractor TokenExtractor, claimsMapper ClaimsMapper) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			checkable, ok := req.Any().(Checkable)
			if !ok {
				return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
			}
			checks := checkable.GetChecks()
			if !checks.IsPublic() {
				claims, err := tokenExtractor.Extract(req)
				if err != nil {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
				}

				user, err := claimsMapper.Map(claims)
				if err != nil {
					return nil, err
				}

				result, err := client.Check(user, checks)
				if err != nil {
					return nil, err
				}
				if !result {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
				}
			}
			return next(ctx, req)
		})
	}
}
