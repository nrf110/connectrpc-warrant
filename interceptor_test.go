package connectrpc_permit

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
)

type stubCheckable struct {
	Checkable
	checks CheckConfig
}

func (r *stubCheckable) GetChecks() CheckConfig {
	return r.checks
}

var _ = Describe("Authorizing a ConnectRPC request", func() {
	When("the request is Checkable", func() {
		It("should invoke the next handler when the check call returns true", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(true, nil)

			claims := mock.Mock[jwt.Claims]()
			mock.When(claims.GetSubject()).ThenReturn("abcde", nil)

			extractor := mock.Mock[TokenExtractor]()
			mock.When(extractor.Extract(mock.Any[connect.AnyRequest]())).ThenReturn(claims, nil)

			claimsMapper := mock.Mock[ClaimsMapper]()
			mock.When(claimsMapper.Map(claims)).ThenReturn(&User{Key: "abcde"}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, claimsMapper)
			result, err := interceptor(next)(ctx, req)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(res))
		})

		It("should invoke the next handler when the CheckConfig is public", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(true, nil)

			extractor := mock.Mock[TokenExtractor]()

			claimsMapper := mock.Mock[ClaimsMapper]()

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{
				Type: PUBLIC,
			}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, claimsMapper)
			result, err := interceptor(next)(ctx, req)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(res))
		})

		It("should return a permission denied error when the check returns false", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, nil)

			claims := mock.Mock[jwt.Claims]()
			mock.When(claims.GetSubject()).ThenReturn("abcde", nil)

			extractor := mock.Mock[TokenExtractor]()
			mock.When(extractor.Extract(mock.Any[connect.AnyRequest]())).ThenReturn(claims, nil)

			claimsMapper := mock.Mock[ClaimsMapper]()
			mock.When(claimsMapper.Map(claims)).ThenReturn(&User{Key: "abcde"}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, claimsMapper)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})

		It("should return a permission denied error when the request is unauthenticated", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, nil)

			extractor := mock.Mock[TokenExtractor]()
			mock.When(extractor.Extract(mock.Any[connect.AnyRequest]())).ThenReturn(nil, fmt.Errorf("unauthenticated"))

			claimsMapper := mock.Mock[ClaimsMapper]()
			mock.When(claimsMapper.Map(mock.Any[jwt.Claims]())).ThenReturn(&User{Key: "abcde"}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, claimsMapper)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})

		It("should return the error when the check call fails", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(false, fmt.Errorf("unknown error"))

			claims := mock.Mock[jwt.Claims]()
			mock.When(claims.GetSubject()).ThenReturn("abcde", nil)

			extractor := mock.Mock[TokenExtractor]()
			mock.When(extractor.Extract(mock.Any[connect.AnyRequest]())).ThenReturn(claims, nil)

			claimsMapper := mock.Mock[ClaimsMapper]()
			mock.When(claimsMapper.Map(claims)).ThenReturn(&User{Key: "abcde"}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn(&stubCheckable{checks: CheckConfig{}})
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, claimsMapper)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("unknown error"))
			Expect(result).To(BeNil())
		})
	})

	When("the request is not Checkable", func() {
		It("should return a permission denied error", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			client := mock.Mock[CheckClient]()
			mock.When(client.Check(mock.Any[*User](), mock.Any[CheckConfig]())).ThenReturn(true, nil)

			claims := mock.Mock[jwt.Claims]()
			mock.When(claims.GetSubject()).ThenReturn("abcde", nil)

			extractor := mock.Mock[TokenExtractor]()
			mock.When(extractor.Extract(mock.Any[connect.AnyRequest]())).ThenReturn(claims, nil)

			claimsMapper := mock.Mock[ClaimsMapper]()
			mock.When(claimsMapper.Map(claims)).ThenReturn(&User{Key: "abcde"}, nil)

			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Any()).ThenReturn("")
			res := mock.Mock[connect.AnyResponse]()
			next := connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
				return res, nil
			})
			interceptor := NewPermitInterceptor(client, extractor, claimsMapper)
			result, err := interceptor(next)(ctx, req)
			Expect(err.Error()).To(Equal("permission_denied: permission denied"))
			Expect(result).To(BeNil())
		})
	})
})
