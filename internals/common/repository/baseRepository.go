package repository

import (
	"backendService/internals/setup/database"
	"errors"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

const defaultPageSize = 10

// BaseModel represents the base model for all entities in the repository.
type BaseModel struct {
	*gorm.Model
	ID        uint64         `gorm:"primary_key"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BaseRepository is a generic repository that provides common database operations.
type BaseRepository[T any] struct {
	db        gorm.DB
	tableName string
}

type Database = *gorm.DB

// NewBaseRepository creates a new instance of the BaseRepository with the specified database connection and table name.
// It returns a pointer to the created BaseRepository.
// The type parameter T represents the model type that the repository will operate on.
func NewBaseRepository[T any](tableName string) *BaseRepository[T] {
	// Create a new BaseRepository
	repo := &BaseRepository[T]{
		db:        *database.Db,
		tableName: tableName,
	}

	// Set table name
	// repo.db = database.Db.Table(tableName)

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

	return r.db.Create(model).Error
}

// updatedAtField represents the field "UpdatedAt" in the modelValue struct.
// It is used to access and manipulate the value of the "UpdatedAt" field.
func (r *BaseRepository[T]) Update(model *T) error {
	now := time.Now()
	modelValue := reflect.ValueOf(model).Elem()

	updatedAtField := modelValue.FieldByName("UpdatedAt")
	if updatedAtField.IsValid() && updatedAtField.CanSet() {
		updatedAtField.Set(reflect.ValueOf(now))
	}
	return r.db.Save(model).Error
}

// FindByID retrieves a record from the database based on the given ID.
// It assigns the result to the provided model pointer.
// Returns an error if the record is not found or if there is an issue with the database operation.
func (r *BaseRepository[T]) Delete(id uint64) error {
	model := new(T)

	err := r.db.First(model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return r.db.Delete(model).Error
}

// FindByID retrieves a record from the database based on the given ID.
// It uses the Unscoped method to include soft-deleted records.
// The retrieved record is stored in the 'model' variable.
// If an error occurs during the retrieval, it is returned.
func (r *BaseRepository[T]) FindByID(id uint64) (*T, error) {
	var model T

	err := r.db.Unscoped().First(&model, id).Error
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
	fmt.Println("GetUsers", r.db.Name())

	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	var models []T

	err := r.db.Scopes(paginateScope(page, pageSize)).Find(&models).Error
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
	err := r.db.Scopes(paginateScope(page, pageSize), func(db *gorm.DB) *gorm.DB {
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
	err := r.db.Scopes(func(db *gorm.DB) *gorm.DB {
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
