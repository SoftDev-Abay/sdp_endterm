package users

type Admin struct {
	Id           int
	UserName     string
	UserPassword string
	Email        string
	PhoneNum     string
}

func (a Admin) Auth() bool {
	return true
}
func (a Admin) GetDetails() (string, string, string, string) {
	return a.UserName, a.UserPassword, a.Email, a.PhoneNum
}
