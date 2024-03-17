package userModule

type UserService struct {
}

type User struct {
	name   string
	age    int
	gender rune
}

func (us *UserService) GetUser(id string) User {

	var users = []User{
		{
			"user1",
			21,
			'M',
		},
		{
			"user2",
			23,
			'F',
		},
	}
	var user User
	for _, u := range users {
		if u.name == id {
			user = u
		}
	}
	println(user.name)
	return user
}
