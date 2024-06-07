package manager

import (
	"sync"

	"github.com/kartelmode/deltix-task/internal/models"
	"github.com/kartelmode/deltix-task/internal/pools"
)

type Pricing interface {
	Update(int64)
	ConvertToUSD(string, float64) float64
	CopyRestPricing() Pricing
}

type Combiner interface {
	AddBalances([]*models.GroupedBalance, *sync.WaitGroup)
}

type Manager struct {
	Converter Pricing
	Data      []*models.UserData
}

var shift int64 = 0

func MakeManager(converter Pricing, data []*models.UserData) *Manager {
	shift = data[0].Timestamp
	return &Manager{
		Converter: converter,
		Data:      data,
	}
}

func GetStampId(row *models.UserData, delta int64) int64 {
	return (row.Timestamp - shift) / delta
}

func Calculate(left, right int,
	wg *WaitGroupCount,
	wgCombine *sync.WaitGroup,
	userData []*models.UserData,
	pricing Pricing,
	timeStampCount, delta int,
	combiner Combiner,
	shiftTimeStampId, shiftTime int64) {
	defer wg.Done()
	balances := make([]*models.GroupedBalance, timeStampCount)
	startTimeStamp := GetStampId(userData[left], int64(delta))
	for i := 0; i < timeStampCount; i++ {
		balances[i] = models.MakeNewGroupBalance(pools.GetNewBalance())
		balances[i].SetId(int64(i) + startTimeStamp - shiftTimeStampId)
	}

	for id := left; id < right; id++ {
		row := userData[id]
		timeStampId := GetStampId(row, int64(delta)) - startTimeStamp
		ok := balances[timeStampId].GetBalance(row.UserId)
		if !ok {
			balances[timeStampId].SetBalance(row.UserId, userData[left].Timestamp)
		}
		pricing.Update(row.Timestamp)
		balances[timeStampId].UpdateBalance(row.UserId, pricing.ConvertToUSD(row.Currency, row.Delta), row.Timestamp)
	}

	for id, stampBalances := range balances {
		timeStampId := int64(id) + startTimeStamp - shiftTimeStampId
		stampBalances.UpdateLastAll((timeStampId+1)*int64(delta) - 1 + shiftTime)
	}

	wgCombine.Add(1)
	go combiner.AddBalances(balances, wgCombine)
}

func GetNextStamp(left, delta int, data []*models.UserData, timeStampCount int) int {
	l := left
	r := min(len(data)-1, l+500*delta)
	var mid int = 0
	for l < r {
		mid = (l + r + 1) / 2
		if GetStampId(data[mid], int64(delta)) > GetStampId(data[left], int64(delta))+int64(timeStampCount)-1 {
			r = mid - 1
		} else {
			l = mid
		}
	}
	return r + 1
}

func (manager *Manager) Run(delta, timeStampCount int, combiner Combiner) {
	wg := &WaitGroupCount{}
	wgCombine := &sync.WaitGroup{}
	pointer := 0

	stampChan := make(chan []*models.GroupedBalance)
	defer close(stampChan)

	shiftStampId := GetStampId(manager.Data[0], int64(delta))
	shiftTime := manager.Data[0].Timestamp
	manager.Converter.Update(manager.Data[0].Timestamp - 1)
	for pointer < len(manager.Data) {
		for wg.GetCount() == 20 {
		}
		right := GetNextStamp(pointer, delta, manager.Data, timeStampCount)
		wg.Add(1)
		restPricing := manager.Converter.CopyRestPricing()

		go Calculate(pointer, right, wg, wgCombine, manager.Data, restPricing, timeStampCount, delta, combiner, shiftStampId, shiftTime)

		pointer = right
		if pointer < len(manager.Data) {
			manager.Converter.Update(manager.Data[pointer].Timestamp)
		}
	}

	wg.Wait()
	wgCombine.Wait()
}
