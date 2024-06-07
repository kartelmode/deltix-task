package models

type GroupedBalance struct {
	Id       int
	Balances *map[string]*Balance
}

func MakeNewGroupBalance(balances *map[string]*Balance) *GroupedBalance {
	(*balances) = make(map[string]*Balance)
	return &GroupedBalance{
		Balances: balances,
	}
}

func (group *GroupedBalance) SetId(id int) {
	group.Id = id
}

func (group *GroupedBalance) GetBalance(userId string) bool {
	_, ok := (*group.Balances)[userId]
	return ok
}

func (group *GroupedBalance) SetBalance(userId string, startTime int) {
	(*group.Balances)[userId] = MakeBalance(startTime)
}

func (group *GroupedBalance) UpdateBalance(userId string, value float64, time int) {
	(*group.Balances)[userId].Update(value, time)
}

func (group *GroupedBalance) UpdateLastAll(lastTime int) {
	for _, balance := range *group.Balances {
		balance.UpdateLast(lastTime)
	}
}
