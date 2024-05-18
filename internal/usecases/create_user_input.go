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

func (uc *CreateUserInputUsecase) Create(ctx context.Context, params []CreateUserInputParams) ([]*entities.UserInput, error) {
	var userInputs []*entities.UserInput
	txFn := func(ctx context.Context) error {
		for _, param := range params {
			userInput := &entities.UserInput{
				UserID:    param.UserID,
				AssetID:   param.AssetID,
				Quantity:  param.Quantity,
				PaidPrice: param.PaidPrice,
			}

			userInput, err := uc.userInputRepo.Create(ctx, userInput)

			if err != nil {
				return err
			}

			userAsset, err := uc.userAssetRepo.FindBy(ctx, &entities.UserAsset{
				UserID:  param.UserID,
				AssetID: param.AssetID,
			})

			if err == database.ErrNotFound {
				_, err := uc.addUserAssetUc.Add(ctx, AddUserAssetParams{
					UserID:   param.UserID,
					AssetID:  param.AssetID,
					Quantity: param.Quantity,
				})

				if err != nil {
					return err
				}

				userInputs = append(userInputs, userInput)
				continue
			}

			if err != nil {
				return err
			}

			userAsset.Quantity += param.Quantity
			_, err = uc.userAssetRepo.Update(ctx, userAsset)

			if err != nil {
				return err
			}

			userInputs = append(userInputs, userInput)
		}

		return nil
	}

	if err := uc.userInputRepo.Tx(ctx, txFn); err != nil {
		return nil, err
	}

	return userInputs, nil
}
