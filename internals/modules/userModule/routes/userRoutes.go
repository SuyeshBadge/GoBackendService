package userModule

import (
	"backendService/internals/common/router"
	"backendService/internals/modules/userModule/userController"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController *userController.UserController
}

func (ur *UserRouter) SetupRoutes(app *gin.Engine) {

	router := router.NewBaseRouter("UserRouter", app)

	userRouter := router.Group("api/v1/user")
	{

		userRouter.GET("/:id", ur.userController.GetUser)
		userRouter.GET("/", ur.userController.GetAllUsers)
		userRouter.POST("/", ur.userController.CreateUser)

	}
}

func NewUserRouter(userController *userController.UserController) *UserRouter {
	return &UserRouter{
		userController: userController,
	}
}
