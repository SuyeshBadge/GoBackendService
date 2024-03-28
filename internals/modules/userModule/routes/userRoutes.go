package userModule

import (
	"backendService/internals/modules/userModule/controller"

	"github.com/gin-gonic/gin"
)

type User_Router struct {
	userController *controller.User_Controller
}

func (ur *User_Router) SetupRoutes(app *gin.Engine) {

	userRouter := app.Group("api/v1/user")
	{
		userRouter.GET("/:id", ur.userController.GetUser)
		userRouter.POST("/", ur.userController.CreateUser)
		userRouter.GET("/", ur.userController.GetAllUsers)
	}
}

func NewUserRouter(userController *controller.User_Controller) *User_Router {
	return &User_Router{
		userController: userController,
	}
}
