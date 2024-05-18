package usecases

import (
	"context"
	"errors"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/cipher"
	"github.com/sousair/go-finance/internal/infra/database"
)

type (
	CreateUserUsecase struct {
		userRepo *database.Repository[entities.User]
		cipher   cipher.Cipher
	}

	CreateUserParams struct {
		Name              string
		Email             string
		PlainTextPassword string
	}
)

var EmailAlreadyExistsError = errors.New("email already exists")

func NewCreateUserUsecase(userRepo *database.Repository[entities.User], cipher cipher.Cipher) *CreateUserUsecase {
	return &CreateUserUsecase{
		userRepo: userRepo,
		cipher:   cipher,
	}
}

func (uc CreateUserUsecase) Create(ctx context.Context, params CreateUserParams) (*entities.User, error) {
	emailAlreadyExists, err := uc.userRepo.FindBy(ctx, &entities.User{Email: params.Email})

	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	if emailAlreadyExists != nil {
		return nil, EmailAlreadyExistsError
	}

	hashedPassword, err := uc.cipher.Hash(params.PlainTextPassword)

	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: hashedPassword,
	}

	return uc.userRepo.Create(ctx, user)
}
