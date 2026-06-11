package repositories

import (
	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/models"
)

type QRCodeRepository struct {
	*GenericRepository[models.QRCode]
}

func NewQRCodeRepository(db *gorm.DB) *QRCodeRepository {
	return &QRCodeRepository{
		GenericRepository: NewGenericRepository[models.QRCode](db, &models.QRCode{}),
	}
}
