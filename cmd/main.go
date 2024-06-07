package main

import (
	"fmt"

	"github.com/kartelmode/deltix-task/internal/combiner"
	"github.com/kartelmode/deltix-task/internal/converter"
	"github.com/kartelmode/deltix-task/internal/manager"
	"github.com/kartelmode/deltix-task/internal/serializer"
	"github.com/kartelmode/deltix-task/internal/writer"
)

func main() {
	userData := serializer.UnmarshalUserData()
	fmt.Println(len(userData))
	market := serializer.UnmarshalMarket()
	fmt.Println(len(market))

	pricing := converter.MakePricing(market)
	manager := manager.MakeManager(pricing, userData)
	combiner := combiner.MakeCombiner()
	writer := writer.MakeWriter()
	writer.CreateCsvIfNotExist("output.csv")
	manager.Run(3600, 200, combiner)
	combiner.Combine(writer, 3600)
	writer.CloseFile()
	fmt.Println("Done")
}
