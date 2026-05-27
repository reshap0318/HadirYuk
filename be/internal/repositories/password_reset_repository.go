package repositories

import (
	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/models"
)

// PasswordResetRepository provides database operations for PasswordReset model.
type PasswordResetRepository struct {
	*GenericRepository[models.PasswordReset]
}

// NewPasswordResetRepository creates a new PasswordResetRepository.
func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{
		GenericRepository: NewGenericRepository(db, &models.PasswordReset{}),
	}
}

// FindByEmail finds password reset records by email.
func (r *PasswordResetRepository) FindByEmail(email string) ([]*models.PasswordReset, error) {
	var resets []*models.PasswordReset
	query := r.DB.Where("email = ? AND used = ?", email, false).Find(&resets)
	return resets, query.Error
}

// FindByToken finds a password reset record by token.
func (r *PasswordResetRepository) FindByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	query := r.DB.Where("token = ? AND used = ?", token, false).First(&reset)
	if query.Error != nil {
		return nil, query.Error
	}
	return &reset, nil
}

// InvalidateByEmail marks all password reset tokens for an email as used.
func (r *PasswordResetRepository) InvalidateByEmail(email string) error {
	query := r.DB.Model(&models.PasswordReset{}).Where("email = ?", email).Update("used", true)
	return query.Error
}

// InvalidateToken marks a specific token as used.
func (r *PasswordResetRepository) InvalidateToken(token string) error {
	query := r.DB.Model(&models.PasswordReset{}).Where("token = ?", token).Update("used", true)
	return query.Error
}
