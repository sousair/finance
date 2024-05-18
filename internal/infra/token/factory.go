package token

import (
	"errors"

	jwt_go "github.com/dgrijalva/jwt-go"
)

type (
	Token[T any] interface {
		Generate(payload *T) (*GenerateTokenResponse[T], error)
		Validate(token string) (*T, error)
	}

	GenerateTokenResponse[T any] struct {
		Token   string `json:"token"`
		Payload *T     `json:"payload"`
	}

	TokenPayload[T any] struct {
		jwt_go.StandardClaims
		Payload *T `json:"payload"`
	}
)

type JWT[T any] struct {
	secret []byte
}

var _ Token[any] = (*JWT[any])(nil)

func NewJWT[T any](secret string) *JWT[T] {
	return &JWT[T]{secret: []byte(secret)}
}

func (j JWT[T]) Generate(payload *T) (*GenerateTokenResponse[T], error) {
	claims := TokenPayload[T]{
		Payload: payload,
	}

	token, err := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, claims).SignedString(j.secret)

	if err != nil {
		return nil, err
	}

	res := &GenerateTokenResponse[T]{
		Token:   token,
		Payload: payload,
	}

	return res, nil
}

func (j JWT[T]) Validate(token string) (*T, error) {
	claims := &TokenPayload[T]{}

	parsedToken, err := jwt_go.ParseWithClaims(token, claims, func(token *jwt_go.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims.Payload, nil
}
