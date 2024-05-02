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

func (ac *AuthController) SendOtp(c *gin.Context) (router.Response, *errors.ApplicationError) {
	var sendOtpData authModule.OtpSendBody
	_, err := ac.TransformAndValidate(c, &sendOtpData)

	if err != nil {
		return router.Response{}, err
	}

	_, err = ac.authService.SendOtp(sendOtpData)

	if err != nil {
		return router.Response{}, err
	}
	response := router.Response{Message: "OTP sent successfully"}
	return response, nil
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
