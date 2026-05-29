package repositories

import (
	"strings"

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

// FindAllExceptSuperAdmin finds all users except those with Super Admin role.
func (r *UserRepository) FindAllExceptSuperAdmin(preloads ...string) ([]models.User, error) {
	var users []models.User
	query := r.DB.Model(&models.User{}).
		Where("id NOT IN (SELECT user_id FROM user_has_roles INNER JOIN roles ON roles.id = user_has_roles.role_id WHERE roles.name = ?)", "Super Admin")

	// Preload relations
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FindAllExceptSuperAdminWithOpts finds all users except Super Admin with query options (pagination, sorting, search, preloads).
func (r *UserRepository) FindAllExceptSuperAdminWithOpts(opts *QueryOptions) (*PagedResult[models.User], error) {
	query := r.DB.Model(&models.User{}).
		Where("id NOT IN (SELECT user_id FROM user_has_roles INNER JOIN roles ON roles.id = user_has_roles.role_id WHERE roles.name = ?)", "Super Admin")

	// Preload relations
	if opts != nil {
		for _, preload := range opts.Preloads {
			query = query.Preload(preload)
		}

		// Sorting
		if opts.SortBy != "" {
			order := "ASC"
			if strings.ToUpper(opts.Order) == "DESC" {
				order = "DESC"
			}
			query = query.Order(opts.SortBy + " " + order)
		}

		// Search
		if opts.Search != "" {
			searchPattern := "%" + opts.Search + "%"
			query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
		}
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	page := 1
	pageSize := 10
	if opts != nil {
		if opts.Page > 0 {
			page = opts.Page
		}
		if opts.PageSize > 0 {
			pageSize = opts.PageSize
		}
	}

	if pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Limit(pageSize).Offset(offset)
	}

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.User]{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}


