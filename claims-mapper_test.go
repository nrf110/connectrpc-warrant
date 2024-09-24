package connectrpc_permit

import (
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	var customClaims map[string]string

	BeforeEach(func() {
		customClaims = map[string]string{
			"foo": "Foo",
			"BAR": "Bar",
		}
	})

	ShouldUseSubjectAsKey := func(claims jwt.Claims) {
		It("should set the subject as the key", func(ctx SpecContext) {
			mapper := NewDefaultClaimsMapper(customClaims)
			user, err := mapper.Map(claims)
			Expect(err).To(BeNil())
			subject, err := claims.GetSubject()
			Expect(err).To(BeNil())
			Expect(user.Key).To(Equal(subject))
		})
	}

	When("claims are MapClaims", func() {
		ShouldUseSubjectAsKey(jwt.MapClaims{
			"sub": "1234",
		})

		Context("for each customClaim", func() {
			It("should add the claim to the attributes if it exists", func(ctx SpecContext) {
				mapper := NewDefaultClaimsMapper(customClaims)
				user, err := mapper.Map(jwt.MapClaims{
					"foo": "baz",
					"BAR": "bla",
				})
				Expect(err).To(BeNil())
				Expect(user.Attributes).To(HaveKeyWithValue("Foo", "baz"))
				Expect(user.Attributes).To(HaveKeyWithValue("Bar", "bla"))
			})

			It("should ignore the claim if it doesn't exist", func(ctx SpecContext) {
				mapper := NewDefaultClaimsMapper(customClaims)
				user, err := mapper.Map(jwt.MapClaims{
					"foo": "baz",
				})
				Expect(err).To(BeNil())
				Expect(user.Attributes).To(HaveKeyWithValue("Foo", "baz"))
				Expect(user.Attributes).ToNot(HaveKey("Bar"))
			})
		})

		Context("for any claims not in customClaims", func() {
			It("should not include the claim in the attributes", func(ctx SpecContext) {
				mapper := NewDefaultClaimsMapper(customClaims)
				user, err := mapper.Map(jwt.MapClaims{
					"foo": "baz",
					"BAR": "bla",
					"bat": "boo",
				})
				Expect(err).To(BeNil())
				Expect(user.Attributes).NotTo(HaveKey("bat"))
				Expect(user.Attributes).NotTo(HaveKey("Bat"))
			})
		})
	})

	When("claims are RegisteredClaims", func() {
		ShouldUseSubjectAsKey(jwt.RegisteredClaims{
			Subject: "1234",
		})
	})
})
