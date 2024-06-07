package serializer

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/kartelmode/deltix-task/internal/models"
)

func UnmarshalMarket() []*models.Market {
	in, err := os.Open("market_data.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	market := []*models.Market{}

	if err := gocsv.UnmarshalFile(in, &market); err != nil {
		panic(err)
	}
	return market
}

func UnmarshalUserData() []*models.UserData {
	in, err := os.Open("user_data.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	userData := []*models.UserData{}

	if err := gocsv.UnmarshalFile(in, &userData); err != nil {
		panic(err)
	}
	return userData
}
