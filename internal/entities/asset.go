package entities

import "github.com/sousair/go-finance/internal/infra/database"

type AssetType string

const (
	AssetTypeStock AssetType = "STOCK"
)

type Asset struct {
	database.BaseEntity
	Type AssetType `json:"type"`
	// NOTE: Maybe add country here to change how to handle inquiry and etc.
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"value"`
	// NOTE: what about assets that don't have a code?
	// like a fixed income asset optional ?
	Code string `json:"code"`
}
