package connectrpc_workos

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"github.com/workos/workos-go/v4/pkg/fga"
)

type CheckClient interface {
	Check(ctx context.Context, opts fga.CheckOpts) (fga.CheckResponse, error)
}

func NewWarrantInterceptor(ctx context.Context, client CheckClient) connect.UnaryInterceptorFunc {
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

			result, err := client.Check(ctx, fga.CheckOpts{
				Op:     fga.CheckOpAllOf, // TODO make this configurable
				Checks: checks,
			})
			if err != nil {
				return nil, err
			}
			if !result.Authorized() {
				return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
			}
			return next(ctx, req)
		})
	}
}
