package main

import (
	"os"

	"github.com/kartelmode/deltix-task/internal/configs"
	"github.com/kartelmode/deltix-task/internal/serializer"
	"github.com/kartelmode/deltix-task/internal/writer"
)

var currencies map[string]float64 = make(map[string]float64)
var users map[string]float64 = make(map[string]float64)
var usersTransactions map[string][]float64 = make(map[string][]float64)
var usersTransactionsTime map[string][]int64 = make(map[string][]int64)

func GetPrice(currency string, delta float64) float64 {
	if currency == "USD" {
		return delta
	}
	return currencies[currency+"USD"]
}

func main() {
	if len(os.Args) > 2 {
		panic("THERE ARE USELESS COMMAND LINE ARGUMENTS")
	}

	configs.UnmarshalOptions()
	flag := os.Args[len(os.Args)-1]
	options := configs.GetRunOption(flag)

	userData := serializer.UnmarshalUserData()
	market := serializer.UnmarshalMarket()

	for _, element := range userData {
		_, ok := users[element.UserId]
		if !ok {
			users[element.UserId] = 0
		}
	}
	writer := writer.MakeWriter()
	writer.CreateCsvIfNotExist(options.AnswerFilename)

	leftTime := userData[0].Timestamp
	pointerUser := 0
	pointerCurrencies := 0
	for leftTime < userData[len(userData)-1].Timestamp {
		for pointerUser < len(userData) && userData[pointerUser].Timestamp < leftTime+int64(options.DeltaTimestamp) {
			for pointerCurrencies < len(market) && market[pointerCurrencies].Timestamp <= userData[pointerUser].Timestamp {
				currencies[market[pointerCurrencies].Currency] = market[pointerCurrencies].Price
				pointerCurrencies++
			}
			id := userData[pointerUser].UserId
			usersTransactions[id] = append(usersTransactions[id], GetPrice(userData[pointerUser].Currency, userData[pointerUser].Delta))
			usersTransactionsTime[id] = append(usersTransactionsTime[id], userData[pointerUser].Timestamp)
			pointerUser++
		}
		for user := range users {
			minv := users[user]
			maxv := users[user]
			sum := float64(0)
			p := 0
			last := leftTime
			for _, transaction := range usersTransactions[user] {
				sum += users[user] * float64(usersTransactionsTime[user][p]-last)
				users[user] += transaction
				minv = min(minv, users[user])
				maxv = max(maxv, users[user])
				last = usersTransactionsTime[user][p]
				p++
			}
			sum += users[user] * float64(leftTime+int64(options.DeltaTimestamp)-last)
			sum /= float64(options.DeltaTimestamp)
			writer.Write(user, minv, maxv, sum, leftTime)
			usersTransactions[user] = make([]float64, 0)
			usersTransactionsTime[user] = make([]int64, 0)
		}
		leftTime += int64(options.DeltaTimestamp)
	}

	return

	// pricing := converter.MakePricing(market)
	// manager := manager.MakeManager(pricing, userData)
	// combiner := combiner.MakeCombiner(userData[0].Timestamp)
	// writer := writer.MakeWriter()
	// writer.CreateCsvIfNotExist(options.AnswerFilename)
	// manager.Run(options.DeltaTimestamp, options.TimestampsPerGoroutine, options.GoroutinesCount, combiner)
	// combiner.Combine(writer, int64(options.DeltaTimestamp))
	// writer.CloseFile()
}
