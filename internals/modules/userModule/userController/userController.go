package userController

import (
	"github.com/gin-gonic/gin"

	controllers "backendService/internals/common/controller"
	"backendService/internals/common/errors"
	"backendService/internals/common/logger"
	"backendService/internals/common/router"
	"backendService/internals/modules/userModule/userModule"
	"backendService/internals/modules/userModule/userService"
)

type UserController struct {
	controllers.BaseController
	userService *userService.UserService
}

// NewUserController creates a new instance of User_Controller with a provided userService
// dependency.
func NewUserController(userService *userService.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetUser retrieves a user from the database.
func (uc *UserController) GetUser(c *gin.Context) (router.Response, *errors.ApplicationError) {
	id := c.Param("id")
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		logger.Error("controller", "user_controller", "GetUser", err.Message)
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
func (uc *UserController) CreateUser(c *gin.Context) (router.Response, *errors.ApplicationError) {
	var createData userModule.CreateUserBody
	// if err := c.ShouldBindJSON(&createData); err != nil {
	// 	return router.Response{}, errors.New("failed to parse request body")
	// }

	_, err := uc.TransformAndValidate(c, &createData)
	if err != nil {
		return router.Response{}, err
	}

	user, err := uc.userService.CreateUser(createData)
	if err != nil {
		logger.Error("controller", "user_controller", "CreateUser", err.Message)
		return router.Response{}, err
	}

	return router.Response{Data: user, Message: "User created successfully"}, nil
}

// GetAllUsers retrieves all users from the database.
func (uc *UserController) GetAllUsers(c *gin.Context) (router.Response, *errors.ApplicationError) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		return router.Response{}, err
	}
	return router.Response{Data: users, Message: "Users retrieved successfully"}, nil
}
