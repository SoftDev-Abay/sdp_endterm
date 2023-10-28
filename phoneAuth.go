package main

import (
	"fmt"
	"time"
)

type PhoneAuth struct {
	iAuth IAuthenticator
}

func (e *PhoneAuth) signIn() bool {
	authResult := e.iAuth.signIn()
	var phoneNumber string
	fmt.Println("Please enter phoneNumber:")

	fmt.Scanln(&phoneNumber)

	var smsCode int
	fmt.Println("Please enter smsCode:")
	time.Sleep(3 * time.Second)

	fmt.Scanln(&smsCode)

	if len(phoneNumber) > 10 && smsCode == 1234 && authResult {
		return true
	} else {
		return false
	}
}
