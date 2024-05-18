package httpxhandlers

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sousair/go-finance/internal/entities"
	httpxutils "github.com/sousair/go-finance/internal/presentation/httpx/utils"
	"github.com/sousair/go-finance/internal/usecases"
)

type (
	CreateUserInputHandler struct {
		validator              *validator.Validate
		createUserInputUsecase *usecases.CreateUserInputUsecase
	}

	CreateUserInputRequest struct {
		// TODO: Get user ID from JWT token
		UserID     string             `json:"user_id" validate:"required"`
		UserInputs []UserInputRequest `json:"inputs" validate:"required,min=1"`
	}

	UserInputRequest struct {
		AssetID   string  `json:"asset_id"`
		Quantity  int     `json:"quantity"`
		PaidPrice float64 `json:"paid_price"`
	}

	CreateUserInputResponse struct {
		UserInput *entities.UserInput `json:"user_input"`
	}
)

var _ httpxutils.EchoHandler = (*CreateUserInputHandler)(nil)

func NewCreateUserInputHandler(validator *validator.Validate, createUserInputUsecase *usecases.CreateUserInputUsecase) *CreateUserInputHandler {
	return &CreateUserInputHandler{
		validator:              validator,
		createUserInputUsecase: createUserInputUsecase,
	}
}

func (h CreateUserInputHandler) Handle(e echo.Context) error {
	// TODO: yeah
	_, err := httpxutils.ValidateEchoRequest[CreateUserInputRequest](e, h.validator)

	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusBadRequest, err.Error())
	}

	return nil
}
