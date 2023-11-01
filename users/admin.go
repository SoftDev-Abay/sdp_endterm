package users

type Admin struct {
	Id           int
	UserName     string
	UserPassword string
	Email        string
	PhoneNum     string
}

func (a *Admin) Auth() bool {
	return true
}
func (a *Admin) GetDetails() (int, string, string, string, string) {
	return a.Id, a.UserName, a.UserPassword, a.Email, a.PhoneNum
}

func (a *Admin) SetId(id int) {
	a.Id = id
}

func (a *Admin) SetUserName(username string) {
	a.UserName = username
}

func (a *Admin) SetUserPassword(userpassword string) {
	a.UserPassword = userpassword
}

func (a *Admin) SetEmail(email string) {
	a.Email = email
}

func (a *Admin) SetPhoneNum(phoneNum string) {
	a.PhoneNum = phoneNum
}
