package userRepository

import (
	"backendService/internals/common/repository"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// UserRepository represents a repository for managing user data.
type User struct {
	repository.BaseModel

	UserId           ulid.ULID  `json:"userId" gorm:"uniqueIndex"`
	Email            *string    `json:"email" gorm:"uniqueIndex"`
	Username         *string    `json:"username" gorm:"uniqueIndex"`
	DOB              *time.Time `json:"dob,omitempty" gorm:"type:timestamp"`
	Password         *string    `json:"password"`
	FirstName        string     `json:"firstName"`
	LastName         string     `json:"lastName"`
	IsEmailVerified  bool       `json:"isEmailVerified" gorm:"type:boolean"`
	EmailVerifiedAt  *time.Time `json:"emailVerifiedAt,omitempty" gorm:"type:timestamp"`
	IsActive         bool       `json:"isActive" gorm:"type:boolean"`
	Mobile           *string    `json:"mobile" gorm:"uniqueIndex"`
	IsMobileVerified bool       `json:"isMobileVerified" gorm:"type:boolean"`
	AuthProvider     string     `json:"authProvider"`
}
type UserRepository struct {
	*repository.BaseRepository[User]
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	db.Migrator().AutoMigrate(&User{})
	return &UserRepository{
		BaseRepository: repository.NewBaseRepository[User](db, "users"),
	}
}

// func (r *User_Repository) GetTableName() string {
// 	log.Println("GetTableName", r.Db.Name())
// 	return r.Db.Migrator().CurrentDatabase() + "." + r.Db.Statement.Table

// }
