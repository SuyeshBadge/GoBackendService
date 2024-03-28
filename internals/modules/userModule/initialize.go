package userModule

import (
	"backendService/internals/modules/userModule/controller"
	repository "backendService/internals/modules/userModule/repositories"
	userModule "backendService/internals/modules/userModule/routes"
	"backendService/internals/modules/userModule/services"
	"backendService/internals/setup/database"
)

var UserRouter *userModule.User_Router

func Initialize() {
	// initialize user module
	userRepository := repository.NewUserRepository(database.Db)
	userService := services.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	userRouter := userModule.NewUserRouter(userController)
	UserRouter = userRouter
}
