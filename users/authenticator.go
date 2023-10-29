package users

type IAuthenticator interface {
	signIn() bool
}
