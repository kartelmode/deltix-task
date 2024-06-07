package pools

import "github.com/kartelmode/deltix-task/internal/models"

func GetNewBalance() *map[string]*models.Balance {
	return new(map[string]*models.Balance)
}
