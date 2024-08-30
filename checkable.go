package connectrpc_workos

import (
	"github.com/workos/workos-go/v4/pkg/fga"
)

type Checkable interface {
	GetChecks() []fga.WarrantCheck
}
