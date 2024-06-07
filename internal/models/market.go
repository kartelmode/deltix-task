package models

type Market struct {
	Currency  string  `csv:"symbol"`
	Timestamp int     `csv:"timestamp"`
	Price     float64 `csv:"price"`
}

func (market Market) GetTime() int {
	return market.Timestamp
}

func (market Market) GetPrice() float64 {
	return market.Price
}

func CopyMarket(market *Market) *Market {
	return &Market{
		Currency:  market.Currency,
		Timestamp: market.Timestamp,
		Price:     market.Price,
	}
}
