package userModule

import "backendService/internals/common/repository"

// UserFields represents the fields of a user.
type UserFields struct {
	Name     string  `json:"name"`             // Name of the user
	Age      int     `json:"age"`              // Age of the user
	UserName string  `json:"username"`         // Username of the user
	Password string  `json:"password"`         // Password of the user
	Mobile   *string `json:"mobile,omitempty"` // Mobile number of the user (optional)
}

// UserRepository represents a repository for managing user data.
type UserRepository struct {
	repository.BaseRepository[UserFields]
}

// NewUserRepository creates a new instance of UserRepository.
// It takes a pointer to a Database and returns a pointer to UserRepository.
// The UserRepository is initialized with a BaseRepository that is created using the provided Database.
func NewUserRepository(db *repository.Database) *UserRepository {
	return &UserRepository{BaseRepository: *repository.NewBaseRepository[UserFields](db)}
}

// FindUserByID retrieves a user from the repository based on the provided ID.
// It returns the user's fields if found, otherwise it returns an error.
func (r *UserRepository) FindUserByID(id uint64) (*UserFields, error) {
	user, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
