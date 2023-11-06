package new_users

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	db "shop/db_f"
)

type User struct {
	Username    string
	Password    string
	Email       string
	PhoneNum    string
	Admin       bool
	Balance     int
	Permissions iPermissionStrategy
}

func RegisterUser(user User) error { // Added isAdmin parameter.
	db := db.GetDBInstance()

	// Hash the password before storing it in your database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if user.Admin {
		user.Permissions = &AdminPermissions{}
	} else {
		// Assuming you have a similar struct for regular users
		user.Permissions = &UserPermissions{}
	}

	// Insert the new user into the database.
	_, err = db.Exec("INSERT INTO users (username, password, email, phone_num, admin, balance) VALUES ($1, $2, $3, $4, $5, $6)", user.Username, string(hashedPassword), user.Email, user.PhoneNum, user.Admin, user.Balance)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(username, password string) (userID int, balance int, permissions iPermissionStrategy, err error) {
	var hashedPassword string
	var isAdmin bool
	db := db.GetDBInstance()

	// Query the database for the hashed password and admin flag based on the username
	err = db.QueryRow("SELECT user_id, balance, password, admin FROM users WHERE username = $1", username).Scan(&userID, &balance, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil, errors.New("user not found")
		}
		return 0, 0, nil, err
	}

	// Compare the hashed password from the database with the one the user provided.
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return 0, 0, nil, errors.New("invalid password")
	}

	// Set permissions based on the isAdmin value
	if isAdmin {
		permissions = &AdminPermissions{}
	} else {
		// Assuming you have a regular user permissions strategy
		permissions = &UserPermissions{}
	}

	return userID, balance, permissions, nil
}
