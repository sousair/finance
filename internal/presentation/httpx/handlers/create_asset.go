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
	CreateAssetHandler struct {
		validator     *validator.Validate
		createAssetUc *usecases.CreateAssetUsecase
	}

	CreateAssetRequest struct {
		Type        string  `json:"type" validate:"required"`
		Name        string  `json:"name" validate:"required"`
		Description string  `json:"description" validate:"required"`
		Code        string  `json:"code" validate:"required"`
		Price       float64 `json:"price" validate:"required,gt=0"`
	}

	CreateAssetResponse struct {
		Asset *entities.Asset `json:"asset"`
	}
)

var _ httpxutils.EchoHandler = (*CreateAssetHandler)(nil)

func NewCreateAssetHandler(validator *validator.Validate, createAssetUc *usecases.CreateAssetUsecase) *CreateAssetHandler {
	return &CreateAssetHandler{
		validator:     validator,
		createAssetUc: createAssetUc,
	}
}

func (h CreateAssetHandler) Handle(e echo.Context) error {
	req, err := httpxutils.ValidateEchoRequest[CreateAssetRequest](e, h.validator)

	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusBadRequest, err.Error())
	}

	asset, err := h.createAssetUc.Create(e.Request().Context(), usecases.CreateAssetParams{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Price:       req.Price,
	})

	if err != nil {
		return httpxutils.NewHttpErrorResponse(e, http.StatusInternalServerError, "internal server error")
	}

	res := &CreateAssetResponse{
		Asset: asset,
	}

	return e.JSON(http.StatusCreated, res)
}
