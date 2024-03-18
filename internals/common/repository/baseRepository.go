package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// BaseModel represents the common fields for all models
type BaseModel struct {
	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}

// BaseRepository is an abstract base repository that provides common functionality for interacting with the database
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository creates a new BaseRepository instance
func NewBaseRepository(db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create creates a new record in the database
func (r *BaseRepository[T]) Create(model *T) error {
	return r.db.Create(model).Error
}

// Update updates an existing record in the database
func (r *BaseRepository[T]) Update(model *T) error {
	return r.db.Save(model).Error
}

// Delete deletes an existing record from the database
func (r *BaseRepository[T]) Delete(model *T) error {
	return r.db.Delete(model).Error
}

// FindByID finds a record by its ID
func (r *BaseRepository[T]) FindByID(id uint64) (*T, error) {
	var model T
	err := r.db.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

// FindAll finds all records in the database
func (r *BaseRepository[T]) FindAll(page, pageSize int) ([]*T, error) {
	var models []*T
	offset := (page - 1) * pageSize
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	err := r.db.Limit(pageSize).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
