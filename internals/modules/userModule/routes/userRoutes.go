package userModule

import (
	"backendService/internals/modules/userModule/controller"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController *controller.User_Controller
}

func (ur *UserRouter) SetupRoutes(app *gin.Engine) {

	userRouter := app.Group("api/v1/user")
	{
		userRouter.GET("/:id", ur.userController.GetUser)
		userRouter.POST("/", ur.userController.CreateUser)
		userRouter.GET("/", ur.userController.GetAllUsers)
	}
}

func NewUserRouter(userController *controller.User_Controller) *UserRouter {
	return &UserRouter{
		userController: userController,
	}
}

var UserRoutes *UserRouter = NewUserRouter(controller.UserController)
