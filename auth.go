package main

import "fmt"

type Auth struct {
}

func (a *Auth) signIn() bool {
	var username string
	fmt.Println("Please enter the username:")

	fmt.Scanln(&username)

	var password string
	fmt.Println("Please enter the password:")

	fmt.Scanln(&password)

	if password == "0000" && username == "abay" {
		return true
	} else {
		return false
	}
}
