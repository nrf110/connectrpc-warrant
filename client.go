package connectrpc_permit

import (
	"fmt"
	"github.com/permitio/permit-golang/pkg/enforcement"
)

type CheckClient interface {
	Check(user *User, config CheckConfig) (bool, error)
}

// PermitInterface subset of the actual PermitInterface.  Permit currently has a bug
// where the `AllTenantsCheck` signature does not match between the interface and implementation.
type PermitInterface interface {
	Check(user enforcement.User, action enforcement.Action, resource enforcement.Resource) (bool, error)
	BulkCheck(requests ...enforcement.CheckRequest) ([]bool, error)
	FilterObjects(user enforcement.User, action enforcement.Action, context map[string]string, resources ...enforcement.ResourceI) ([]enforcement.ResourceI, error)
	GetUserPermissions(user enforcement.User, tenants ...string) (enforcement.UserPermissions, error)
}

type permitCheckClient struct {
	Client PermitInterface
}

func NewCheckClient(client PermitInterface) CheckClient {
	return &permitCheckClient{
		Client: client,
	}
}

func (client *permitCheckClient) Check(user *User, config CheckConfig) (bool, error) {
	switch config.Type {
	case SINGLE:
		return client.check(user, config)
	case BULK:
		return client.bulkCheck(user, config)
	default:
		return false, fmt.Errorf("unexpected CheckType %s", config.Type)
	}
}

func (client *permitCheckClient) check(user *User, config CheckConfig) (bool, error) {
	request := config.Checks[0].toCheckRequest(user)
	return client.Client.Check(
		request.User,
		request.Action,
		request.Resource,
	)
}

func (client *permitCheckClient) bulkCheck(user *User, config CheckConfig) (bool, error) {
	var requests []enforcement.CheckRequest
	for _, check := range config.Checks {
		requests = append(requests, check.toCheckRequest(user))
	}

	results, err := client.Client.BulkCheck(requests...)
	if err != nil {
		return false, err
	}

	switch config.Mode {
	case ALL_OF:
		for _, result := range results {
			if !result {
				return false, nil
			}
		}
		return true, nil
	default:
		for _, result := range results {
			if result {
				return true, nil
			}
		}
		return false, nil
	}
}
