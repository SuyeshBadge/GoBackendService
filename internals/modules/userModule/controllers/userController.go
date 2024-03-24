package userModule

import (
	userModule "backendService/internals/modules/userModule/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService userModule.UserService
	module      string
}

func (uc *UserController) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func NewUserController() *UserController {
	return &UserController{
		module: "User",
	}
}

var UserControllers UserController = *NewUserController()
