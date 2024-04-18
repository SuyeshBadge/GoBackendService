package userService

import (
	appError "backendService/internals/common/errors"
	"backendService/internals/modules/userModule/userModule"
	repository "backendService/internals/modules/userModule/userRepository"

	"strconv"

	"github.com/oklog/ulid/v2"
)

type Filter = map[string]interface{}

// UserService is a struct that represents the service for the user model
type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService creates a new instance of UserService.
// It takes a pointer to a UserRepository and returns a pointer to UserService.
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

// CreateUser creates a new user with the provided user data.
// It takes a CreateUserBody object as input and returns the created user and an error interface.
func (us *UserService) CreateUser(createUserData userModule.CreateUserBody) (*repository.User, *appError.ApplicationError) {
	// Generate a new UUID for the user ID
	userId := ulid.Make()

	exixtingUser, err := us.userRepository.FindOneBy(Filter{
		"email": createUserData.Email,
	})
	if err != nil {
		return nil, appError.NewApplicationError("internal_error", "failed to find user")
	}
	if exixtingUser != nil {
		return nil, appError.NewApplicationError("user_exists", "user with this email already exists")
	}

	// Map the request data to a User struct
	user := &repository.User{
		UserId:    userId,
		FirstName: createUserData.FirstName,
		LastName:  createUserData.LastName,
		Email:     &createUserData.Email,
		Password:  &createUserData.Password,
		DOB:       createUserData.DOB,
		Mobile:    createUserData.Mobile,
		IsActive:  true,
	}
	createdUser, err := us.userRepository.Create(user)
	if err != nil {
		return nil, appError.NewApplicationError("internal_error", "failed to create user")
	}
	return createdUser, nil
}

// GetUserByID retrieves a user from the repository based on the provided ID.
// It takes an ID as a string parameter and returns a pointer to a User struct and an error.
// If the ID cannot be parsed or the user is not found, an error is returned.
func (us *UserService) GetUserByID(id string) (*repository.User, *appError.ApplicationError) {

	num, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return nil, appError.NewBadRequestError("invalid_id", "invalid user ID")
	}
	user, err := us.userRepository.FindByID(num)
	if err != nil {
		return nil, appError.NewApplicationError("internal_error", "failed to retrieve user")
	}
	if user == nil {
		return nil, appError.NewNotFoundError("user_not_found", "user not found")
	}
	return user, nil
}

// list of users
func (us *UserService) GetUsers() ([]repository.User, *appError.ApplicationError) {
	users, err := us.userRepository.FindAll(1, 10)
	if err != nil {
		return nil, appError.NewBadRequestError("internal_error", "failed to retrieve users")
	}
	return users, nil

}
