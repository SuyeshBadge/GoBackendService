package userModule

import (
	"backendService/internals/common/repository"
	"backendService/internals/setup/database"

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
	repository.BaseRepository[User]
}

// NewUserRepository creates a new instance of UserRepository.
// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *User_Repository {
	return &User_Repository{
		BaseRepository: *repository.NewBaseRepository[User]("users"),
	}
}

// FindUserByID retrieves a user from the repository based on the provided ID.
func (r *User_Repository) FindUserByID(id uint64) (*User, error) {
	user, err := r.FindByID(id)
	return user, err
}

func (r *User_Repository) CreateUser(user *User) error {
	return r.Create(user)
}

var UserRepository = NewUserRepository(database.Db)
