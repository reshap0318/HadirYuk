package repositories

import (
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

type UserProfileRepository struct {
	*GenericRepository[models.UserProfile]
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{
		GenericRepository: NewGenericRepository(db, &models.UserProfile{}),
	}
}
