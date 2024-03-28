package services

import (
	repository "backendService/internals/modules/userModule/repositories"
	"fmt"
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

func (us *User_Service) CreateUser(createUserData CreateUserData) error {
	user := repository.User{
		Name:     createUserData.name,
		Age:      createUserData.age,
		Username: createUserData.username,
		Password: createUserData.password,
		Mobile:   &createUserData.mobile,
	}

	if err := us.userRepository.CreateUser(&user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (us *User_Service) GetUserByID(id uint64) (*repository.User, error) {
	user, err := us.userRepository.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}
	return user, nil
}

// list of users
func (us *User_Service) GetUsers() ([]repository.User, error) {
	users, err := us.userRepository.FindAll(1, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err)
	}
	return users, nil
}
