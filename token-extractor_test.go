package connectrpc_permit

import (
	"connectrpc.com/connect"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
	"net/http"
)

const (
	token    string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	badToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.jM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

var _ = Describe("Extract", func() {
	When("the authorization header is empty", func() {
		It("should return an unauthorized error", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			req := mock.Mock[connect.AnyRequest]()
			mock.When(req.Header()).ThenReturn(map[string][]string{})

			extractor := NewAuthorizationHeaderTokenExtractor()
			_, err := extractor.Extract(req)
			Expect(err.Error()).To(Equal("unauthenticated"))
		})
	})

	When("the authorization header is present", func() {
		Context("and the token is valid", func() {
			It("should return the claims", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())
				req := mock.Mock[connect.AnyRequest]()
				mock.When(req.Header()).ThenReturn(http.Header{
					"Authorization": {
						fmt.Sprintf("bearer %s", token),
					},
				})

				extractor := NewAuthorizationHeaderTokenExtractor()
				claims, err := extractor.Extract(req)
				Expect(err).To(BeNil())
				Expect(claims.GetSubject()).To(Equal("1234567890"))
			})
		})

		Context("and the token is invalid", func() {
			It("should return an unauthorized error", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())
				req := mock.Mock[connect.AnyRequest]()
				mock.When(req.Header()).ThenReturn(http.Header{
					"Authorization": {
						fmt.Sprintf("bearer %s", badToken),
					},
				})

				extractor := NewAuthorizationHeaderTokenExtractor()
				_, err := extractor.Extract(req)
				Expect(err).ToNot(BeNil())
			})
		})

		Context("and the token type is not bearer", func() {
			It("should return an unauthorized error", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())
				req := mock.Mock[connect.AnyRequest]()
				mock.When(req.Header()).ThenReturn(http.Header{
					"Authorization": {
						fmt.Sprintf("bearerer %s", token),
					},
				})

				extractor := NewAuthorizationHeaderTokenExtractor()
				_, err := extractor.Extract(req)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
