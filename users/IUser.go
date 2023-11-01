package users

type IUser interface {
	Auth() bool
	GetDetails() (string, string, string, string)
}
