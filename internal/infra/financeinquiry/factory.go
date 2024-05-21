package financeinquiry

import (
	"context"
	"errors"

	scraper "github.com/sousair/go-brazilian-investments-scraper"
	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/database"
)

type (
	Finance interface {
		UpdateAssetPrice(assetID string) (float64, error)
	}

	BrazilianFinance struct {
		assetRepo *database.Repository[entities.Asset]
	}
)

var _ Finance = (*BrazilianFinance)(nil)

func NewBrazilianFinance() *BrazilianFinance {
	return &BrazilianFinance{}
}

// TODO: fix this to a better method and add more data on response and save
func (b BrazilianFinance) UpdateAssetPrice(assetID string) (float64, error) {
	asset, err := b.assetRepo.FindById(context.TODO(), assetID)
	if err != nil {
		return 0, err
	}

	// TODO: This would be a switch in the future
	if asset.Type == "STOCK" {
		price, err := getStockData(asset.Code)
		if err != nil {
			return 0, err
		}

		return price, nil
	}

	// FIX:
	return 0, nil
}

func getStockData(symbol string) (float64, error) {
	stockData := scraper.GetStockData(symbol)

	if stockData == nil {
		return 0, errors.New("stock data not found")
	}

	return stockData.Price, nil
}
