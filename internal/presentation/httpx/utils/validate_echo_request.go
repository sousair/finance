package httpxutils

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ValidateEchoRequest[T any](e echo.Context, validator *validator.Validate) (*T, error) {
	var req T

	if err := e.Bind(&req); err != nil {
		return nil, errors.New("invalid request")
	}

	if err := validator.Struct(req); err != nil {
		return nil, errors.New("invalid request")
	}

	return &req, nil
}
