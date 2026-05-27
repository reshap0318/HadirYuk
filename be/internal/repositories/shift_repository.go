package repositories

import (
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

type ShiftRepository struct {
	*GenericRepository[models.Shift]
}

func NewShiftRepository(db *gorm.DB) *ShiftRepository {
	return &ShiftRepository{
		GenericRepository: NewGenericRepository(db, &models.Shift{}),
	}
}
