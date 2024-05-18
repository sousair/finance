package usecases

import (
	"context"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/database"
)

type (
	CreateUserInputUsecase struct {
		userInputRepo  *database.Repository[entities.UserInput]
		userAssetRepo  *database.Repository[entities.UserAsset]
		addUserAssetUc *AddUserAssetUsecase
	}

	CreateUserInputParams struct {
		UserID    string
		AssetID   string
		Quantity  int
		PaidPrice float64
	}
)

func NewCreateUserInput(
	userInputRepo *database.Repository[entities.UserInput],
	userAssetRepo *database.Repository[entities.UserAsset],
	addUserAssetUc *AddUserAssetUsecase,
) *CreateUserInputUsecase {
	return &CreateUserInputUsecase{
		userInputRepo:  userInputRepo,
		userAssetRepo:  userAssetRepo,
		addUserAssetUc: addUserAssetUc,
	}
}

func (uc *CreateUserInputUsecase) Create(ctx context.Context, params CreateUserInputParams) (*entities.UserInput, error) {
	userInput := &entities.UserInput{
		UserID:    params.UserID,
		AssetID:   params.AssetID,
		Quantity:  params.Quantity,
		PaidPrice: params.PaidPrice,
	}

	userInput, err := uc.userInputRepo.Create(ctx, userInput)

	if err != nil {
		return nil, err
	}

	userAsset, err := uc.userAssetRepo.FindBy(ctx, &entities.UserAsset{
		UserID:  params.UserID,
		AssetID: params.AssetID,
	})

	if err == database.ErrNotFound {
		_, err := uc.addUserAssetUc.Add(ctx, AddUserAssetParams{
			UserID:   params.UserID,
			AssetID:  params.AssetID,
			Quantity: params.Quantity,
		})

		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	userAsset.Quantity += params.Quantity
	_, err = uc.userAssetRepo.Update(ctx, userAsset)

	if err != nil {
		return nil, err
	}

	return userInput, nil
}
