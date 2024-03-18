package repository

import (
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// setBaseModelFields sets the BaseModel fields in the input model
func setBaseModelFields(model interface{}, baseModel *BaseModel) error {
	// Use reflection to set the BaseModel fields in the input model
	modelValue := reflect.ValueOf(model).Elem()
	if modelValue.Kind() != reflect.Struct {
		return errors.New("input model must be a struct")
	}

	baseModelType := reflect.TypeOf(BaseModel{})
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelValue.Type().Field(i)

		if fieldType.Anonymous && fieldType.Type == baseModelType {
			baseModelValue := reflect.ValueOf(baseModel).Elem()
			field.Set(baseModelValue)
			break
		}
	}

	return nil
}

// BaseModel represents the base model for all entities in the repository.
type BaseModel struct {
	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	IsDeleted bool      `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}

type Pagination struct {
	Page     int
	PageSize int
}

type Field map[interface{}]interface{}

// BaseRepository is an abstract base repository that provides common functionality for interacting with the database
type BaseRepository[Fields Field] struct {
	db *gorm.DB
}

// NewBaseRepository creates a new BaseRepository instance
func NewBaseRepository[Fields Field](db *gorm.DB) *BaseRepository[Fields] {
	return &BaseRepository[Fields]{db: db}
}

// Create creates a new record in the database
func (r *BaseRepository[T]) Create(model *T) error {
	baseModel := BaseModel{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	// Set the BaseModel fields in the input model
	if err := setBaseModelFields(model, &baseModel); err != nil {
		return err
	}

	return r.db.Create(model).Error
}

// Update updates an existing record in the database by given filter
func (r *BaseRepository[T]) Update(id string, fields Field) error {
	return r.db.Model(&T{}).Where("id = ?", id).Updates(fields).Error
}

// Delete soft deletes an existing record from the database by setting the flag isDeleted to true
func (r *BaseRepository[T]) Delete(id string) error {
	var record T
	result := r.db.First(&record, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}

	result = r.db.Model(&record).Update("isDeleted", true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByID finds a record by its ID, filtering out records where isDeleted = true
func (r *BaseRepository[T]) FindByID(id uint64) (*T, error) {
	var model T
	err := r.db.Where("isDeleted = ?", false).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

// FindAll finds all records in the database, filtering out records where isDeleted = true
func (r *BaseRepository[T]) FindAll(page, pageSize int) ([]*T, error) {
	var models []*T
	offset := (page - 1) * pageSize
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	err := r.db.Where("isDeleted = ?", false).Limit(pageSize).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// FindAllBy finds all records in the database that match the given conditions, filtering out records where isDeleted = true
func (r *BaseRepository[T]) FindAllBy(conditions map[string]interface{}, page, pageSize int) ([]*T, error) {
	var models []*T
	offset := (page - 1) * pageSize
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	err := r.db.Where(conditions).Where("isDeleted = ?", false).Limit(pageSize).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// FindOneBy finds one record in the database that matches the given conditions, filtering out records where isDeleted = true
func (r *BaseRepository[T]) FindOneBy(conditions map[string]interface{}) (T, error) {
	var record T
	result := r.db.Where(conditions).Where("isDeleted = ?", false).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return record, result.Error
	}

	return record, nil
}
