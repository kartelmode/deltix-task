package configs

import (
	"encoding/json"
	"os"

	"github.com/kartelmode/deltix-task/internal/models"
)

var runOptions map[string]models.Options

func UnmarshalOptions() {
	jsonData, err := os.ReadFile("./configs/options.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(jsonData, &runOptions); err != nil {
		panic(err)
	}
}

func GetRunOption(cmdFlag string) models.Options {
	option, ok := runOptions[cmdFlag]
	if !ok {
		panic("THERE IS NO SUCH FLAG")
	}
	return option
}
