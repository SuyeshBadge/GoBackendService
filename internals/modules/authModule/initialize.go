package authModule

import (
	"backendService/internals/common/cache"
	authController "backendService/internals/modules/authModule/controller"
	authRoutes "backendService/internals/modules/authModule/routes"
	authService "backendService/internals/modules/authModule/service"
	"backendService/internals/modules/userModule"
)

var (
	AuthRouter  *authRoutes.AuthRoutes
	AuthService *authService.AuthService
)

func Initialize() {
	otpService := authService.NewOtpService(cache.Cache)
	authService := authService.NewAuthService(*userModule.UserService, *otpService)
	authController := authController.NewAuthController(*authService)
	authRouter := authRoutes.NewAuthRoutes(authController)

	// Export
	AuthRouter = authRouter
	AuthService = authService

}
