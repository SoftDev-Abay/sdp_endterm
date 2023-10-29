package db

import (
	"database/sql"
	"fmt"
	"sync"

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

func InsertUser(user users.IUser) {
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
		panic("Invalid user type")
	}

	if err_user != nil {
		panic(err_user)
	}
}

// func CheckUser(user users.IUser) {
// 	db := GetDBInstance()

// 	sqlCheckUser := `
// 	SELECT username, password
// 	FROM users
// 	WHERE username = $1 AND password = $2`
// 	err_user := error(nil)

// 	switch u := user.(type) {
// 	case *users.Admin:
// 		_, err_user = db.Exec(sqlCheckUser, u.UserName, u.UserPassword)

// 	case *users.Customer:
// 		_, err_user = db.Exec(sqlCheckUser, u.UserName, u.UserPassword)

// 	default:
// 		panic("Invalid user type")
// 	}

// 	if err_user != nil {
// 		panic(err_user)
// 	}

// }

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "031216551248"
	dbname   = "db_shop"
)
