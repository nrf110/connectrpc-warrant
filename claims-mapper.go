package connectrpc_permit

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsMapper interface {
	Map(claims jwt.Claims) (*User, error)
}

type defaultClaimsMapper struct {
	ClaimsMapper
	customClaims map[string]string
}

func NewDefaultClaimsMapper(customClaims map[string]string) ClaimsMapper {
	return defaultClaimsMapper{
		customClaims: customClaims,
	}
}

func (mapper defaultClaimsMapper) Map(claims jwt.Claims) (*User, error) {
	subject, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	attributes := make(Attributes)
	switch t := claims.(type) {
	case jwt.MapClaims:
		for k, v := range t {
			if name, ok := mapper.customClaims[k]; ok {
				attributes[name] = v
			}
		}
	}

	return &User{
		Key:        subject,
		Attributes: attributes,
	}, nil
}
