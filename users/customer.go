package users

type Customer struct {
	id           int
	UserName     string
	UserPassword string
	Email        string
	PhoneNum     string
}

func (c *Customer) Auth() bool {
	return true
}
