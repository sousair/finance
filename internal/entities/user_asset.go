package entities

import (
	"github.com/sousair/go-finance/internal/infra/database"
)

type UserAsset struct {
	database.BaseEntity
	UserID   string `json:"user_id"`
	AssetID  string `json:"asset_id"`
	Quantity int    `json:"quantity"`
}
