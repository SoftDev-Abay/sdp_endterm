package main

import (
	db "shop/db_f"
	"shop/users"
)

// "bufio"
// "fmt"
// "os"

func main() {
	user := &users.Admin{
		Name:         "Customer2",
		UserName:     "CustomerName2",
		UserPassword: "CustomerPassword2",
		Email:        "example2@.com",
		PhoneNum:     "*7777777777",
	}
	db.InsertUser(user)

	// // Create a new scanner to read from standard input
	// scanner := bufio.NewScanner(os.Stdin)

	// fmt.Println("Enter some text (write 'end' to end):")

	// // Read input line by line
	// for scanner.Scan() {
	// 	text := scanner.Text() // Get the current line of text
	// 	if text == "end" {
	// 		break // Exit loop if an empty line is entered
	// 	}
	// 	fmt.Println("You entered:", text)

	// 	authUsername := &Auth{}

	// 	authWithUserAndEmail := &EmailAuth{
	// 		iAuth: authUsername,
	// 	}

	// 	authWithEmailAndPhone := &PhoneAuth{
	// 		iAuth: authWithUserAndEmail,
	// 	}

	// 	result := authWithEmailAndPhone.signIn()
	// 	if result {
	// 		fmt.Println("Sign in successfully")
	// 	} else {
	// 		fmt.Println("Sign in failed")
	// 	}
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Println("Error:", err)
	// }
}
