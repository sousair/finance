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
		UserInputs []UserInputRequest `json:"inputs" validate:"required,min=1"`
	}

	UserInputRequest struct {
		AssetID   string  `json:"asset_id" validate:"required"`
		Quantity  int     `json:"quantity" validate:"required,gt=0"`
		PaidPrice float64 `json:"paid_price" validate:"required,gt=0"`
	}

	CreateUserInputResponse struct {
		UserInput []*entities.UserInput `json:"inputs"`
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
	userId := e.Get("user_id").(string)

	req, err := httpxutils.ValidateEchoRequest[CreateUserInputRequest](e, h.validator)
	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusBadRequest, err.Error())
	}

	userInputs := make([]usecases.CreateUserInputParams, len(req.UserInputs))
	for i, ui := range req.UserInputs {
		userInputs[i] = usecases.CreateUserInputParams{
			UserID:    userId,
			AssetID:   ui.AssetID,
			Quantity:  ui.Quantity,
			PaidPrice: ui.PaidPrice,
		}
	}

	inputs, err := h.createUserInputUsecase.Create(e.Request().Context(), userInputs)
	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusInternalServerError, "internal server error")
	}

	res := &CreateUserInputResponse{
		UserInput: inputs,
	}

	return e.JSON(200, res)
}
