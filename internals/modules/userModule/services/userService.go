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
	println(users)
	for _, u := range users {
		println(u.name)
		if u.name == id {
			user = u
		}
	}
	println(1234)
	return user
}
