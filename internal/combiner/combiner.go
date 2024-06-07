package combiner

import (
	"sync"

	"github.com/kartelmode/deltix-task/internal/models"
)

type Writer interface {
	Write(string, float64, float64, float64, int64)
}

type Combiner struct {
	GroupedStatByStampId *map[int]*map[string]*models.Balance
	Mutex                *sync.Mutex
}

var currentUserBalance []float64

func SetEmptyBalances() {
	for i := 0; i < len(currentUserBalance); i++ {
		currentUserBalance[i] = 0
	}
}

func MakeCombiner() *Combiner {
	combiner := &Combiner{}
	combiner.GroupedStatByStampId = new(map[int]*map[string]*models.Balance)
	*combiner.GroupedStatByStampId = make(map[int]*map[string]*models.Balance)
	combiner.Mutex = new(sync.Mutex)
	return combiner
}

func (combiner Combiner) AddBalances(stamps []*models.GroupedBalance, wg *sync.WaitGroup) {
	defer wg.Done()
	combiner.Mutex.Lock()
	for _, stamp := range stamps {
		stampId := stamp.Id
		(*combiner.GroupedStatByStampId)[int(stampId)] = stamp.Balances
	}
	combiner.Mutex.Unlock()
}

func GetStartTimeStamp(time, delta int64) int64 {
	return (time / delta) * delta
}

func UpdateBalance(userId string, id int, balance *models.Balance, wg *sync.WaitGroup, writer Writer, delta int64) {
	defer wg.Done()

	currentBalance := currentUserBalance[id]
	minBalance := currentBalance + balance.Min
	maxBalance := currentBalance + balance.Max
	averageBalance := currentBalance + balance.Sum/float64(delta)
	writer.Write(userId, minBalance, maxBalance, averageBalance, GetStartTimeStamp(balance.LastUpdate, delta))
	currentUserBalance[id] += balance.Current
}

func UpdateBalances(usersBalances *map[string]*models.Balance, userIds *map[string]int, writer Writer, delta int64) {
	var wg sync.WaitGroup
	for userId, balance := range *usersBalances {
		wg.Add(1)
		id := (*userIds)[userId]
		go UpdateBalance(userId, id, balance, &wg, writer, delta)
	}
	wg.Wait()
}

func (combiner *Combiner) Combine(writer Writer, delta int64) {
	userIds := map[string]int{}
	maxTimeStamp := 0
	for stampId, users := range *combiner.GroupedStatByStampId {
		maxTimeStamp = max(maxTimeStamp, int(stampId))
		for userId := range *users {
			_, ok := userIds[userId]
			if !ok {
				userIds[userId] = len(userIds)
			}
		}
	}
	currentUserBalance = make([]float64, len(userIds))
	SetEmptyBalances()
	//map in go is a hash table :D
	for stampId := 0; stampId <= maxTimeStamp; stampId++ {
		for key, usersBalances := range *combiner.GroupedStatByStampId {
			if key != stampId {
				continue
			}
			UpdateBalances(usersBalances, &userIds, writer, delta)
		}
	}
}
