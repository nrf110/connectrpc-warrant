package connectrpc_permit

import (
	"connectrpc.com/connect"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenExtractor interface {
	Extract(req connect.AnyRequest) (jwt.Claims, error)
}

type authorizationHeaderTokenExtractor struct {
	parser *jwt.Parser
}

func NewAuthorizationHeaderTokenExtractor() TokenExtractor {
	return authorizationHeaderTokenExtractor{
		parser: jwt.NewParser(),
	}
}

func (extractor authorizationHeaderTokenExtractor) Extract(req connect.AnyRequest) (jwt.Claims, error) {
	header := req.Header().Get("Authorization")
	if header == "" {
		return nil, fmt.Errorf("unauthenticated")
	}
	tokenString := header[7:]
	token, _, err := extractor.parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}
