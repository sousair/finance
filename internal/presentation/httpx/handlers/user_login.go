package httpxhandlers

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	httpxutils "github.com/sousair/go-finance/internal/presentation/httpx/utils"
	"github.com/sousair/go-finance/internal/usecases"
)

type (
	UserLoginHandler struct {
		validator        *validator.Validate
		userLoginUsecase *usecases.UserLoginUsecase
	}

	UserLoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	UserLoginResponse struct {
		usecases.UserLoginUsecaseResponse
	}
)

var _ httpxutils.EchoHandler = (*UserLoginHandler)(nil)

func NewUserLoginHandler(validator *validator.Validate, userLoginUsecase *usecases.UserLoginUsecase) *UserLoginHandler {
	return &UserLoginHandler{
		validator:        validator,
		userLoginUsecase: userLoginUsecase,
	}
}

func (h UserLoginHandler) Handle(e echo.Context) error {
	req, err := httpxutils.ValidateEchoRequest[UserLoginRequest](e, h.validator)

	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusBadRequest, err.Error())
	}

	params := usecases.UserLoginUsecaseParams{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.userLoginUsecase.Login(e.Request().Context(), params)

	if err != nil {
		return err
	}

	return e.JSON(200, res)
}
