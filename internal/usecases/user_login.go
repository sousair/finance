package usecases

import (
	"context"
	"errors"

	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/cipher"
	"github.com/sousair/go-finance/internal/infra/database"
	"github.com/sousair/go-finance/internal/infra/token"
)

type (
	UserLoginUsecase struct {
		userRepo *database.Repository[entities.User]
		cipher   cipher.Cipher
		token    token.Token[UserTokenPayload]
	}

	UserLoginUsecaseParams struct {
		Email    string
		Password string
	}

	UserLoginUsecaseResponse struct {
		Token   string            `json:"token"`
		Payload *UserTokenPayload `json:"payload"`
	}

	UserTokenPayload struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
	}
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewUserLoginUsecase(userRepo *database.Repository[entities.User], cipher cipher.Cipher, token token.Token[UserTokenPayload]) *UserLoginUsecase {
	return &UserLoginUsecase{
		userRepo: userRepo,
		cipher:   cipher,
		token:    token,
	}
}

func (u UserLoginUsecase) Login(ctx context.Context, params UserLoginUsecaseParams) (*UserLoginUsecaseResponse, error) {
	user, err := u.userRepo.FindBy(ctx, &entities.User{Email: params.Email})

	if err != nil {
		return nil, err
	}

	if err := u.cipher.Compare(user.Password, params.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	payload := &UserTokenPayload{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	token, err := u.token.Generate(payload)

	if err != nil {
		return nil, err
	}

	return &UserLoginUsecaseResponse{
		Token:   token,
		Payload: payload,
	}, nil
}
