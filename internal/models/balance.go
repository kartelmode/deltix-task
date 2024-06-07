package models

type Balance struct {
	Current    float64
	Min        float64
	Max        float64
	Sum        float64
	LastUpdate int64
}

func MakeBalance(startTime int64) *Balance {
	return &Balance{
		Current:    0,
		Min:        0,
		Max:        0,
		Sum:        0,
		LastUpdate: startTime,
	}
}

func (balance *Balance) Update(value float64, time int64) {
	balance.Sum += balance.Current * float64(time-balance.LastUpdate)
	balance.Current += value
	balance.Min = min(balance.Min, balance.Current)
	balance.Max = max(balance.Max, balance.Current)
	balance.LastUpdate = time
}

func (balance *Balance) UpdateLast(delta int64) {
	balance.Sum += balance.Current * float64(delta-balance.LastUpdate)
}
