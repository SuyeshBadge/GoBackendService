package authService

import "backendService/internals/common/errors"

type OtpService struct {
}

func (os *OtpService) SendOtp(sendOtpdata interface{}) (bool, *errors.ApplicationError) {
	// generate otp

	// save otp in cache for specific time frame

	// send otp to the user
	return false, nil
}

func (os *OtpService) VerifyOtp(verifyOtpData interface{}) (bool, *errors.ApplicationError) {
	// fetch the otp from the cache service

	// if found verify and return the success flag

	// if not found or incorrect otp return failure flag.
	return false, nil
}

func generateOtp(length int) int {
	// genrate a random number with the provided length and return the number
	return 0
}

func (os *OtpService) saveOtpInCache(key string, value int) {
	// save the otp in cache with the proper key prefix using the cacheService
}
