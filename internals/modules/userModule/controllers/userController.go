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
	user := uc.userService.GetUser(c.Param("id"))
	c.JSON(200, user)
}

func NewUserController() *UserController {
	return &UserController{
		module: "User",
	}
}

var UserControllers UserController = *NewUserController()
