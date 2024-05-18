package httpxhandlers

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sousair/go-finance/internal/entities"
	httpxuitls "github.com/sousair/go-finance/internal/presentation/httpx/utils"
	httpxutils "github.com/sousair/go-finance/internal/presentation/httpx/utils"
	"github.com/sousair/go-finance/internal/usecases"
)

type (
	CreateUserHandler struct {
		validator    *validator.Validate
		createUserUc *usecases.CreateUserUsecase
	}

	CreateUserRequest struct {
		Name              string `json:"name" validate:"required"`
		Email             string `json:"email" validate:"required,email"`
		PlainTextPassword string `json:"password" validate:"required,min=8"`
	}

	CreateUserResponse struct {
		User *entities.User `json:"user"`
	}
)

var _ httpxuitls.EchoHandler = (*CreateUserHandler)(nil)

func NewCreateUserHandler(validator *validator.Validate, createUserUc *usecases.CreateUserUsecase) *CreateUserHandler {
	return &CreateUserHandler{
		validator:    validator,
		createUserUc: createUserUc,
	}
}

func (h CreateUserHandler) Handle(e echo.Context) error {
	req, err := httpxutils.ValidateEchoRequest[CreateUserRequest](e, h.validator)

	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusBadRequest, err.Error())
	}

	user, err := h.createUserUc.Create(e.Request().Context(), usecases.CreateUserParams{
		Name:              req.Name,
		Email:             req.Email,
		PlainTextPassword: req.PlainTextPassword,
	})

	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusInternalServerError, "internal server error")
	}

	res := &CreateUserResponse{
		User: user,
	}

	return e.JSON(http.StatusCreated, res)
}
