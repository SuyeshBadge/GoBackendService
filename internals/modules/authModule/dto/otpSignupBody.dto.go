package authModule

type OtpVerifyBody struct {
	Mobile *string `json:"mobile,omitempty" validate:"omitempty,len=10"` // Mobile should be 10 characters long if present
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`   // Email should be a valid email address if present
	OTP    string  `json:"otp" validate:"required,len=6"`                // OTP is required and should be 6 characters long
}

type OtpSendBody struct {
	Mobile *string `json:"mobile,omitempty" validate:"omitempty,len=10"` // Mobile should be 10 characters long if present
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`   // Email should be a valid email address if present
}
