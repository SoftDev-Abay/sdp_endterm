package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"shop/products"
	"shop/users"

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

func InsertUser(user users.IUser) error {
	db := GetDBInstance()

	sqlInsertUser := `
	INSERT INTO users (username, password, email, phone_num, admin)
	VALUES ($1, $2, $3, $4, $5)`
	err_user := error(nil)

	switch u := user.(type) {
	case *users.Admin:
		_, err_user = db.Exec(sqlInsertUser, u.UserName, u.UserPassword, u.Email, u.PhoneNum, true)

	case *users.Customer:
		_, err_user = db.Exec(sqlInsertUser, u.UserName, u.UserPassword, u.Email, u.PhoneNum, false)

	default:
		return errors.New("invalid user type")
	}

	if err_user != nil {
		return err_user
	}
	return nil
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
