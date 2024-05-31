package financeinquiry

import (
	"context"
	"errors"

	scraper "github.com/sousair/go-brazilian-investments-scraper"
	"github.com/sousair/go-finance/internal/entities"
)

type (
	Finance interface {
		StockAssetDataInquiry(ctx context.Context, asset *entities.Asset) (*StockData, error)
	}

	StockData struct {
		Price float64
	}

	BrazilianFinance struct{}
)

var _ Finance = (*BrazilianFinance)(nil)

func NewBrazilianFinance() *BrazilianFinance {
	return &BrazilianFinance{}
}

func (b BrazilianFinance) StockAssetDataInquiry(ctx context.Context, asset *entities.Asset) (*StockData, error) {
	// NOTE: this should be a switch statement in the future
	if asset.Type == "STOCK" {
		price, err := getStockData(asset.Code)
		if err != nil {
			return nil, err
		}

		return &StockData{Price: price}, nil
	}

	return nil, nil
}

func getStockData(symbol string) (float64, error) {
	stockData := scraper.GetStockData(symbol)

	if stockData == nil {
		return 0, errors.New("stock data not found")
	}

	return stockData.Price, nil
}
