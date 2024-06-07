package models

type Transaction struct {
	Timestamp int64
	Price     float32
}

func MakeTransaction(time int64, price float32) *Transaction {
	return &Transaction{
		Timestamp: time,
		Price:     price,
	}
}
