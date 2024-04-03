package services

import (
	"backendService/internals/modules/userModule/dto"
	repository "backendService/internals/modules/userModule/repositories"
	"errors"
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

func (us *User_Service) CreateUser(createUserData dto.CreateUserBody) (*repository.User, error) {
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

func (us *User_Service) GetUserByID(id string) (*repository.User, error) {
	num, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %v", err)
	}
	user, err := us.userRepository.FindByID(num)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
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
