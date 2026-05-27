package repositories

import (
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

type LeaveTypeRepository struct {
	*GenericRepository[models.LeaveType]
}

func NewLeaveTypeRepository(db *gorm.DB) *LeaveTypeRepository {
	return &LeaveTypeRepository{
		GenericRepository: NewGenericRepository(db, &models.LeaveType{}),
	}
}
