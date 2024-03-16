package userModule

import (
	userModule "backendService/internals/modules/userModule/controllers"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController userModule.UserController
}

func (ur *UserRouter) SetupRoutes(app *gin.Engine) {

	userRouter := app.Group("api/v1/user")
	{
		userRouter.GET("/", ur.userController.GetUser)
	}
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		userController: *userModule.NewUserController(),
	}
}

var UserRoutes *UserRouter = NewUserRouter()
