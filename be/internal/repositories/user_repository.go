package repositories

import (
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// UserRepository provides database operations for User model.
type UserRepository struct {
	*GenericRepository[models.User]
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		GenericRepository: NewGenericRepository(db, &models.User{}),
	}
}

// FindByEmail finds a user by email.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	users, err := r.FindByFieldMap(nil, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &users[0], nil
}
