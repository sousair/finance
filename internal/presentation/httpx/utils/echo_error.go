package httpxutils

import "github.com/labstack/echo/v4"

type HttpErrorResponse struct {
	Message string `json:"message"`
}

func NewHttpErrorResponse(e echo.Context, httpStatus int, message string) error {
	httpError := &HttpErrorResponse{
		Message: message,
	}

	return e.JSON(httpStatus, httpError)
}
