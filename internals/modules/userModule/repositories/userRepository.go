package userModule

type UserRepository struct {
}

type User struct {
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	UserName string  `json:"username"`
	Password string  `json:"password"`
	Mobile   *string `json:"mobile,omitempty"`
}
