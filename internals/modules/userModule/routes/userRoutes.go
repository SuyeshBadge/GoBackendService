package userModule

import (
	"backendService/internals/common/router"
	"backendService/internals/modules/userModule/controller"

	"github.com/gin-gonic/gin"
)

type User_Router struct {
	userController *controller.User_Controller
}

func (ur *User_Router) SetupRoutes(app *gin.Engine) {

	router := router.NewBaseRouter("UserRouter", app)

	userRouter := router.Group("api/v1/user")
	{

		userRouter.GET("/:id", ur.userController.GetUser)

	}
}

func NewUserRouter(userController *controller.User_Controller) *User_Router {
	return &User_Router{
		userController: userController,
	}
}
