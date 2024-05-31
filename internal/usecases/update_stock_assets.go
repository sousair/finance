package usecases

import (
	"context"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/database"
	"github.com/sousair/go-finance/internal/infra/financeinquiry"
)

type UpdateStockAssetsUsecase struct {
	assetRepository      *database.Repository[entities.Asset]
	assetPriceRepository *database.Repository[entities.AssetPrice]
	finance              financeinquiry.Finance
}

func NewUpdateStockAssetsUsecase(
	assetRepository *database.Repository[entities.Asset],
	assetPriceRepo *database.Repository[entities.AssetPrice],
	finance financeinquiry.Finance,
) *UpdateStockAssetsUsecase {
	return &UpdateStockAssetsUsecase{
		assetRepository:      assetRepository,
		assetPriceRepository: assetPriceRepo,
		finance:              finance,
	}
}

func (uc UpdateStockAssetsUsecase) UpdateStockAssets(ctx context.Context) error {
	// TODO: Define a reasonable interval to update the stock assets
	assetQuery := &entities.Asset{
		Type: entities.AssetTypeStock,
	}

	stockAssets, err := uc.assetRepository.FindAll(ctx, assetQuery)
	if err != nil {
		if err == database.ErrNotFound {
			return nil
		}
		return err
	}

	for _, asset := range stockAssets {
		stockData, err := uc.finance.StockAssetDataInquiry(ctx, asset)
		if err != nil {
			// WARN: A return here is bad, maybe sent a warning log should be better
			return err
		}

		asset.Price = stockData.Price

		if _, err := uc.assetRepository.Update(ctx, asset); err != nil {
			// WARN: Same here
			return err
		}

		assetPrice := &entities.AssetPrice{
			AssetID: asset.ID,
			Price:   stockData.Price,
		}

		if _, err := uc.assetPriceRepository.Create(ctx, assetPrice); err != nil {
			// WARN: Same here
			return err
		}
	}

	return nil
}
