package services

import (
	appError "backendService/internals/common/errors"
	"backendService/internals/modules/userModule/dto"
	repository "backendService/internals/modules/userModule/repositories"

	"fmt"
	"strconv"
)

// UserService is a struct that represents the service for the user model
type User_Service struct {
	userRepository *repository.User_Repository
}

// NewUserService creates a new instance of UserService.
// It takes a pointer to a UserRepository and returns a pointer to UserService.
func NewUserService(userRepository *repository.User_Repository) *User_Service {
	return &User_Service{userRepository: userRepository}
}

// CreateUser creates a new user with the provided user data.
// It takes a CreateUserBody object as input and returns the created user and an error interface.
func (us *User_Service) CreateUser(createUserData dto.CreateUserBody) (*repository.User, interface{}) {
	userData := repository.User{
		Name:     createUserData.Name,
		Age:      createUserData.Age,
		Username: createUserData.UserName,
		Password: createUserData.Password,
		Mobile:   &createUserData.Mobile,
	}

	user, err := us.userRepository.Create(&userData)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	return user, nil
}

// GetUserByID retrieves a user from the repository based on the provided ID.
// It takes an ID as a string parameter and returns a pointer to a User struct and an error.
// If the ID cannot be parsed or the user is not found, an error is returned.
func (us *User_Service) GetUserByID(id string) (*repository.User, *appError.ApplicationError) {
	// log.Println("inside GetUserByID", id)

	num, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return nil, appError.NewBadRequestError("invalid_id", "invalid user ID")
	}
	user, err := us.userRepository.FindByID(num)
	if err != nil {
		return nil, appError.NewApplicationError("internal_error", "failed to retrieve user")
	}
	if user == nil {
		return nil, appError.NewBadRequestError("user_not_found", "user not found")
	}
	return user, nil
}

// list of users
func (us *User_Service) GetUsers() ([]repository.User, error) {
	// log.Println(us.userRepository.GetTableName())
	users, err := us.userRepository.FindAll(1, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err.Error())
	}
	return users, nil

}
