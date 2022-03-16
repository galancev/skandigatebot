package users

type PACSUserResponse struct {
	Users [][]string
}

type PACSUser struct {
	Phone  string
	Name   string
	Active bool
}
