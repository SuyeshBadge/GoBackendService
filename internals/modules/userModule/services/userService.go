package userModule

type UserService struct {
}

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

func (us *UserService) GetUser(id string) User {

	var users = []User{
		{
			"user1",
			21,
			"M",
		},
		{
			"user2",
			23,
			"F",
		},
	}
	var user User
	for _, u := range users {

		if u.Name == id {
			user = u
		}
	}
	return user
}
