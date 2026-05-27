package repositories

import (
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

type OfficeLocationRepository struct {
	*GenericRepository[models.OfficeLocation]
}

func NewOfficeLocationRepository(db *gorm.DB) *OfficeLocationRepository {
	return &OfficeLocationRepository{
		GenericRepository: NewGenericRepository(db, &models.OfficeLocation{}),
	}
}
