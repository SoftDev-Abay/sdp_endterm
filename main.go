package main

import (
	"bufio"
	"fmt"
	"os"
	db "shop/db_f"
	"shop/new_users"
	"shop/products"
	"strconv"
	"strings"
)

// "bufio"
// "fmt"
// "os"

func main() {
	db.GetDBInstance()

	var currentUserID int
	var currentUserIsAdmin bool

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\nWelcome to the Go Shop!")
	fmt.Println("1. Login in")
	fmt.Println("2. Registration")
	fmt.Println("3. Exit")
	fmt.Print("Enter option: ")

	scanner.Scan()
	action := scanner.Text()

	switch action {
	case "1":
		fmt.Println("Please enter your username:")
		scanner.Scan()
		username := scanner.Text()

		fmt.Println("Please enter your password:")
		scanner.Scan()
		password := scanner.Text()

		userID, isAdmin, err := db.LoginUser(username, password)
		if err != nil {
			fmt.Println("Error logging in:", err)
			return
		}
		currentUserID = userID
		currentUserIsAdmin = isAdmin
		fmt.Println("Logged in successfully!")

	case "2":
		fmt.Println("Choose a username:")
		scanner.Scan()
		username := scanner.Text()

		fmt.Println("Choose a password:")
		scanner.Scan()
		password := scanner.Text()

		fmt.Println("Enter your email:")
		scanner.Scan()
		email := scanner.Text()

		fmt.Println("Enter your phone number:")
		scanner.Scan()
		phoneNum := scanner.Text()

		if err := db.RegisterUser(username, password, email, phoneNum, false); err != nil {
			fmt.Println("Error registering:", err)
			return
		}

		fmt.Println("Registered successfully! You can now login.")

	case "3":
		fmt.Println("Bye!")
		return

	default:
		fmt.Println("Invalid option, exiting.")
		return
	}

	for {
		fmt.Println("1. Buy Product")
		fmt.Println("2. Cart")
		fmt.Println("3. Exit")
		fmt.Println("4. Admin Panel (if you are admin)")
		fmt.Print("Enter option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			for {
				fmt.Println("Choose products by ID to add to cart:")
				viewProducts() // Make sure this function prints out products with their IDs
				fmt.Println("Enter 0 to exit")
				fmt.Print("Enter product ID: ")
				productIDStr, _ := reader.ReadString('\n')
				productIDStr = strings.TrimSpace(productIDStr)
				productID, err := strconv.Atoi(productIDStr)
				if productID == 0 {
					break
				}
				if err != nil {
					fmt.Println("Invalid product ID")
					continue
				}

				fmt.Print("Enter quantity: ")
				quantityStr, _ := reader.ReadString('\n')
				quantityStr = strings.TrimSpace(quantityStr)
				quantity, err := strconv.Atoi(quantityStr)
				if err != nil {
					fmt.Println("Invalid quantity")
					continue
				}

				// Call the AddToCart function
				err = db.AddToCart(currentUserID, productID, quantity)
				if err != nil {
					fmt.Println("Error adding to cart:", err)
					continue
				}

				fmt.Println("Product added to cart successfully!")
				break
			}
		case "2":
			fmt.Println("Viewing your cart items:")
			cartItems, err := new_users.ViewCart(currentUserID)
			if err != nil {
				fmt.Println("Error retrieving cart items:", err)
				continue
			}

			if len(cartItems) == 0 {
				fmt.Println("Your cart is empty.")
			} else {
				fmt.Println("Your Cart:")
				for _, item := range cartItems {
					fmt.Printf("Product ID: %d, Product Name: %s, Quantity: %d\n", item.ProductID, item.ProductName, item.Quantity)
				}
			}
		case "3":
			fmt.Println("Thank you for visiting Go Shop!")
			return
		case "4":
			if currentUserIsAdmin {
				for {
					fmt.Println("Hello, Admin!")
					fmt.Println("1. Delete Product by iD")
					fmt.Println("2. Add Product")
					fmt.Println("3. Exit Admin console")
				}
			} else {
				fmt.Println("You are not an admin!")
				continue
			}
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}

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

	//product := products.Product{
	//	Name:  "iphoe 12 pro",
	//	Desc:  "modern phone, cool resolution, 128GB",
	//	Price: 3000,
	//	Category: map[int]string{
	//		1: "clothes",
	//		4: "mobile phones",
	//		5: "food",
	//	},
	//}
	//product, err := db.InsertProduct(product)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("Insert product successfully")
	//	fmt.Println(product)
	//}

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

func viewProducts() {
	// Use the GetProducts function that you have previously provided
	productsList, err := db.GetProducts() // Assuming this function is in the 'products' package
	if err != nil {
		fmt.Println("Error retrieving products:", err)
		return
	}

	fmt.Println("\nList of Products:")
	for _, product := range productsList {
		id, name, price, desc, categories := product.GetDetails()
		fmt.Printf("ID: %d, Name: %s, Price: %d, Description: %s, Categories: %v\n", id, name, price, desc, categories)
	}
}

func addProduct(reader *bufio.Reader) {
	// Add product details
	fmt.Print("Enter product name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter product description: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)

	fmt.Print("Enter product price: ")
	priceStr, _ := reader.ReadString('\n')
	priceStr = strings.TrimSpace(priceStr)
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		fmt.Println("Invalid price. Please enter a number.")
		return
	}

	newProduct := products.Product{Name: name, Desc: desc, Price: price}
	_, err = db.InsertProduct(newProduct)
	if err != nil {
		fmt.Println("Error adding product:", err)
		return
	}

	fmt.Println("Product added successfully!")
}
