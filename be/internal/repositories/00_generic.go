package repositories

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/helpers"
)

// QueryOptions holds options for querying records.
type QueryOptions struct {
	Page         int      // Page number (default: 1)
	PageSize     int      // Items per page (default: 10, 0 = no pagination)
	SortBy       string   // Field to sort by
	Order        string   // "ASC" or "DESC" (default: "ASC")
	Search       string   // Search keyword
	SearchFields []string // Fields to search
	Preloads     []string // Relations to preload
}

// PagedResult holds paginated query results.
type PagedResult[T any] struct {
	Data       []T   `json:"data"`        // Data records
	Total      int64 `json:"total"`       // Total records
	Page       int   `json:"page"`        // Current page
	PageSize   int   `json:"page_size"`   // Items per page
	TotalPages int   `json:"total_pages"` // Total pages
}

// GetData returns the data slice.
func (p *PagedResult[T]) GetData() interface{} {
	return p.Data
}

// GetMetadata returns the pagination metadata.
func (p *PagedResult[T]) GetMetadata() helpers.PaginationMeta {
	return helpers.PaginationMeta{
		Total:      p.Total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: p.TotalPages,
	}
}

// GenericRepository provides generic CRUD operations for any model
type GenericRepository[T any] struct {
	DB    *gorm.DB
	Model *T
}

// NewGenericRepository creates a new generic repository
func NewGenericRepository[T any](db *gorm.DB, model *T) *GenericRepository[T] {
	return &GenericRepository[T]{
		DB:    db,
		Model: model,
	}
}

// getDB returns the transaction DB if provided, otherwise the default DB
func (r *GenericRepository[T]) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.DB
}

// applyOptions applies query options to a DB query
func (r *GenericRepository[T]) applyOptions(db *gorm.DB, opts *QueryOptions) *gorm.DB {
	if opts == nil {
		return db
	}

	// Preload relations
	for _, preload := range opts.Preloads {
		db = db.Preload(preload)
	}

	// Sorting
	if opts.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(opts.Order) == "DESC" {
			order = "DESC"
		}
		db = db.Order(opts.SortBy + " " + order)
	}

	// Search
	if opts.Search != "" && len(opts.SearchFields) > 0 {
		searchPattern := "%" + opts.Search + "%"
		var conditions []string
		for _, field := range opts.SearchFields {
			conditions = append(conditions, field+" LIKE ?")
		}
		searchCondition := strings.Join(conditions, " OR ")
		args := make([]interface{}, len(opts.SearchFields))
		for i := range args {
			args[i] = searchPattern
		}
		db = db.Where(searchCondition, args...)
	}

	return db
}

// FindByID finds a record by ID
func (r *GenericRepository[T]) FindByID(tx *gorm.DB, id uint, preloads ...string) (*T, error) {
	db := r.getDB(tx)
	var instance *T
	query := db.Model(&instance).Where("id = ?", id)

	// Preload relations
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}
	return instance, nil
}

// FindByIDWithOpts finds a record by ID with query options
func (r *GenericRepository[T]) FindByIDWithOpts(tx *gorm.DB, id uint, opts *QueryOptions) (*T, error) {
	db := r.getDB(tx)
	var instance *T
	query := r.applyOptions(db, opts)

	if err := query.Where("id = ?", id).First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}
	return instance, nil
}

// Create creates a new record
func (r *GenericRepository[T]) Create(tx *gorm.DB, request *T) (*T, error) {
	db := r.getDB(tx)
	var instance *T
	if err := db.Model(&instance).Create(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

// CreateMany creates multiple records
func (r *GenericRepository[T]) CreateMany(tx *gorm.DB, request []T) error {
	db := r.getDB(tx)
	var instance *T
	return db.Model(&instance).Create(&request).Error
}

// Update updates a record by filter
func (r *GenericRepository[T]) Update(tx *gorm.DB, filter *T, update *T) (*T, error) {
	db := r.getDB(tx)
	var instance *T
	if err := db.Model(&instance).Where(filter).First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}

	if err := db.Model(&instance).Where(filter).Updates(update).Error; err != nil {
		return nil, err
	}

	return instance, nil
}

// UpdateMap updates a record by filter using a map for partial updates (supports zero values)
func (r *GenericRepository[T]) UpdateMap(tx *gorm.DB, filter *T, update map[string]interface{}) (*T, error) {
	db := r.getDB(tx)
	var instance *T
	if err := db.Model(&instance).Where(filter).First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}

	if err := db.Model(&instance).Where(filter).Updates(update).Error; err != nil {
		return nil, err
	}

	return instance, nil
}

// Delete deletes a record by ID
func (r *GenericRepository[T]) Delete(tx *gorm.DB, id uint) (*T, error) {
	db := r.getDB(tx)
	var instance *T
	if err := db.Model(&instance).Where("id = ?", id).First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound
		}
		return nil, err
	}

	if err := db.Model(&instance).Where("id = ?", id).Delete(&instance).Error; err != nil {
		return nil, err
	}

	return instance, nil
}

// FindAll finds all records
func (r *GenericRepository[T]) FindAll(tx *gorm.DB, preloads ...string) ([]T, error) {
	db := r.getDB(tx)
	var instance *T
	query := db.Model(&instance)

	// Preload relations
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}

// FindAllWithOpts finds all records with query options (supports pagination, sorting, preloads, search)
func (r *GenericRepository[T]) FindAllWithOpts(tx *gorm.DB, opts *QueryOptions) (*PagedResult[T], error) {
	db := r.getDB(tx)
	var instance *T
	query := r.applyOptions(db, opts)

	// Get total count
	var total int64
	if err := query.Model(&instance).Count(&total).Error; err != nil {
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

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[T]{
		Data:       datas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindByField finds records by filter (note: ignores zero values like 0, false, "")
func (r *GenericRepository[T]) FindByField(tx *gorm.DB, filter *T, preloads ...string) ([]T, error) {
	db := r.getDB(tx)
	var instance *T
	query := db.Model(&instance).Where(filter)

	// Preload relations
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}

// FindByFieldWithOpts finds records by filter with query options (supports pagination, sorting, preloads)
func (r *GenericRepository[T]) FindByFieldWithOpts(tx *gorm.DB, filter *T, opts *QueryOptions) (*PagedResult[T], error) {
	db := r.getDB(tx)
	var instance *T
	query := r.applyOptions(db, opts).Where(filter)

	// Get total count
	var total int64
	if err := query.Model(&instance).Count(&total).Error; err != nil {
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

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[T]{
		Data:       datas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindByFieldMap finds records by filter map (supports zero values like 0, false, "")
func (r *GenericRepository[T]) FindByFieldMap(tx *gorm.DB, filter map[string]interface{}, preloads ...string) ([]T, error) {
	db := r.getDB(tx)
	var instance *T
	query := db.Model(&instance).Where(filter)

	// Preload relations
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}

// FindByFieldMapWithOpts finds records by filter map with query options
func (r *GenericRepository[T]) FindByFieldMapWithOpts(tx *gorm.DB, filter map[string]interface{}, opts *QueryOptions) (*PagedResult[T], error) {
	db := r.getDB(tx)
	var instance *T
	query := r.applyOptions(db, opts).Where(filter)

	// Get total count
	var total int64
	if err := query.Model(&instance).Count(&total).Error; err != nil {
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

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[T]{
		Data:       datas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Count counts all records
func (r *GenericRepository[T]) Count(tx *gorm.DB) (int64, error) {
	db := r.getDB(tx)
	var instance *T
	var count int64
	if err := db.Model(&instance).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists checks if a record exists by filter map (supports zero values like 0, false, "")
func (r *GenericRepository[T]) Exists(tx *gorm.DB, filter map[string]interface{}) (bool, error) {
	db := r.getDB(tx)
	var instance *T
	var count int64

	if err := db.Model(&instance).Where(filter).Limit(1).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByField checks if a record exists by struct filter (note: ignores zero values)
func (r *GenericRepository[T]) ExistsByField(tx *gorm.DB, filter *T) (bool, error) {
	db := r.getDB(tx)
	var instance *T
	var count int64

	if err := db.Model(&instance).Where(filter).Limit(1).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByCondition checks if a record exists using a raw WHERE condition
func (r *GenericRepository[T]) ExistsByCondition(tx *gorm.DB, condition string, args ...interface{}) (bool, error) {
	db := r.getDB(tx)
	var instance *T
	var count int64

	if err := db.Model(&instance).Where(condition, args...).Limit(1).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByCondition finds records using a raw WHERE condition with optional preloads
func (r *GenericRepository[T]) FindByCondition(tx *gorm.DB, condition string, preloads ...string) ([]T, error) {
	db := r.getDB(tx)
	var instance *T
	query := db.Model(&instance).Where(condition)

	// Preload relations
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	datas := []T{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}
