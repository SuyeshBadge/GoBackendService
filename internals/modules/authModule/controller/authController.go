package authController

import (
	controllers "backendService/internals/common/controller"
	"backendService/internals/common/errors"
	"backendService/internals/common/router"
	authModule "backendService/internals/modules/authModule/dto"
	authService "backendService/internals/modules/authModule/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	controllers.BaseController
	authService authService.AuthService
}

func NewAuthController(authService authService.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) OtpSignUp(c *gin.Context) (router.Response, *errors.ApplicationError) {
	var signupData authModule.OtpVerifyBody
	_, err := ac.TransformAndValidate(c, &signupData)

	if err != nil {
		return router.Response{}, err
	}

	_, err = ac.authService.VerifyOtp(signupData)

	if err != nil {
		return router.Response{}, err
	}

	return router.Response{Message: "OTP verified successfully"}, nil

}
