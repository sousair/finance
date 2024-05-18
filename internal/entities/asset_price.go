package entities

import "github.com/sousair/go-finance/internal/infra/database"

type AssetPrice struct {
	database.BaseEntity
	AssetID string  `json:"asset_id"`
	Price   float64 `json:"price"`
}
