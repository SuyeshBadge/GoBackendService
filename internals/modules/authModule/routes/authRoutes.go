package authModule

import (
	"backendService/internals/common/router"
	authController "backendService/internals/modules/authModule/controller"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	AuthController *authController.AuthController
}

func NewAuthRoutes(authController *authController.AuthController) *AuthRoutes {
	return &AuthRoutes{
		AuthController: authController,
	}
}

func (ar *AuthRoutes) SetupRoutes(app *gin.Engine) {
	router := router.NewBaseRouter("AuthRoutes", app)

	authRouter := router.Group("api/v1/auth")
	{
		authRouter.POST("/otp/send", ar.AuthController.SendOtp)
	}

	// authRouter.POST("/signup/otp", ar.AuthController.OtpSignUp)

}
