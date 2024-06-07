package converter

import (
	"github.com/kartelmode/deltix-task/internal/manager"
	"github.com/kartelmode/deltix-task/internal/models"
)

type Pricing struct {
	Pointer    int
	Currencies map[string]float64
	Market     []*models.Market
}

func MakePricing(market []*models.Market) *Pricing {
	return &Pricing{
		Pointer:    0,
		Currencies: map[string]float64{},
		Market:     market,
	}
}

func (pricing Pricing) CopyRestPricing() manager.Pricing {
	newMarket := make([]*models.Market, len(pricing.Market)-pricing.Pointer)
	for i := 0; i < len(newMarket); i++ {
		newMarket[i] = models.CopyMarket(pricing.Market[pricing.Pointer+i])
	}
	newCurrencies := map[string]float64{}
	for key, value := range pricing.Currencies {
		newCurrencies[key] = value
	}
	return &Pricing{
		Pointer:    0,
		Currencies: newCurrencies,
		Market:     newMarket,
	}
}

func (pricing Pricing) ChangePrice(currency string, price float64) {
	pricing.Currencies[currency] = price
}

func (pricing Pricing) Update(time int) {
	for pricing.Pointer < len(pricing.Market) && pricing.Market[pricing.Pointer].Timestamp <= time {
		curMarket := pricing.Market[pricing.Pointer]
		pricing.ChangePrice(curMarket.Currency, curMarket.Price)
		pricing.Pointer++
	}
}

func (pricing Pricing) ConvertToUSD(currency string, price float64) float64 {
	if currency == "USD" {
		return price
	}
	currency += "USD"
	return price * pricing.Currencies[currency]
}
