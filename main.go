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
	var currentUserBalance int

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

		user, err := new_users.LoginUser(username, password)
		if err != nil {
			fmt.Println("Error logging in:", err)
			return
		}

		currentUserID = user.UserID
		currentUserBalance = user.Balance
		currentUserIsAdmin = false // assume the user is not an admin by default.
		if _, ok := user.Permissions.(*new_users.AdminPermissions); ok {
			currentUserIsAdmin = true // set to true if the permissions type is AdminPermissions
		}

		fmt.Println("Logged in successfully!")

	case "2":
		var factory new_users.IUserFactory

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

		fmt.Println("Enter your balance:")
		scanner.Scan()
		balanceStr := scanner.Text()

		balance, err := strconv.Atoi(balanceStr)
		if err != nil {
			fmt.Println("Error parsing balance:", err)
			return
		}

		factory = &new_users.RegularUserFactory{}

		// use the factory to create and register the user
		err = new_users.Register(factory, username, password, email, phoneNum, balance)
		if err != nil {
			fmt.Println("Error registering:", err)
			return
		}

		fmt.Println("Registered successfully! You can now login.")
		return
	case "3":
		fmt.Println("Bye!")
		return

	default:
		fmt.Println("Invalid option, exiting.")
		return
	}

	for {
		fmt.Printf("Your current balance: %v\n", currentUserBalance)
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
				fmt.Printf("Your current balance: %v\n", currentUserBalance)
				fmt.Println("Choose products by ID to add to cart:")
				viewProducts()
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

				err = db.AddToCart(currentUserID, productID, quantity)
				if err != nil {
					fmt.Println("Error adding to cart:", err)
					continue
				}

				fmt.Println("Product added to cart successfully!")
				break
			}
		case "2":
			fmt.Printf("Your current balance: %v\n", currentUserBalance)
			fmt.Println("Viewing your cart items:")
			cartItems, err := new_users.ViewCart(currentUserID)
			if err != nil {
				fmt.Println("Error retrieving cart items:", err)
				continue
			}

			if len(cartItems) == 0 {
				fmt.Println("Your cart is empty.")
			} else {
				for {
					fmt.Printf("Your current balance: %v\n", currentUserBalance)
					fmt.Println("Your Cart:")
					for _, item := range cartItems {
						fmt.Printf("Product ID: %d, Product Name: %s, Quantity: %d, Total Price: %d\n", item.ProductID, item.ProductName, item.Quantity, item.TotalPrice)
					}
					fmt.Println("1. Buy")
					fmt.Println("2. Quit from cart")
					fmt.Print("Your option: ")
					cartOption, _ := reader.ReadString('\n')
					cartOption = strings.TrimSpace(cartOption)

					if cartOption == "1" {
						currentBalance, err := db.BuyProducts(currentUserID)
						currentUserBalance = currentBalance
						if err != nil {
							fmt.Println("Error during purchase:", err)
						} else {
							fmt.Println("Purchase successful!")
						}
					} else if cartOption == "2" {
						fmt.Println("Back to the menu")
						break
					}
					fmt.Println("Invalid option")
					break
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
					fmt.Println("3. Add Category")
					fmt.Println("4. Exit Admin console")
					adminOption, _ := reader.ReadString('\n')
					adminOption = strings.TrimSpace(adminOption)

					if adminOption == "4" {
						fmt.Println("bye!")
						break
					}
					switch adminOption {
					case "1":
						fmt.Println("Enter ID please")
						productIDStr, _ := reader.ReadString('\n')
						productIDStr = strings.TrimSpace(productIDStr)
						productID, err := strconv.Atoi(productIDStr)
						if err != nil {
							fmt.Println("Invalid product ID")
							continue
						}

						db.DeleteProduct(productID)
						fmt.Printf("Product: %v has deleted!", productID)

					case "2":
						// Абай ебашь
						// haha okay
						fmt.Println("Adding product:")
						addProduct(reader)
					case "3":
						fmt.Println("Enter new category name please")
						categoryNameInput, _ := reader.ReadString('\n')
						categoryNameInput = strings.TrimSpace(categoryNameInput)
						err := db.InsertCategory(categoryNameInput)
						if err != nil {
							fmt.Println("Error adding category", err)
							continue
						}
						fmt.Println("successfully added category !")
					}
				}
			} else {
				fmt.Println("You are not an admin!")
				continue
			}
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func viewProducts() {
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
	allCategoriesMap, err := db.GetCategoriesMap()
	if err != nil {
		fmt.Println("Error getting categories:", err)
		return
	}
	fmt.Println("Choose a category to add: ")
	productCategoryMap := make(map[int]string)
	for {
		for categoryId, categoryName := range allCategoriesMap {
			fmt.Println(categoryId, categoryName)
		}
		fmt.Println("Write a category num (write `-1` if you are done)")
		categoryIdInputStr, _ := reader.ReadString('\n')

		categoryIdInputStr = strings.TrimSpace(categoryIdInputStr)
		categoryIdInputInt, err := strconv.Atoi(categoryIdInputStr)

		if categoryIdInputInt < 0 {
			break
		}
		if err != nil {
			fmt.Println("Error getting category id input:", err)
		}
		// check if id is inside category map, id is correct
		if !mapContains(allCategoriesMap, categoryIdInputInt) {
			fmt.Println("Id is invalid:")
			continue
		}
		productCategoryMap[categoryIdInputInt] = allCategoriesMap[categoryIdInputInt]
	}

	newProduct := products.Product{Name: name, Desc: desc, Price: price, Category: productCategoryMap}
	_, err = db.InsertProduct(newProduct)
	if err != nil {
		fmt.Println("Error adding product:", err)
		return
	}

	fmt.Println("Product added successfully!")
}

func mapContains(mapInput map[int]string, elem int) bool {
	for id, _ := range mapInput {
		if id == elem {
			return true
		}
	}
	return false
}
