package usecases

import (
	"context"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/database"
)

type (
	AddUserAssetUsecase struct {
		userAssetRepo *database.Repository[entities.UserAsset]
	}

	AddUserAssetParams struct {
		UserID   string
		AssetID  string
		Quantity int
	}
)

func NewAddUserAssetUsecase(userAssetRepo *database.Repository[entities.UserAsset]) *AddUserAssetUsecase {
	return &AddUserAssetUsecase{userAssetRepo: userAssetRepo}
}

func (uc AddUserAssetUsecase) Add(ctx context.Context, params AddUserAssetParams) (*entities.UserAsset, error) {
	userAsset := &entities.UserAsset{
		UserID:   params.UserID,
		AssetID:  params.AssetID,
		Quantity: params.Quantity,
	}

	return uc.userAssetRepo.Create(ctx, userAsset)
}
