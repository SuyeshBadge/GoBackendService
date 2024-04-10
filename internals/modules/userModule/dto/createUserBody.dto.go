package dto

// CreateUserBody represents the request body for creating a new user.
// It includes the user's name, age, username, password, and mobile number.
// All fields are required and have specific validation rules.
type CreateUserBody struct {
	Name     string `json:"name" validate:"required"`              // Name is the user's full name and is required.
	Age      int    `json:"age" validate:"required,min=1,max=100"` // Age is the user's age and must be between 1 and 100.
	UserName string `json:"username" validate:"required,alphanum"` // UserName is the user's chosen username and must be alphanumeric.
	Password string `json:"password" validate:"required"`          // Password is the user's chosen password and is required.
	Mobile   string `json:"mobile" validate:"len=10"`              // Mobile is the user's mobile phone number and must be 10 digits long.
}
