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
			user = users.Admin{Id: id, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}
		} else {
			user = users.Customer{Id: id, UserName: username, UserPassword: userpassword, Email: email, PhoneNum: phoneNum}

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
		return errors.New("Invalid user type")
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
			return &users.Admin{userId, username, userpassword, email, phoneNum}
		} else {
			return &users.Customer{userId, username, userpassword, email, phoneNum}
		}
	default:
		panic(err)
	}
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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "031216551248"
	dbname   = "db_shop"
)
