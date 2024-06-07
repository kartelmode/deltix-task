package models

type Transaction struct {
	Timestamp int
	Price     float32
}

func MakeTransaction(time int, price float32) *Transaction {
	return &Transaction{
		Timestamp: time,
		Price:     price,
	}
}
