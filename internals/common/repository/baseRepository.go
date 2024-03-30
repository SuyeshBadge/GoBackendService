package repository

import (
	"backendService/internals/setup/database"
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"
)

const defaultPageSize = 10

// BaseModel represents the base model for all entities in the repository.
type BaseModel struct {
	ID        uint64         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	IsDeleted bool           `gorm:"boolean" json:"isDeleted" default:"false"`
}

// BaseRepository is a generic repository that provides common database operations.
type BaseRepository[T any] struct {
	Db        *gorm.DB
	tableName string
}

type Database = *gorm.DB

// NewBaseRepository creates a new instance of the BaseRepository with the specified database connection and table name.
// It returns a pointer to the created BaseRepository.
// The type parameter T represents the model type that the repository will operate on.
func NewBaseRepository[T any](db *gorm.DB, tableName string) *BaseRepository[T] {
	// Create a new BaseRepository
	repo := &BaseRepository[T]{
		Db:        db,
		tableName: tableName,
	}

	// Set table name
	repo.Db = database.Db.Table(tableName)

	return repo
}
func (r *BaseRepository[T]) Create(model *T) error {

	// Set the BaseModel fields in the input model
	modelValue := reflect.ValueOf(model).Elem()
	baseModelType := reflect.TypeOf(BaseModel{})
	// Create a slice of reflect.StructField with the same length as the number of fields in baseModelType.
	baseModelFields := make([]reflect.StructField, baseModelType.NumField())

	for i := 0; i < baseModelType.NumField(); i++ {
		baseModelFields[i] = baseModelType.Field(i)
	}

	baseModelStructType := reflect.StructOf(baseModelFields)

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelValue.Type().Field(i)

		if fieldType.Anonymous && fieldType.Type == baseModelType {
			field.Set(reflect.Zero(baseModelStructType))
			break
		}
	}

	return r.Db.Create(model).Error
}

func (r *BaseRepository[T]) Update(filter any, update any) error {
	//find the record with filter and update the record with update
	model := new(T)
	return r.Db.Model(model).Where(filter).Updates(update).Error

}

// FindByID retrieves a record from the database based on the given ID.
// It assigns the result to the provided model pointer.
// Returns an error if the record is not found or if there is an issue with the database operation.
func (r *BaseRepository[T]) Delete(id uint64) error {
	model := new(T)

	err := r.Db.First(model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	//update the record to soft delete
	now := time.Now()
	modelValue := reflect.ValueOf(model).Elem()

	deletedAtField := modelValue.FieldByName("DeletedAt")
	if deletedAtField.IsValid() && deletedAtField.CanSet() {
		deletedAtField.Set(reflect.ValueOf(gorm.DeletedAt{Time: now}))
	}

	isDeletedField := modelValue.FieldByName("IsDeleted")
	if isDeletedField.IsValid() && isDeletedField.CanSet() {
		isDeletedField.Set(reflect.ValueOf(true))
	}

	return r.Db.Save(model).Error
}

// FindByID retrieves a record from the database based on the given ID.
// It uses the Unscoped method to include soft-deleted records.
// The retrieved record is stored in the 'model' variable.
// If an error occurs during the retrieval, it is returned.
func (r *BaseRepository[T]) FindByID(id uint64) (*T, error) {
	var model T

	err := r.Db.Scopes(AllowNonDeletedRecords).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

// FindAllWithPagination retrieves a paginated list of models from the database.
// It applies pagination using the provided page number and page size.
// The retrieved models are stored in the 'models' slice.
// Returns an error if there was a problem executing the query.
func (r *BaseRepository[T]) FindAll(page, pageSize int) ([]T, error) {

	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	var models []T

	err := r.Db.Scopes(paginateScope(page, pageSize), AllowNonDeletedRecords).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *BaseRepository[T]) FindAllBy(conditions map[string]interface{}, page, pageSize int) ([]T, error) {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	var models []T
	err := r.Db.Scopes(paginateScope(page, pageSize), func(db *gorm.DB) *gorm.DB {
		for key, value := range conditions {
			db = db.Where(key, value)
		}
		return db.Unscoped()
	}).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *BaseRepository[T]) FindOneBy(conditions map[string]interface{}) (*T, error) {
	var model T
	err := r.Db.Scopes(func(db *gorm.DB) *gorm.DB {
		for key, value := range conditions {
			db = db.Where(key, value)
		}
		return db.Unscoped()
	}).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

// paginateScope is a helper function that returns a closure function to be used as a scope in GORM queries for pagination.
// It takes in two parameters: page and pageSize, which specify the current page number and the number of records per page, respectively.
// The returned closure function takes a GORM DB instance as input and returns a modified DB instance with pagination applied.
// The pagination is based on the provided page and pageSize values.
func paginateScope(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
}

// Filter deleted records
func AllowNonDeletedRecords(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = ?", false)
}
