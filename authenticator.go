package main

type IAuthenticator interface {
	signIn() bool
}
