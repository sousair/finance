package financeinquiry

type (
	Finance interface {
		// TODO: Check this method signature maybe change it to
		// return a better and typed response not just a float64
		UpdateAssetPrice(assetID string) (float64, error)
	}
	GoogleFinance struct {
		apiKey string
	}
)

var _ Finance = (*GoogleFinance)(nil)

func NewGoogleFinance(apiKey string) *GoogleFinance {
	return &GoogleFinance{apiKey: apiKey}
}

// TODO:: Define a response for this to update AssetPrice entity
func (gf GoogleFinance) UpdateAssetPrice(assetID string) (float64, error) {
	return 100, nil
}
