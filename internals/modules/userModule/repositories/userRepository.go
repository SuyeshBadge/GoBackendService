package repository

import (
	"backendService/internals/common/repository"

	"gorm.io/gorm"
)

// User represents a user entity.
type User struct {
	repository.BaseModel
	Name     string  `json:"name"`             // Name of the user
	Age      int     `json:"age"`              // Age of the user
	Username string  `json:"username"`         // Username of the user
	Password string  `json:"password"`         // Password of the user
	Mobile   *string `json:"mobile,omitempty"` // Mobile number of the user (optional)
}

// UserRepository represents a repository for managing user data.
type User_Repository struct {
	*repository.BaseRepository[User]
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *User_Repository {
	db.Migrator().AutoMigrate(&User{})
	return &User_Repository{
		BaseRepository: repository.NewBaseRepository[User](db, "users"),
	}
}

// func (r *User_Repository) GetTableName() string {
// 	log.Println("GetTableName", r.Db.Name())
// 	return r.Db.Migrator().CurrentDatabase() + "." + r.Db.Statement.Table

// }
