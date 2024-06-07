package models

type UserData struct {
	UserId    string  `csv:"user_id"`
	Currency  string  `csv:"currency"`
	Timestamp int     `csv:"timestamp"`
	Delta     float64 `csv:"delta"`
}
