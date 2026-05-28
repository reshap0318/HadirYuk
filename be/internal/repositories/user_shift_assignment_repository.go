package repositories

import (
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

type UserShiftAssignmentRepository struct {
	*GenericRepository[models.UserShiftAssignment]
}

func NewUserShiftAssignmentRepository(db *gorm.DB) *UserShiftAssignmentRepository {
	return &UserShiftAssignmentRepository{
		GenericRepository: NewGenericRepository(db, &models.UserShiftAssignment{}),
	}
}

// FindByUserID finds the active shift assignment for a user.
func (r *UserShiftAssignmentRepository) FindByUserID(tx *gorm.DB, userID uint, preloads ...string) (*models.UserShiftAssignment, error) {
	db := r.getDB(tx)
	var instance *models.UserShiftAssignment
	query := db.Model(&instance).Where("user_id = ? AND is_active = ?", userID, true)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.First(&instance).Error; err != nil {
		return nil, err
	}
	return instance, nil
}

// FindByUserIDWithHistory finds all shift assignments for a user (including inactive).
func (r *UserShiftAssignmentRepository) FindByUserIDWithHistory(tx *gorm.DB, userID uint, opts *QueryOptions, preloads ...string) (*PagedResult[models.UserShiftAssignment], error) {
	db := r.getDB(tx)
	var instance *models.UserShiftAssignment
	query := db.Model(&instance).Where("user_id = ?", userID)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = r.applyOptions(query, opts)

	var total int64
	if err := query.Model(&instance).Count(&total).Error; err != nil {
		return nil, err
	}

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

	datas := []models.UserShiftAssignment{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.UserShiftAssignment]{
		Data:       datas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindAllWithSearch finds all assignments with search on user name, email, and shift name.
func (r *UserShiftAssignmentRepository) FindAllWithSearch(tx *gorm.DB, opts *QueryOptions) (*PagedResult[models.UserShiftAssignment], error) {
	db := r.getDB(tx)
	var instance *models.UserShiftAssignment

	query := db.Model(&instance)

	// Preload relations
	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	// Search with JOINs
	if opts.Search != "" {
		searchPattern := "%" + opts.Search + "%"
		query = query.Joins("LEFT JOIN users ON users.id = user_shift_assignments.user_id").
			Joins("LEFT JOIN shifts ON shifts.id = user_shift_assignments.shift_id").
			Where("users.name LIKE ? OR users.email LIKE ? OR shifts.name LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Sorting
	if opts.SortBy != "" {
		order := "ASC"
		if opts.Order == "DESC" {
			order = "DESC"
		}
		query = query.Order("user_shift_assignments." + opts.SortBy + " " + order)
	} else {
		query = query.Order("user_shift_assignments.id DESC")
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	page := 1
	pageSize := 10
	if opts.Page > 0 {
		page = opts.Page
	}
	if opts.PageSize > 0 {
		pageSize = opts.PageSize
	}

	// Pagination
	if pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Limit(pageSize).Offset(offset)
	}

	datas := []models.UserShiftAssignment{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.UserShiftAssignment]{
		Data:       datas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
