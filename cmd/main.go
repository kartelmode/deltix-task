package main

import (
	"os"

	"github.com/kartelmode/deltix-task/internal/combiner"
	"github.com/kartelmode/deltix-task/internal/configs"
	"github.com/kartelmode/deltix-task/internal/converter"
	"github.com/kartelmode/deltix-task/internal/manager"
	"github.com/kartelmode/deltix-task/internal/serializer"
	"github.com/kartelmode/deltix-task/internal/writer"
)

func main() {
	if len(os.Args) > 2 {
		panic("THERE ARE USELESS COMMAND LINE ARGUMENTS")
	}

	configs.UnmarshalOptions()
	flag := os.Args[len(os.Args)-1]
	options := configs.GetRunOption(flag)

	userData := serializer.UnmarshalUserData()
	market := serializer.UnmarshalMarket()

	pricing := converter.MakePricing(market)
	manager := manager.MakeManager(pricing, userData)
	combiner := combiner.MakeCombiner(userData[0].Timestamp)
	writer := writer.MakeWriter()
	writer.CreateCsvIfNotExist(options.AnswerFilename)
	manager.Run(options.DeltaTimestamp, options.TimestampsPerGoroutine, options.GoroutinesCount, combiner)
	combiner.Combine(writer, int64(options.DeltaTimestamp))
	writer.CloseFile()
}
