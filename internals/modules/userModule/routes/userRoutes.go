package userModule

import (
	userModule "backendService/internals/modules/userModule/controllers"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController *userModule.User_Controller
}

func (ur *UserRouter) SetupRoutes(app *gin.Engine) {

	userRouter := app.Group("api/v1/user")
	{
		userRouter.GET("/:id", ur.userController.GetUser)
		userRouter.POST("/", ur.userController.CreateUser)
		userRouter.GET("/", ur.userController.GetAllUsers)
	}
}

func NewUserRouter(userController *userModule.User_Controller) *UserRouter {
	return &UserRouter{
		userController: userController,
	}
}

var UserRoutes *UserRouter = NewUserRouter(userModule.UserController)
