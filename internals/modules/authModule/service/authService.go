package authService

import (
	"backendService/internals/common/errors"
	authModule "backendService/internals/modules/authModule/dto"
	"backendService/internals/modules/userModule/userService"
	"strconv"
)

type AuthService struct {
	userService *userService.UserService
	otpService  *OtpService
}

// NewAuthService creates a new instance of AuthService with the provided UserService and OtpService.
// The returned AuthService will use the given services to handle user and OTP operations.
func NewAuthService(userService userService.UserService, otpService OtpService) *AuthService {
	return &AuthService{userService: &userService, otpService: &otpService}
}

// SendOtp sends an OTP (One-Time Password) to the provided mobile or email address.
// It returns true if the OTP was successfully sent, or an ApplicationError if there was an error sending the OTP.
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

// VerifyOtp verifies the provided OTP for the given mobile or email address.
// It returns true if the OTP is valid, or an ApplicationError if the OTP is invalid or other errors occur.
func (as *AuthService) VerifyOtp(verifyOtpData authModule.OtpVerifyBody) (any, *errors.ApplicationError) {

	if verifyOtpData.Mobile == nil && verifyOtpData.Email == nil {
		return nil, errors.NewBadRequestError("missing_data", "mobile or email is required")
	}

	var recipient string
	if verifyOtpData.Mobile != nil {
		recipient = *verifyOtpData.Mobile
	} else {
		recipient = *verifyOtpData.Email
	}

	otp, _ := strconv.Atoi(verifyOtpData.OTP)
	otpVerifyRequest := VerifyOtpRequest{
		Key: recipient,
		Otp: otp,
	}

	_, err := as.otpService.VerifyOtp(otpVerifyRequest)
	if err != nil {
		return nil, err
	}

	return true, nil

}
