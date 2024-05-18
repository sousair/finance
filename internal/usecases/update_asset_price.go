package usecases

import (
	"context"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/database"
)

type (
	UpdateAssetPriceUsecase struct {
		assetRepo      *database.Repository[entities.Asset]
		assetPriceRepo *database.Repository[entities.AssetPrice]
	}

	UpdateAssetPriceParams struct {
		AssetID string
		Price   float64
	}
)

func NewUpdateAssetPrice(
	assetRepo *database.Repository[entities.Asset],
	assetPriceRepo *database.Repository[entities.AssetPrice],
) *UpdateAssetPriceUsecase {
	return &UpdateAssetPriceUsecase{
		assetRepo:      assetRepo,
		assetPriceRepo: assetPriceRepo,
	}
}

func (uc UpdateAssetPriceUsecase) Update(ctx context.Context, params UpdateAssetPriceParams) (*entities.AssetPrice, error) {
	asset, err := uc.assetRepo.FindById(ctx, params.AssetID)

	if err != nil {
		return nil, err
	}

	asset.Price = params.Price
	asset, err = uc.assetRepo.Update(ctx, asset)

	if err != nil {
		return nil, err
	}

	assetPrice := &entities.AssetPrice{
		AssetID: params.AssetID,
		Price:   params.Price,
	}

	assetPrice, err = uc.assetPriceRepo.Create(ctx, assetPrice)

	if err != nil {
		return nil, err
	}

	return assetPrice, nil
}
