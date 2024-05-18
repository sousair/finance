package entities

import "github.com/sousair/go-finance/internal/infra/database"

type UserInput struct {
	database.BaseEntity
	UserID    string  `json:"user_id"`
	AssetID   string  `json:"asset_id"`
	Quantity  int     `json:"quantity"`
	PaidPrice float64 `json:"paid_price"`
}
