package users

type Customer struct {
	Id           int
	UserName     string
	UserPassword string
	Email        string
	PhoneNum     string
}

func (c Customer) Auth() bool {
	return true
}

func (c Customer) GetDetails() (string, string, string, string) {
	return c.UserName, c.UserPassword, c.Email, c.PhoneNum
}
