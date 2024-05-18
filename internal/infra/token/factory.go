package token

import (
	"errors"
	"fmt"

	jwt_go "github.com/dgrijalva/jwt-go"
)

type Token[T any] interface {
	Generate(payload *T) (string, error)
	Validate(token string) (*T, error)
}

type GenerateTokenResponse[T any] struct {
	Token   string `json:"token"`
	Payload *T     `json:"payload"`
}

type JWT[T any] struct {
	secret []byte
}

var _ Token[any] = (*JWT[any])(nil)

func NewJWT[T any](secret string) *JWT[T] {
	return &JWT[T]{secret: []byte(secret)}
}

func (j JWT[T]) Generate(payload *T) (string, error) {
	claims := jwt_go.MapClaims{
		"payload": payload,
	}

	token, err := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, claims).SignedString(j.secret)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (j JWT[T]) Validate(token string) (*T, error) {
	claims := &jwt_go.MapClaims{}

	parsedToken, err := jwt_go.ParseWithClaims(token, claims, func(token *jwt_go.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	// TODO: This doesn't work
	payload, ok := (*claims)["payload"]

	if !ok {
		return nil, errors.New("invalid token")
	}

	fmt.Println(payload)

	tokenPayload, ok := payload.(*T)

	if !ok {
		return nil, errors.New("invalid token")
	}

	return tokenPayload, nil
}
