package authService

import (
	"backendService/internals/common/cache"
	"backendService/internals/common/errors"
	"backendService/internals/common/logger"
	"context"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type OtpService struct {
	cacheService cache.CacheService
}

type OtpSendRequest struct {
	Recipient string
}

type VerifyOtpRequest struct {
	Key string
	Otp int
}

func NewOtpService(cacheService cache.CacheService) *OtpService {
	return &OtpService{cacheService: cacheService}
}

func (os *OtpService) SendOtp(req OtpSendRequest) (bool, *errors.ApplicationError) {
	// Generate OTP
	length := 6 // Specify the desired length of the OTP
	otp := generateOtp(length)

	// Extract necessary data from the request
	recipient := req.Recipient

	// Save OTP in cache for a specific time frame
	os.saveOtpInCache(recipient, otp)

	// Send OTP to the user
	err := os.sendOtpToUser(recipient, otp)
	if err != nil {
		// Handle the error appropriately (e.g., log the error)
		logger.Error("Auth", "OtpService", "SendOtp", "failed to send OTP to user", err)
		return false, errors.NewBadRequestError("failed_to_send_otp", "failed to send OTP to user")
	}

	return true, nil
}

func (os *OtpService) VerifyOtp(verifyOtpData VerifyOtpRequest) (bool, *errors.ApplicationError) {
	// Extract necessary data from the request
	key := verifyOtpData.Key
	userProvidedOtp := verifyOtpData.Otp

	// Prefix the key with a specific identifier for OTPs
	cacheKey := "otp:" + key

	// Fetch the OTP from the cache service
	var storedOtp int
	err := os.cacheService.Get(context.Background(), cacheKey, &storedOtp)
	if err != nil {
		if err == redis.Nil {
			// OTP not found in the cache
			return false, errors.NewBadRequestError("otp_not_found", "OTP not found in cache")
		}
		// Handle other cache retrieval errors
		logger.Error("Auth", "OtpService", "VerifyOtp", "failed to retrieve OTP from cache", err)
		return false, errors.NewInternalServerError("failed_to_retrieve_otp", "failed to retrieve OTP from cache")
	}

	// Verify the OTP
	if userProvidedOtp == storedOtp {
		// OTP is correct
		// Remove the OTP from the cache to prevent reuse
		err := os.cacheService.Delete(context.Background(), cacheKey)
		if err != nil {
			// Handle cache deletion error
			logger.Error("Auth", "OtpService", "VerifyOtp", "failed to delete OTP from cache", err)
		}
		return true, nil
	}

	// OTP is incorrect
	return false, errors.NewBadRequestError("otp_incorrect", "OTP is incorrect")
}

func generateOtp(length int) int {
	// Create a new random number generator with the current time as the seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Calculate the minimum and maximum values for the OTP
	min := int(math.Pow10(length - 1))
	max := int(math.Pow10(length) - 1)

	// Generate a random number between min and max
	otp := rng.Intn(max-min+1) + min

	return otp
}

func (os *OtpService) saveOtpInCache(key string, value int) {
	// Set the expiration time for the OTP (e.g., 5 minutes)
	expiration := time.Minute * 5

	// Prefix the key with a specific identifier for OTPs
	cacheKey := "otp:" + key

	// Save the OTP in the cache using the cacheService
	err := os.cacheService.Set(context.Background(), cacheKey, value, expiration)
	if err != nil {
		// Handle the error appropriately (e.g., log the error)
		logger.Error("Auth", "OtpService", "saveOtpInCache", "failed to save OTP in cache", err)
	}
}

func (os *OtpService) sendOtpToUser(recipient string, otp int) error {
	// Implement the logic to send the OTP to the user (e.g., via SMS or email)
	// For demonstration purposes, we'll just log the OTP
	logger.Info("Auth", "OtpService", "sendOtpToUser", "Sending OTP "+strconv.Itoa(otp)+" to user: "+recipient)
	return nil
}
