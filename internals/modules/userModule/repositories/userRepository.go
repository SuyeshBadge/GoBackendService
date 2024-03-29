package repository

import (
	"backendService/internals/common/repository"
	"log"

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

// CreateUser creates a new user in the database.
func (ur *User_Repository) CreateUser(user *User) error {
	if err := ur.Db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindUserByID retrieves a user by ID from the database.
func (ur *User_Repository) FindUserByID(id uint64) (*User, error) {
	var user User
	if err := ur.Db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers retrieves a list of users from the database.
func (ur *User_Repository) GetUsers() ([]User, error) {
	var users []User
	rows, err := ur.Db.Model(&User{}).Select("*").Rows()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	for rows.Next() {
		var user User
		ur.Db.ScanRows(rows, &user)
		users = append(users, user)
	}
	return users, nil
}
