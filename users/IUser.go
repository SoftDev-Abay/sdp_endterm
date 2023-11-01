package users

type IUser interface {
	Auth() bool
	GetDetails() (int, string, string, string, string)
	SetId(int)
	SetUserName(string)
	SetUserPassword(string)
	SetEmail(string)
	SetPhoneNum(string)
}
