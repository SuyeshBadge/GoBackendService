package userModule

import (
	userModule "backendService/internals/modules/userModule/routes"
	"backendService/internals/modules/userModule/userController"
	repository "backendService/internals/modules/userModule/userRepository"
	"backendService/internals/modules/userModule/userService"
	"backendService/internals/setup/server"
)

var (
	UserRouter  *userModule.UserRouter
	UserService *userService.UserService
)

func Initialize() {

	userRepository := repository.NewUserRepository(server.Server.Db)
	userService := userService.NewUserService(userRepository)
	userController := userController.NewUserController(userService)
	userRouter := userModule.NewUserRouter(userController)

	// Export
	UserService = userService
	UserRouter = userRouter
}
