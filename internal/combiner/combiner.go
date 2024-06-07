package combiner

import (
	"sync"

	"github.com/kartelmode/deltix-task/internal/models"
)

type Writer interface {
	Write(string, float64, float64, float64, int)
}

type Combiner struct {
	GroupedStatByStampId *map[int]*map[string]*models.Balance
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
	return combiner
}

func (combiner Combiner) StartListen(stampChan chan []*models.GroupedBalance) {
	for {
		stamps := <-stampChan
		for _, stamp := range stamps {
			// fmt.Println(stamp)
			stampId := stamp.Id
			(*combiner.GroupedStatByStampId)[stampId] = stamp.Balances
		}
	}
}

func GetStartTimeStamp(time, delta int) int {
	return (time / delta) * delta
}

func UpdateBalance(userId string, id int, balance *models.Balance, wg *sync.WaitGroup, writer Writer, delta int) {
	defer wg.Done()

	currentBalance := currentUserBalance[id]
	minBalance := currentBalance + balance.Min
	maxBalance := currentBalance + balance.Max
	averageBalance := currentBalance + balance.Sum/float64(delta)
	writer.Write(userId, minBalance, maxBalance, averageBalance, GetStartTimeStamp(balance.LastUpdate, delta))
	currentUserBalance[id] += balance.Current
}

func UpdateBalances(usersBalances *map[string]*models.Balance, userIds *map[string]int, writer Writer, delta int) {
	var wg sync.WaitGroup
	for userId, balance := range *usersBalances {
		wg.Add(1)
		id := (*userIds)[userId]
		go UpdateBalance(userId, id, balance, &wg, writer, delta)
	}
	wg.Wait()
}

func (combiner *Combiner) Combine(writer Writer, delta int) {
	userIds := map[string]int{}
	maxTimeStamp := 0
	for stampId, users := range *combiner.GroupedStatByStampId {
		maxTimeStamp = max(maxTimeStamp, stampId)
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
