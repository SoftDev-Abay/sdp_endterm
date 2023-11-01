package users

type Customer struct {
	Id           int
	UserName     string
	UserPassword string
	Email        string
	PhoneNum     string
}

func (c *Customer) Auth() bool {
	return true
}

func (c *Customer) GetDetails() (int, string, string, string, string) {
	return c.Id, c.UserName, c.UserPassword, c.Email, c.PhoneNum
}

func (c *Customer) SetId(id int) {
	c.Id = id
}

func (c *Customer) SetUserName(username string) {
	c.UserName = username
}

func (c *Customer) SetUserPassword(userpassword string) {
	c.UserPassword = userpassword
}

func (c *Customer) SetEmail(email string) {
	c.Email = email
}

func (c *Customer) SetPhoneNum(phoneNum string) {
	c.PhoneNum = phoneNum
}
