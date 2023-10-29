package users

type Admin struct {
	id           int
	Name         string
	UserName     string
	UserPassword string
	Email        string
	PhoneNum     string
}

func (a *Admin) Auth() bool {
	return true
}
