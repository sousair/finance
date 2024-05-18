package httpxutils

import "github.com/labstack/echo/v4"

type EchoHandler interface {
	Handle(ctx echo.Context) error
}
