package dto

import "time"

// CreateUserBody represents the request body for creating a new user.
// It includes the user's name, age, username, password, and mobile number.
// All fields are required and have specific validation rules.
type CreateUserBody struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	//optional fields
	DOB    *time.Time `json:"dob,omitempty"`
	Mobile *string    `json:"mobile,omitempty" validate:"omitempty,min=10,max=10"`
}
