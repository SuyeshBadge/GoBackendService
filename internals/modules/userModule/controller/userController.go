package controller

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"

	controllers "backendService/internals/common/controller"
	"backendService/internals/common/router"
	"backendService/internals/modules/userModule/dto"
	"backendService/internals/modules/userModule/services"
)

type User_Controller struct {
	controllers.BaseController
	userService *services.User_Service
}

// NewUserController creates a new instance of User_Controller with a provided userService
// dependency.
func NewUserController(userService *services.User_Service) *User_Controller {
	return &User_Controller{userService: userService}
}

// GetUser retrieves a user from the database.
func (uc *User_Controller) GetUser(c *gin.Context) (router.Response, any) {
	id := c.Param("id")
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		log.Println("Error in GetUserByID: ", err)
		return router.Response{}, err
	}
	message := "User retrieved successfully"
	if user == nil {
		message = "User not found"
	}
	return router.Response{Data: user, Message: message}, nil
}

// CreateUser handles the creation of a user. It reads the request body, parses it into a CreateUserData struct,
// and passes the data to the UserService's CreateUser method. CreateUser validates the request body and creates a new user.
func (uc *User_Controller) CreateUser(c *gin.Context) (router.Response, any) {
	var createData dto.CreateUserBody
	// if err := c.ShouldBindJSON(&createData); err != nil {
	// 	return router.Response{}, errors.New("failed to parse request body")
	// }

	_, err := uc.TransformAndValidate(c, createData)
	if err != nil {
		return router.Response{}, err
	}

	user, err1 := uc.userService.CreateUser(createData)
	if err1 != nil {
		return router.Response{}, errors.New("failed to create user")
	}

	return router.Response{Data: user, Message: "User created successfully"}, nil
}

// GetAllUsers retrieves all users from the database.
func (uc *User_Controller) GetAllUsers(c *gin.Context) (router.Response, any) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		return router.Response{}, errors.New("failed to retrieve users")
	}
	return router.Response{Data: users, Message: "Users retrieved successfully"}, nil
}
