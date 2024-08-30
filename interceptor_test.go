package connectrpc_workos

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/workos/workos-go/v4/pkg/fga"
)

type stubCheckable struct {
	checks []fga.WarrantCheck
}

func (r *stubCheckable) GetChecks() []fga.WarrantCheck {
	return r.checks
}

var _ = Describe("Authorizing a ConnectRPC request", func() {
	When("the request is Checkable", func() {
		It("should invoke the next handler when the check call returns authorized", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[context.Context](), mock.Any[fga.CheckOpts]())).ThenReturn(fga.CheckResponse{
				Result:     fga.CheckResultAuthorized,
				IsImplicit: false,
			}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: []fga.WarrantCheck{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewWarrantInterceptor(ctx, client)
			result, err := interceptor(next)(ctx, req)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(res))
		})

		It("should return a permission denied error when the check returns unauthorized", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[context.Context](), mock.Any[fga.CheckOpts]())).ThenReturn(fga.CheckResponse{
				Result:     fga.CheckResultNotAuthorized,
				IsImplicit: false,
			}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: []fga.WarrantCheck{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewWarrantInterceptor(ctx, client)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})

		It("should return the error when the check call fails", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[context.Context](), mock.Any[fga.CheckOpts]())).ThenReturn(nil, fmt.Errorf("unknown error"))

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: []fga.WarrantCheck{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewWarrantInterceptor(ctx, client)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("unknown error"))
			Expect(result).To(BeNil())
		})
	})

	When("the request is not Checkable", func() {
		It("should return a permission denied error", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[context.Context](), mock.Any[fga.CheckOpts]())).ThenReturn(fga.CheckResponse{
				Result:     fga.CheckResultNotAuthorized,
				IsImplicit: false,
			}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn("")
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewWarrantInterceptor(ctx, client)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})
	})
})
