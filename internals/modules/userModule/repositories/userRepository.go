package repository

import (
	"backendService/internals/common/repository"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// User represents a user entity.
type User struct {
	repository.BaseModel
	UserId          ulid.ULID  `json:"userId" gorm:"uniueIndex"`
	Email           string     `json:"email" gorm:"uniqueIndex"`
	DOB             *time.Time `json:"dob,omitempty" gorm:"type:timestamp"`
	Password        string     `json:"password"`
	FirstName       string     `json:"firstName"`
	LastName        string     `json:"lastName" `
	IsEmailVerified bool       `json:"isEmailVerified" gorm:"type:boolean"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt,omitempty" gorm:"type:timestamp"`
	IsActive        bool       `json:"isActive" gorm:"type:boolean"`
	Mobile          *string    `json:"mobile"`
	IsMoileVerified bool       `json:"isMobileVerified" gorm:"type:boolean"`
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
