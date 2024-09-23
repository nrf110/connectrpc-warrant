package connectrpc_permit_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPermit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Permit Suite")
}
