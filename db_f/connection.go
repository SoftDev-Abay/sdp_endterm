package db

import (
	"database/sql"
	"errors"
	"fmt"
	"shop/products"
	"shop/users"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db     *sql.DB
	once   sync.Once
	dbOnce sync.Once
)

func getDB() *sql.DB {
	once.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		dbOnce.Do(func() {
			var err error
			db, err = sql.Open("postgres", psqlInfo)
			if err != nil {
				panic(err)
			}

			err = db.Ping()
			if err != nil {
				panic(err)
			}
		})
	})

	return db
}

func GetDBInstance() *sql.DB {
	return getDB()
}

func GetUsers() ([]users.IUser, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	var usersArr []users.IUser
	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var userpassword string
		var email string
		var phoneNum string
		var isAdmin bool
		err = rows.Scan(&id, &username, &userpassword, &email, &phoneNum, &isAdmin)
		if err != nil {
			// handle this error
			return nil, err
		}
		var user users.IUser
		if isAdmin {
			user = &users.Admin{Id: id, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}
		} else {
			user = &users.Customer{Id: id, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}

		}
		usersArr = append(usersArr, user)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return usersArr, nil
}

// LoginUser checks user credentials against the database.
func LoginUser(username, password string) (userID, balance int, isAdmin bool, err error) {
	var hashedPassword string

	// Query the database for the hashed password and admin flag based on the username
	err = db.QueryRow("SELECT user_id, balance, password, admin FROM users WHERE username = $1", username).Scan(&userID, &balance, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, false, errors.New("user not found")
		}
		return 0, 0, false, err
	}

	// Compare the hashed password from the database with the one the user provided.
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return 0, 0, false, errors.New("invalid password")
	}

	return userID, balance, isAdmin, nil
}

// RegisterUser adds a new user to the database.
func RegisterUser(username, password, email, phoneNum string, admin bool, balance int) error {
	admin = false
	// You would hash the password before storing it in your database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Replace with your actual database insertion logic.
	_, err = db.Exec("INSERT INTO users (username, password, email, phone_num, admin, balance) VALUES ($1, $2, $3, $4, $5, $6)", username, string(hashedPassword), email, phoneNum, admin, balance)
	if err != nil {
		return err
	}

	return nil
}

func InsertUser(username, password, email, phoneNum string, admin bool) error {
	var user users.IUser

	if admin {
		user = &users.Admin{
			UserName:     username,
			UserPassword: password,
			Email:        email,
			PhoneNum:     phoneNum,
		}
	} else {
		user = &users.Customer{
			UserName:     username,
			UserPassword: password,
			Email:        email,
			PhoneNum:     phoneNum,
		}
	}

	// Now use InsertUser function to insert the user into the DB.
	err := UserInsert(user)
	if err != nil {
		return err
	}

	return nil
}

// InsertUser inserts a user into the database.
func UserInsert(user users.IUser) error {
	db := GetDBInstance()

	sqlInsertUser := `
	INSERT INTO users (username, password, email, phone_num, admin)
	VALUES ($1, $2, $3, $4, $5)`
	var errUser error

	switch u := user.(type) {
	case *users.Admin:
		_, errUser = db.Exec(sqlInsertUser, u.UserName, u.UserPassword, u.Email, u.PhoneNum, true)

	case *users.Customer:
		_, errUser = db.Exec(sqlInsertUser, u.UserName, u.UserPassword, u.Email, u.PhoneNum, false)

	default:
		return errors.New("invalid user type")
	}

	if errUser != nil {
		return errUser
	}
	return nil
}

func BuyProducts(userID int) (int, error) {
	db := GetDBInstance() // Get your DB instance.

	// Start a database transaction.
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Flag to check if the commit was successful.
	commitSuccess := false
	defer func() {
		if !commitSuccess {
			tx.Rollback()
		}
	}()

	// Calculate the total cost of the items in the cart.
	var totalCost int
	err = tx.QueryRow("SELECT SUM(total_price) FROM cart WHERE user_id = $1", userID).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	// Retrieve the current balance of the user.
	var currentBalance int
	err = tx.QueryRow("SELECT balance FROM users WHERE user_id = $1", userID).Scan(&currentBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user balance: %w", err)
	}

	// Check if the user has enough balance to cover the purchase.
	if currentBalance < totalCost {
		return 0, errors.New("insufficient balance")
	}

	// Update the user's balance.
	newBalance := currentBalance - totalCost
	_, err = tx.Exec("UPDATE users SET balance = $1 WHERE user_id = $2", newBalance, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to update user balance: %w", err)
	}

	_, err = tx.Exec("DELETE FROM cart WHERE user_id = $1", userID)
	if err != nil {
		return 0, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	commitSuccess = true

	return newBalance, nil
}

func CheckUser(username string, userpassword string) users.IUser {
	db := GetDBInstance()

	sqlCheckUser := `
	SELECT *
	FROM users
	WHERE username = $1 AND password = $2`
	var userId int
	var phoneNum string
	var email string
	var isAdmin bool

	row := db.QueryRow(sqlCheckUser, username, userpassword)
	switch err := row.Scan(&userId, &username, &userpassword, &email, &phoneNum, &isAdmin); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil
	case nil:
		if isAdmin {
			return &users.Admin{Id: userId, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}
		} else {
			return &users.Customer{Id: userId, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}
		}
	default:
		panic(err)
	}
}

func GetUserById(id int) (users.IUser, error) {
	db := GetDBInstance()

	sqlGetUserById := `
	SELECT *
	FROM users
	WHERE user_id = $1`
	var userId int
	var username string
	var userpassword string
	var email string
	var phoneNum string
	var isAdmin bool

	row := db.QueryRow(sqlGetUserById, id)
	switch err := row.Scan(&userId, &username, &userpassword, &email, &phoneNum, &isAdmin); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, err
	case nil:
		if isAdmin {
			return &users.Admin{Id: userId, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}, nil
		} else {
			return &users.Customer{Id: userId, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}, nil
		}
	default:
		panic(err)
	}
}

func UpdateUser(user users.IUser) (users.IUser, error) {
	db := GetDBInstance()
	sqlUpdateUser := `UPDATE users SET username = $2, password = $3, email = $4, phone_num = $5, admin = $6 WHERE user_id = $1`
	id, username, password, email, phoneNum := user.GetDetails()
	var isAdmin bool
	switch user.(type) {
	case *users.Admin:
		isAdmin = true
	case *users.Customer:
		isAdmin = false
	}

	_, err := db.Exec(sqlUpdateUser, id, username, password, email, phoneNum, isAdmin)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetProductById(id int) (products.Product, error) {
	db := GetDBInstance()
	var name string
	var desc string
	var price int
	sqlGetUserById := `
	SELECT *
	FROM products
	WHERE id = $1`

	row := db.QueryRow(sqlGetUserById, id)
	err := row.Scan(&name, &desc, &price)
	if err != nil {
		return products.Product{}, err
	}
	product := products.Product{Id: id, Name: name, Desc: desc, Price: price}
	return product, nil
}

func GetProducts() ([]products.Product, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	var productsArr []products.Product
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var desc string
		var price int
		err = rows.Scan(&id, &name, &desc, &price)
		if err != nil {
			// handle this error
			return nil, err
		}
		product := products.Product{
			Id:    id,
			Name:  name,
			Desc:  desc,
			Price: price,
		}
		productsArr = append(productsArr, product)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return productsArr, nil
}

func InsertProduct(product products.Product) (products.Product, error) {
	db := GetDBInstance()

	sqlInsertProduct := `
	INSERT INTO products (name, description, price)
	VALUES ($1, $2, $3) RETURNING id`

	row := db.QueryRow(sqlInsertProduct, product.Name, product.Desc, product.Price)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return products.Product{}, err
	}

	for category_id, _ := range product.Category {
		err := AddCategoryToProduct(id, category_id)
		if err != nil {
			return products.Product{}, err
		}
	}

	product.Id = id
	return product, nil
}

func AddToCart(userID, productID, quantity int) error {
	db := GetDBInstance() // Get your DB instance.

	// Retrieve the price of the product.
	var price float64
	err := db.QueryRow("SELECT price FROM products WHERE id = $1", productID).Scan(&price)
	if err != nil {
		return err // Return error if product does not exist or query failed.
	}

	// Calculate the total price.
	totalPrice := price * float64(quantity)

	// Prepare statement for inserting data into the cart table.
	stmt, err := db.Prepare("INSERT INTO cart (user_id, product_id, quantity, price, total_price) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the user's data.
	_, err = stmt.Exec(userID, productID, quantity, price, totalPrice)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(product products.Product) (products.Product, error) {
	db := GetDBInstance()

	sqlUpdateProduct := `
	UPDATE products
	SET name = $2, description = $3, price = $4
	WHERE id = $1`

	_, err := db.Exec(sqlUpdateProduct, product.Id, product.Name, product.Desc, product.Price)
	if err != nil {
		return products.Product{}, err
	}
	var categoriesToDelete []int
	var categoriesToAdd []int
	oldCategories, err := GetProductCategoriesIds(product.Id)
	if err != nil {
		return products.Product{}, err
	}
	for category_id, _ := range product.Category {
		if !contains(oldCategories, category_id) {
			categoriesToAdd = append(categoriesToAdd, category_id)
		}
	}
	for _, category_id := range oldCategories {
		if !contains(categoriesToAdd, category_id) {
			categoriesToDelete = append(categoriesToDelete, category_id)
		}
	}
	for _, category_id := range categoriesToAdd {
		err := AddCategoryToProduct(product.Id, category_id)
		if err != nil {
			return products.Product{}, err
		}
	}
	for _, category_id := range categoriesToDelete {
		err := DeleteCategoryFromProductById(category_id)
		if err != nil {
			return products.Product{}, err
		}
	}

	return product, nil
}

func DeleteProduct(id int) error {
	db := GetDBInstance()

	sqlDeleteProduct := `
	DELETE FROM products
	WHERE id = $1`

	_, err := db.Exec(sqlDeleteProduct, id)
	if err != nil {
		return err
	}
	return nil
}

// category
func GetCategories() ([]string, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	var categoriesArr []string
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		err = rows.Scan(&id, &name)
		if err != nil {
			// handle this error
			return nil, err
		}
		categoriesArr = append(categoriesArr, name)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categoriesArr, nil
}

func GetCategoriesMap() (map[int]string, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	categoriesMap := make(map[int]string)
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		err = rows.Scan(&id, &name)
		if err != nil {
			// handle this error
			return nil, err
		}
		categoriesMap[id] = name
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categoriesMap, nil
}

func InsertCategory(category string) error {
	db := GetDBInstance()

	sqlInsertCategory := `
	INSERT INTO categories (category_name)
	VALUES ($1)`

	_, err := db.Exec(sqlInsertCategory, category)
	if err != nil {
		return err
	}
	return nil
}

func GetCategoryId(category string) (int, error) {
	db := GetDBInstance()

	sqlGetCategoryId := `
	SELECT category_id
	FROM categories
	WHERE category_name = $1`

	row := db.QueryRow(sqlGetCategoryId, category)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetProductCategoriesIds(product_id int) ([]int, error) {
	db := GetDBInstance()
	sqlProductCategories := `
		select c.category_id, c.category_name 
		from products_categories as pc
		inner join products as p
		on pc.product_id = p.id
		inner join categories as c
		on c.category_id = pc.category_id
		where p.id = $1;`
	rows, err := db.Query(sqlProductCategories, product_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categoriesIds []int
	for rows.Next() {
		var categoryId int
		err = rows.Scan(&categoryId)

		if err != nil {
			return nil, err
		}
		categoriesIds = append(categoriesIds, categoryId)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categoriesIds, nil
}

func GetProductCategoriesMap(product_id int) (map[int]string, error) {
	db := GetDBInstance()
	sqlProductCategories := `
		select c.category_id,c.category_name 
		from products_categories as pc
		inner join products as p
		on pc.product_id = p.id
		inner join categories as c
		on c.category_id = pc.category_id
		where p.id = $1;`
	rows, err := db.Query(sqlProductCategories, product_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categoriesMap := make(map[int]string)
	for rows.Next() {
		var categoryId int
		var categoryName string
		err = rows.Scan(&categoryId, &categoryName)

		if err != nil {
			return nil, err
		}
		categoriesMap[categoryId] = categoryName
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categoriesMap, nil
}

func AddCategoriesToProduct(categories_id []int, product_id int) error {
	for _, category_id := range categories_id {
		err := AddCategoryToProduct(product_id, category_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddCategoryToProduct(productId int, categoryId int) error {
	db := GetDBInstance()

	sqlInsertCategory := `
	INSERT INTO products_categories (product_id, category_id)
	VALUES ($1, $2)`

	_, err := db.Exec(sqlInsertCategory, productId, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategoryFromProductByName(category string) error {
	db := GetDBInstance()

	sqlDeleteCategory := `
	DELETE FROM products_categories WHERE category_id = $1 AND product_id = $2`

	_, err := db.Exec(sqlDeleteCategory, category)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategoryFromProductById(categoryId int) error {
	db := GetDBInstance()

	sqlDeleteCategory := `
	DELETE FROM products_categories WHERE category_id = $1 AND product_id = $2`

	_, err := db.Exec(sqlDeleteCategory, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func contains(arr []int, elem int) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123412"
	dbname   = "db_shop"
)
