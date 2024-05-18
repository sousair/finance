package httpxmiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sousair/go-finance/internal/infra/token"
	httpxutils "github.com/sousair/go-finance/internal/presentation/httpx/utils"
	"github.com/sousair/go-finance/internal/usecases"
)

type UserAuthMiddleware struct {
	token token.Token[usecases.UserTokenPayload]
}

func NewUserAuthMiddleware(token token.Token[usecases.UserTokenPayload]) UserAuthMiddleware {
	return UserAuthMiddleware{token: token}
}

func (um UserAuthMiddleware) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		authHeader := e.Request().Header.Get("Authorization")

		if authHeader == "" {
			return httpxutils.NewHttpErrorResponse(e, http.StatusUnauthorized, "missing Authorization header")
		}

		authHeaderParts := strings.Split(authHeader, " ")

		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			return httpxutils.NewHttpErrorResponse(e, http.StatusUnauthorized, "invalid Authorization header")
		}

		token := authHeaderParts[1]

		if token == "" {
			return httpxutils.NewHttpErrorResponse(e, http.StatusUnauthorized, "missing token")
		}

		userPayload, err := um.token.Validate(token)

		if err != nil {
			return httpxutils.NewHttpErrorResponse(e, http.StatusUnauthorized, "invalid token")
		}

		e.Set("user_id", userPayload.UserID)
		e.Set("user_email", userPayload.Email)
		e.Set("user_name", userPayload.Name)

		return nil
	}
}
