package usecases

import (
	"context"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/database"
)

type (
	CreateAssetUsecase struct {
		assetRepo *database.Repository[entities.Asset]
	}

	CreateAssetParams struct {
		Type        entities.AssetType
		Name        string
		Description string
		Code        string
		Price       float64
	}
)

func NewCreateAssetUsecase(assetRepo *database.Repository[entities.Asset]) *CreateAssetUsecase {
	return &CreateAssetUsecase{assetRepo: assetRepo}
}

func (uc CreateAssetUsecase) Create(ctx context.Context, params CreateAssetParams) (*entities.Asset, error) {
	// NOTE: Maybe inquiry to google finance API to check if asset exists (STOCK, CRYPTO, FII, etc)

	// TODO: Get asset price here if stock price is searchable
	asset := &entities.Asset{
		Type:        params.Type,
		Name:        params.Name,
		Description: params.Description,
		Code:        params.Code,
		Price:       params.Price,
	}

	return uc.assetRepo.Create(ctx, asset)
}
