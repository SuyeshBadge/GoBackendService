package authService

import (
	"backendService/internals/common/errors"
	authModule "backendService/internals/modules/authModule/dto"
	"backendService/internals/modules/userModule/userService"
)

type AuthService struct {
	userService userService.UserService
}

func (as *AuthService) OtpSignUp(signupData authModule.OtpVerifyBody) (any, *errors.ApplicationError) {
	if signupData.Mobile == nil && signupData.Email == nil {
		return nil, errors.NewBadRequestError("missing_data", "mobile or email is required")
	}
	if signupData.Mobile != nil {
		// Verify OTP for mobile
		// _, err := as.userService.VerifyMobile(signupData.Mobile, signupData.OTP)
		// if err != nil {
		// 	return nil, err
		// }
		// if signupData.Email != nil {
		// 	// Verify OTP for email
		// 	_, err := as.userService.VerifyEmail(signupData.Email, signupData.OTP)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// }
	}
	return nil, nil
}
