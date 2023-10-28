package main

import "fmt"

type EmailAuth struct {
	iAuth IAuthenticator
}

func (e *EmailAuth) signIn() bool {
	authResult := e.iAuth.signIn()
	var email string
	fmt.Println("Please enter email:")

	fmt.Scanln(&email)
	validEmail := e.validateEmail(email)

	if validEmail && authResult {
		return true
	} else {
		return false
	}
}

func (e *EmailAuth) validateEmail(email string) bool {
	if len(email) < 5 {
		return false
	}
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			return true
		}
	}

	return false

}
