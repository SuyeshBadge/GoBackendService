package userModule

import "github.com/gin-gonic/gin"

type UserController struct {
	module string
}

func (this *UserController) GetUser(c *gin.Context) {
	c.JSON(200, "This is user this ")
}

func NewUserController() *UserController {
	return &UserController{
		module: "User",
	}
}

var UserControllers UserController = *NewUserController()
