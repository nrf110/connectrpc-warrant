package connectrpc_permit

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
	. "github.com/wiremock/wiremock-testcontainers-go"
)

var _ = Describe("Check", func() {
	var container *WireMockContainer

	BeforeEach(func(ctx SpecContext) {
		var err error
		container, err = RunContainer(ctx, WithImage("wiremock/wiremock:3.9.1"))
		Expect(err).To(BeNil())
	})

	AfterEach(func(ctx SpecContext) {
		err := container.Terminate(ctx)
		Expect(err).To(BeNil())
	})

	When("config type is single", func() {
		It("should do things", func(ctx SpecContext) {
			mock.SetUp(GinkgoT())
			Expect(true).To(BeTrue())
		})
	})
})
