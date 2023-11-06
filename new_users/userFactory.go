package new_users

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	db "shop/db_f"
)

// Define the IUserFactory interface
type IUserFactory interface {
	CreateUser(username, password, email, phoneNum string, balance int) (User, error)
}

// Define the AdminUserFactory which will create users with admin permissions
type AdminUserFactory struct{}

func (f *AdminUserFactory) CreateUser(username, password, email, phoneNum string, balance int) (User, error) {
	user := User{
		Username:    username,
		Password:    password,
		Email:       email,
		PhoneNum:    phoneNum,
		Admin:       true,
		Balance:     balance,
		Permissions: &AdminPermissions{},
	}
	return user, nil
}

// Define the RegularUserFactory which will create users with regular permissions
type RegularUserFactory struct{}

func (f *RegularUserFactory) CreateUser(username, password, email, phoneNum string, balance int) (User, error) {
	user := User{
		Username:    username,
		Password:    password,
		Email:       email,
		PhoneNum:    phoneNum,
		Admin:       false,
		Balance:     balance,
		Permissions: &UserPermissions{},
	}
	return user, nil
}

// User struct represents a user in the system.
type User struct {
	UserID      int
	Username    string
	Password    string
	Email       string
	PhoneNum    string
	Admin       bool
	Balance     int
	Permissions iPermissionStrategy
}

// The Register function takes a factory, which will provide the mechanism to create a User with the correct permissions.
func Register(factory IUserFactory, username, password, email, phoneNum string, balance int) error {
	user, err := factory.CreateUser(username, password, email, phoneNum, balance)
	if err != nil {
		return err
	}

	// Hash the password before storing it in your database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert the new user into the database.
	db := db.GetDBInstance()
	_, err = db.Exec("INSERT INTO users (username, password, email, phone_num, admin, balance) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Username, string(hashedPassword), user.Email, user.PhoneNum, user.Admin, user.Balance)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(username, password string) (User, error) {
	var user User
	var hashedPassword string
	var isAdmin bool
	dbInstance := db.GetDBInstance()

	// Query the database for the hashed password and admin flag based on the username
	err := dbInstance.QueryRow("SELECT user_id, balance, password, admin FROM users WHERE username = $1", username).Scan(&user.UserID, &user.Balance, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	// Compare the hashed password from the database with the one the user provided.
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return User{}, errors.New("invalid password")
	}

	// Set permissions based on the isAdmin value
	if isAdmin {
		user.Permissions = &AdminPermissions{}
	} else {
		user.Permissions = &UserPermissions{}
	}

	return user, nil
}

func (u *User) HasAdminPermissions() bool {
	_, ok := u.Permissions.(*AdminPermissions)
	return ok
}
