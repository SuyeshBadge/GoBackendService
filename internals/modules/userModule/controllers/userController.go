package userModule

import "github.com/gin-gonic/gin"

type UserController struct {
	module string
}

func (uc *UserController) GetUser(c *gin.Context) {
	c.JSON(200, "This is new user controller.")
}

func NewUserController() *UserController {
	return &UserController{
		module: "User",
	}
}

var UserControllers UserController = *NewUserController()
