package db

import (
	"database/sql"
	"errors"
	"fmt"
	"shop/products"
	"sync"

	_ "github.com/lib/pq"
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

func BuyProducts(userID int) (int, error) {
	db := GetDBInstance() // Get your DB instance.

	// start database transaction.
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// flag to check if the commit was successful.
	commitSuccess := false
	defer func() {
		if !commitSuccess {
			tx.Rollback()
		}
	}()

	// calculate the total cost of the items in the cart
	var totalCost int
	err = tx.QueryRow("SELECT SUM(total_price) FROM cart WHERE user_id = $1", userID).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	// retrieve the current balance of the user
	var currentBalance int
	err = tx.QueryRow("SELECT balance FROM users WHERE user_id = $1", userID).Scan(&currentBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user balance: %w", err)
	}

	// check if the user has enough balance to cover the purchase
	if currentBalance < totalCost {
		return 0, errors.New("insufficient balance")
	}

	// update the user's balance
	newBalance := currentBalance - totalCost
	_, err = tx.Exec("UPDATE users SET balance = $1 WHERE user_id = $2", newBalance, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to update user balance: %w", err)
	}

	_, err = tx.Exec("DELETE FROM cart WHERE user_id = $1", userID)
	if err != nil {
		return 0, fmt.Errorf("failed to clear cart: %w", err)
	}

	// commit the transaction.
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	commitSuccess = true

	return newBalance, nil
}

func GetUserId(username string) (int, error) {
	db := GetDBInstance()

	sqlGetUserById := `
	SELECT user_id
	FROM users
	WHERE username = $1`
	var id int
	row := db.QueryRow(sqlGetUserById, username)

	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
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

	CreateNotification(fmt.Sprintf("New product was added to the shop! Product name: %v", product.Name))

	product.Id = id
	return product, nil
}

func CreateNotification(text string) error {
	db := GetDBInstance()

	sqlInsert := `
	INSERT INTO notifications(text) VALUES($1) RETURNING id`

	_, err := db.Exec(sqlInsert, text)
	if err != nil {
		return err
	}
	return nil

}

func GetNotifications() (map[int]string, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT id,text FROM notifications")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	notifications := make(map[int]string)
	defer rows.Close()
	for rows.Next() {
		var notification_id int
		var text string
		err = rows.Scan(&notification_id, &text)
		if err != nil {
			// handle this error
			return nil, err
		}
		notifications[notification_id] = text
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func AddNotificationToUser(userID, notificationID int) error {
	db := GetDBInstance()

	_, err := db.Exec("INSERT INTO users_notifications(user_id, notification_id) VALUES($1, $2)", userID, notificationID)
	if err != nil {
		return err
	}
	return nil
}

// Function to get notifications for a user
func GetNotificationsForUserByID(userID int) (map[int]string, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT nu.id, n.text FROM notifications n INNER JOIN users_notifications nu ON n.id = nu.notification_id WHERE nu.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make(map[int]string)
	for rows.Next() {
		var text string
		var id int
		if err := rows.Scan(&id, &text); err != nil {
			return nil, err
		}
		notifications[id] = text
	}
	return notifications, nil
}

func GetNotificationsForUserByUsername(username string) (map[int]string, error) {
	db := GetDBInstance()
	rows, err := db.Query("SELECT nu.id, n.text FROM notifications n INNER JOIN users_notifications nu ON n.id = nu.notification_id INNER JOIN users u ON u.id = nu.user_id WHERE u.username = $1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make(map[int]string)
	for rows.Next() {
		var text string
		var id int
		if err := rows.Scan(&id, &text); err != nil {
			return nil, err
		}
		notifications[id] = text
	}
	return notifications, nil
}

func ClearAllNotifications() error {
	db := GetDBInstance()
	_, err := db.Exec("DELETE FROM notifications")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users_notifications")
	if err != nil {
		return err
	}
	return nil
}

// Function to mark a notification as seen by a user
// func MarkNotificationAsSeen(userID, notificationID int) error {
// 	db := GetDBInstance()
// 	_, err := db.Exec("UPDATE users_notifications SET seen = true WHERE user_id = $1 AND notification_id = $2", userID, notificationID)
// 	return err
// }

func AddToCart(userID, productID, quantity int) error {
	db := GetDBInstance()

	// retrieve the price of the product.
	var price float64
	err := db.QueryRow("SELECT price FROM products WHERE id = $1", productID).Scan(&price)
	if err != nil {
		return err // Return error if product does not exist or query failed.
	}

	// calculate the total price.
	totalPrice := price * float64(quantity)

	// prepare statement for inserting data into the cart table.
	stmt, err := db.Prepare("INSERT INTO cart (user_id, product_id, quantity, price, total_price) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// execute the prepared statement with the user's data.
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
	password = "031216551248"
	dbname   = "db_shop"
)
