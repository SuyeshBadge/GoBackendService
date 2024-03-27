package userModule

import (
	userModule "backendService/internals/modules/userModule/repositories"
	"fmt"
)

// UserService is a struct that represents the service for the user model
type User_Service struct {
	userRepository *userModule.User_Repository
}

// NewUserService creates a new instance of UserService.
// It takes a pointer to a UserRepository and returns a pointer to UserService.
func NewUserService(userRepository *userModule.User_Repository) *User_Service {
	return &User_Service{userRepository: userRepository}
}

func (us *User_Service) CreateUser(createUserData CreateUserData) error {
	user := userModule.User{
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

var UserService = NewUserService(userModule.UserRepository)
