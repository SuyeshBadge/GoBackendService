package dto

type CreateUserBody struct {
	Name     string `json:"name" validate:"required"`
	Age      int    `json:"age" validate:"required,min=1,max=100"`
	UserName string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required"`
	Mobile   string `json:"mobile" validate:"len=10"`
}
