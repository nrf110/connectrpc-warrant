package connectrpc_permit

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/permitio/permit-golang/pkg/enforcement"
)

var _ = Describe("toCheckRequest", func() {
	When("attributes are present", func() {
		It("should include them on the CheckRequest", func(ctx SpecContext) {
			user := User{
				Key: "abcde",
				Attributes: map[string]any{
					"foo": "bar",
				},
			}

			check := Check{
				Action: "edit",
				Resource: Resource{
					Type:   "Widget",
					Key:    "1234",
					Tenant: "fghi",
					Attributes: map[string]any{
						"baz": "blah",
					},
				},
			}
			req := check.toCheckRequest(&user)

			Expect(req).To(BeEquivalentTo(enforcement.CheckRequest{
				User: enforcement.User{
					Key: "abcde",
					Attributes: map[string]any{
						"foo": "bar",
					},
				},
				Action: "edit",
				Resource: enforcement.Resource{
					Type:   "Widget",
					Key:    "1234",
					Tenant: "fghi",
					Attributes: map[string]any{
						"baz": "blah",
					},
				},
			}))
		})
	})

	When("resource key is empty", func() {
		It("should set the key to *", func(ctx SpecContext) {
			user := User{
				Key: "abcde",
				Attributes: map[string]any{
					"foo": "bar",
				},
			}

			check := Check{
				Action: "edit",
				Resource: Resource{
					Type:   "Widget",
					Key:    "",
					Tenant: "fghi",
					Attributes: map[string]any{
						"baz": "blah",
					},
				},
			}
			req := check.toCheckRequest(&user)

			Expect(req).To(BeEquivalentTo(enforcement.CheckRequest{
				User: enforcement.User{
					Key: "abcde",
					Attributes: map[string]any{
						"foo": "bar",
					},
				},
				Action: "edit",
				Resource: enforcement.Resource{
					Type:   "Widget",
					Key:    "*",
					Tenant: "fghi",
					Attributes: map[string]any{
						"baz": "blah",
					},
				},
			}))
		})
	})
})
