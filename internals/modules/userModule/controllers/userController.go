package userModule

import (
	userModule "backendService/internals/modules/userModule/services"
	"encoding/json"
	"log"

	"io"

	"github.com/gin-gonic/gin"
)

type User_Controller struct {
	userService *userModule.User_Service
}

func (uc *User_Controller) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func NewUserController(userService *userModule.User_Service) *User_Controller {
	return &User_Controller{
		userService: userService,
	}
}

func (uc *User_Controller) CreateUser(c *gin.Context) {
	body, err := io.ReadAll(io.Reader(c.Request.Body))
	if err != nil {
		// Handle error
		return
	}
	defer c.Request.Body.Close()

	// Parse request body into CreateUserData struct
	var createData userModule.CreateUserData
	if err := json.Unmarshal(body, &createData); err != nil {
		// Handle error
		return
	}
	log.Print(createData)

	// Pass createData to UserService's CreateUser method
	if err := uc.userService.CreateUser(createData); err != nil {
		// Handle error
		return
	}
}

func (uc *User_Controller) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetUsers()
	if err != nil {
		// Handle error
		return
	}
	c.JSON(200, users)
}

var UserController *User_Controller = NewUserController(userModule.UserService)
