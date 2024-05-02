package authService

import (
	"backendService/internals/common/errors"
	authModule "backendService/internals/modules/authModule/dto"
	"backendService/internals/modules/userModule/userService"
	"fmt"
)

type AuthService struct {
	userService *userService.UserService
	otpService  *OtpService
}

func NewAuthService(userService userService.UserService, otpService OtpService) *AuthService {
	return &AuthService{userService: &userService, otpService: &otpService}
}

func (as *AuthService) SendOtp(sendOtpData authModule.OtpSendBody) (any, *errors.ApplicationError) {
	if sendOtpData.Mobile == nil && sendOtpData.Email == nil {
		return nil, errors.NewBadRequestError("missing_data", "mobile or email is required")
	}

	var recipient string
	if sendOtpData.Mobile != nil {
		recipient = *sendOtpData.Mobile
	} else {
		recipient = *sendOtpData.Email
	}

	otpSendRequest := OtpSendRequest{
		Recipient: recipient,
	}

	_, err := as.otpService.SendOtp(otpSendRequest)
	if err != nil {
		return nil, err
	}

	return true, nil
}

func (as *AuthService) VerifyOtp(signupData authModule.OtpVerifyBody) (any, *errors.ApplicationError) {

	fmt.Println(signupData)
	if signupData.Mobile == nil && signupData.Email == nil {
		return nil, errors.NewBadRequestError("missing_data", "mobile or email is required")
	}
	if signupData.Mobile != nil {
	}
	return nil, nil
}
