package controller

import (
	"backendService/internals/common/router"
	"backendService/internals/modules/userModule/services"
	"errors"

	"encoding/json"
	"log"

	"io"

	"github.com/gin-gonic/gin"
)

type User_Controller struct {
	userService *services.User_Service
}

// The NewUserController function creates a new instance of User_Controller with a provided userService
// dependency.
func NewUserController(userService *services.User_Service) *User_Controller {
	return &User_Controller{
		userService: userService,
	}
}

// GetUser retrieves a user from the database.
// It takes a gin.Context object as a parameter and returns a router.Response containing a repository.User object and an error.
func (uc *User_Controller) GetUser(c *gin.Context) (router.Response, error) {
	id := c.Param("id")
	log.Println(id)
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		return router.Response{}, errors.New("failed to retrieve user")
	}
	return router.Response{
		Data:    user,
		Message: "User retrieved successfully",
	}, nil
}

// CreateUser handles the creation of a user.
// It reads the request body, parses it into a CreateUserData struct,
// and passes the data to the UserService's CreateUser method.
func (uc *User_Controller) CreateUser(c *gin.Context) error {
	body, err := io.ReadAll(io.Reader(c.Request.Body))
	if err != nil {
		return errors.New("failed to read request body")
	}
	defer c.Request.Body.Close()

	// Parse request body into CreateUserData struct
	var createData services.CreateUserData
	if err := json.Unmarshal(body, &createData); err != nil {
		return errors.New("failed to parse request body")
	}
	log.Println(createData)

	// Pass createData to UserService's CreateUser method
	if err := uc.userService.CreateUser(createData); err != nil {
		return errors.New("failed to create user")
	}
	return nil
}

// GetAllUsers retrieves all users from the database.
// It returns a slice of User objects and an error if any.
func (uc *User_Controller) GetAllUsers(c *gin.Context) (router.Response, error) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		return router.Response{}, errors.New("failed to retrieve users")
	}
	return router.Response{
		Data:    users,
		Message: "Users retrieved successfully",
	}, nil
}
