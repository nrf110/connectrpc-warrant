package connectrpc_permit

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/permitio/permit-golang/pkg/enforcement"
)

var _ = Describe("Check", func() {
	const ACTION = "edit"

	var (
		stubUser     *User
		stubResource Resource
	)

	ShouldPropagateClientErrors := func() {
		It("should return errors from the client", func() {
			mock.SetUp(GinkgoT())
			permitClient := mock.Mock[PermitInterface]()
			mock.When(permitClient.Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())).
				ThenReturn(nil, fmt.Errorf("server error"))
		})
	}

	BeforeEach(func() {
		stubUser = &User{
			Key: "1234",
			Attributes: map[string]any{
				"foo": "bar",
			},
		}
		stubResource = Resource{
			Type:   "widget",
			Key:    "abcde",
			Tenant: "fghij",
			Attributes: map[string]any{
				"baz": "bap",
			},
		}
	})

	When("config type is single", func() {
		It("should invoke Check", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			permitClient := mock.Mock[PermitInterface]()
			mock.When(permitClient.Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())).
				ThenReturn(true, nil)

			checkClient := NewCheckClient(permitClient)
			result, err := checkClient.Check(stubUser, CheckConfig{
				Type: SINGLE,
				Checks: []Check{
					{
						Action:   ACTION,
						Resource: stubResource,
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeTrue())

			mock.Verify(permitClient, mock.Once()).Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())
			mock.Verify(permitClient, mock.Never()).BulkCheck(mock.Any[enforcement.CheckRequest]())
		})
	})

	When("config type is bulk", func() {
		Context("and the mode is all_of", func() {
			It("should invoke BulkCheck and return true when all results are true", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())

				permitClient := mock.Mock[PermitInterface]()
				mock.When(permitClient.BulkCheck(
					mock.Any[enforcement.CheckRequest](),
					mock.Any[enforcement.CheckRequest](),
				)).ThenReturn([]bool{true, true}, nil)

				checkClient := NewCheckClient(permitClient)
				result, err := checkClient.Check(stubUser, CheckConfig{
					Type: BULK,
					Mode: ALL_OF,
					Checks: []Check{
						{
							Action:   ACTION,
							Resource: stubResource,
						},
						{
							Action:   ACTION,
							Resource: stubResource,
						},
					},
				})

				Expect(err).To(BeNil())
				Expect(result).To(BeTrue())

				mock.Verify(permitClient, mock.Once()).BulkCheck(mock.Any[enforcement.CheckRequest](), mock.Any[enforcement.CheckRequest]())
				mock.Verify(permitClient, mock.Never()).Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())
			})

			It("should invoke BulkCheck and return false when at least one result is false", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())

				permitClient := mock.Mock[PermitInterface]()
				mock.When(permitClient.BulkCheck(
					mock.Any[enforcement.CheckRequest](),
					mock.Any[enforcement.CheckRequest](),
				)).ThenReturn([]bool{true, false, true}, nil)

				checkClient := NewCheckClient(permitClient)
				result, err := checkClient.Check(stubUser, CheckConfig{
					Type: BULK,
					Mode: ALL_OF,
					Checks: []Check{
						{
							Action:   ACTION,
							Resource: stubResource,
						},
						{
							Action:   ACTION,
							Resource: stubResource,
						},
					},
				})

				Expect(err).To(BeNil())
				Expect(result).To(BeFalse())

				mock.Verify(permitClient, mock.Once()).BulkCheck(mock.Any[enforcement.CheckRequest](), mock.Any[enforcement.CheckRequest]())
				mock.Verify(permitClient, mock.Never()).Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())
			})

			ShouldPropagateClientErrors()
		})

		Context("and the mode is any_of", func() {
			It("should invoke BulkCheck and return true when at least one result is true", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())

				permitClient := mock.Mock[PermitInterface]()
				mock.When(permitClient.BulkCheck(
					mock.Any[enforcement.CheckRequest](),
					mock.Any[enforcement.CheckRequest](),
				)).ThenReturn([]bool{false, false, true}, nil)

				checkClient := NewCheckClient(permitClient)
				result, err := checkClient.Check(stubUser, CheckConfig{
					Type: BULK,
					Mode: ANY_OF,
					Checks: []Check{
						{
							Action:   ACTION,
							Resource: stubResource,
						},
						{
							Action:   ACTION,
							Resource: stubResource,
						},
					},
				})

				Expect(err).To(BeNil())
				Expect(result).To(BeTrue())

				mock.Verify(permitClient, mock.Once()).BulkCheck(mock.Any[enforcement.CheckRequest](), mock.Any[enforcement.CheckRequest]())
				mock.Verify(permitClient, mock.Never()).Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())
			})

			It("should invoke BulkCheck and return false when all results are false", func(ctx SpecContext) {
				mock.SetUp(GinkgoT())

				permitClient := mock.Mock[PermitInterface]()
				mock.When(permitClient.BulkCheck(
					mock.Any[enforcement.CheckRequest](),
					mock.Any[enforcement.CheckRequest](),
				)).ThenReturn([]bool{false, false, false}, nil)

				checkClient := NewCheckClient(permitClient)
				result, err := checkClient.Check(stubUser, CheckConfig{
					Type: BULK,
					Mode: ANY_OF,
					Checks: []Check{
						{
							Action:   ACTION,
							Resource: stubResource,
						},
						{
							Action:   ACTION,
							Resource: stubResource,
						},
					},
				})

				Expect(err).To(BeNil())
				Expect(result).To(BeFalse())

				mock.Verify(permitClient, mock.Once()).BulkCheck(mock.Any[enforcement.CheckRequest](), mock.Any[enforcement.CheckRequest]())
				mock.Verify(permitClient, mock.Never()).Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())
			})

			ShouldPropagateClientErrors()
		})
	})

	When("config type is public", func() {
		It("should return an error", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			permitClient := mock.Mock[PermitInterface]()

			checkClient := NewCheckClient(permitClient)
			result, err := checkClient.Check(stubUser, CheckConfig{
				Type: PUBLIC,
				Checks: []Check{
					{
						Action:   ACTION,
						Resource: stubResource,
					},
				},
			})

			Expect(err.Error()).To(Equal("unexpected CheckType public"))
			Expect(result).To(BeFalse())

			mock.Verify(permitClient, mock.Never()).Check(mock.Any[enforcement.User](), mock.Any[enforcement.Action](), mock.Any[enforcement.Resource]())
			mock.Verify(permitClient, mock.Never()).BulkCheck(mock.Any[enforcement.CheckRequest]())
		})
	})
})
