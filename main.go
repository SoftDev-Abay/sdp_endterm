package main

import (
	"fmt"
	db "shop/db_f"
	"shop/products"
)

// "bufio"
// "fmt"
// "os"

func main() {
	// user := &users.Admin{
	// 	UserName:     "CustomerName2",
	// 	UserPassword: "CustomerPassword2",
	// 	Email:        "example2@.com",
	// 	PhoneNum:     "*7777777777",
	// }
	// db.InsertUser(user)
	// user := db.CheckUser("CustomerName2", "CustomerPassword2")
	// if user != nil {
	// 	fmt.Println(user)
	// } else {
	// 	fmt.Println("User not found")
	// }
	// var usersArr []users.IUser
	// usersArr, err_users := db.GetUsers()
	// if err_users != nil {
	// 	fmt.Println(err_users)
	// } else {
	// 	fmt.Println("Get users successfully")
	// 	fmt.Println(usersArr)
	// }
	// for _, v := range usersArr {
	// 	fmt.Println(v)
	// }

	// user, userIdErr := db.GetUserById(4)
	// if userIdErr != nil {
	// 	fmt.Println(userIdErr)
	// } else {
	// 	fmt.Println("Get user by id successfully")
	// 	fmt.Println(user)
	// }
	// user.SetUserName("Damir")
	// user.SetUserPassword("Damir123")
	// newUser, updateErr := db.UpdateUser(user)
	// if updateErr != nil {
	// 	fmt.Println(updateErr)
	// } else {
	// 	fmt.Println("Update user successfully")
	// 	fmt.Println(newUser)
	// }

	// err := db.InsertCategory("pens")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Insert category successfully")
	// }

	// categories, err := db.GetCategories()
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Get categories successfully")
	// 	fmt.Println(categories)
	// }

	// categories, err = db.GetProductCategories(1)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("GetProductCategories successfully")
	// 	fmt.Print(categories)
	// }

	product := products.Product{
		Name:  "iphoe 12 pro",
		Desc:  "modern phone, cool resolution, 128GB",
		Price: 3000,
		Category: map[int]string{
			1: "clothes",
			4: "mobile phones",
			5: "food",
		},
	}
	product, err := db.InsertProduct(product)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Insert product successfully")
		fmt.Println(product)
	}

	// db.AddCategoryToProduct(1, 5)

	// categoriesMap, err := db.GetProductCategoriesMap(1)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Get successfully")
	// 	fmt.Println(categoriesMap)
	// }
	// product := products.Product{
	// 	Id:    3,
	// 	Name:  "laptop",
	// 	Desc:  "modern laptop, 32GB RAM, 1TB SSD",
	// 	Price: 300,
	// }
	// newProduct, update_err := db.UpdateProduct(product)
	// if update_err != nil {
	// 	fmt.Println(update_err)
	// } else {
	// 	fmt.Println("Update product successfully")
	// 	fmt.Println(newProduct)
	// }
	// id := 2
	// delete_err := db.DeleteProduct(id)
	// if delete_err != nil {
	// 	fmt.Println(delete_err)
	// } else {
	// 	fmt.Println("Deleted product successfully", id)
	// }

	// products, err := db.GetProducts()
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Get products successfully")
	// 	fmt.Print(products)
	// }
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
