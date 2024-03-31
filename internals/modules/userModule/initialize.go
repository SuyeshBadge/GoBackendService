package userModule

import (
	"backendService/internals/modules/userModule/controller"
	repository "backendService/internals/modules/userModule/repositories"
	userModule "backendService/internals/modules/userModule/routes"
	"backendService/internals/modules/userModule/services"
	"backendService/internals/setup/server"
)

var (
	UserRouter  *userModule.User_Router
	UserService *services.User_Service
)

func Initialize() {

	userRepository := repository.NewUserRepository(server.Server.Db)
	userService := services.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	userRouter := userModule.NewUserRouter(userController)

	// Export
	UserService = userService
	UserRouter = userRouter
}
