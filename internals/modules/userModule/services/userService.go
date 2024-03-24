package userModule

import userModule "backendService/internals/modules/userModule/repositories"

// UserService is a struct that represents the service for the user model
type UserService struct {
	userRepository userModule.UserRepository
}

// NewUserService creates a new instance of UserService.
// It takes a pointer to a UserRepository and returns a pointer to UserService.
func NewUserService(userRepository userModule.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
