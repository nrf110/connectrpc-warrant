package connectrpc_permit

import (
	"fmt"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/enforcement"
	"github.com/permitio/permit-golang/pkg/permit"
)

type CheckClient interface {
	Check(user *User, config CheckConfig) (bool, error)
}

type permitCheckClient struct {
	Client *permit.Client
}

func NewPermitCheckClient(cfg config.PermitConfig) CheckClient {
	return &permitCheckClient{
		Client: permit.NewPermit(cfg),
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
