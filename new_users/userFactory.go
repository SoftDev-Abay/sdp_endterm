package new_users

import (
	"database/sql"
	"errors"
	db "shop/db_f"

	"golang.org/x/crypto/bcrypt"
)

type IUserFactory interface {
	CreateUser(username, password, email, phoneNum string, balance int) (User, error)
}

// AdminUserFactory will create users with admin permissions
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

// RegularUserFactory will create users with regular permissions
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

func (u *User) Update() error {
	userId, err := db.GetUserId(u.Username)

	if err != nil {
		return err
	}
	userNotifications, err := db.GetNotificationsForUserByID(userId)
	if err != nil {
		return err
	}
	allNotifications, err := db.GetNotifications()
	if err != nil {
		return err
	}
	for id, _ := range allNotifications {
		if !mapContains(userNotifications, id) {
			err = db.AddNotificationToUser(userId, id)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

// Register function takes a factory, which will provide the mechanism to create a User with the correct permissions.
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

	// query the database for the hashed password and admin flag based on the username
	err := dbInstance.QueryRow("SELECT user_id, balance, password, admin FROM users WHERE username = $1", username).Scan(&user.UserID, &user.Balance, &hashedPassword, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	// compare the hashed password from the database with the one the user provided.
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return User{}, errors.New("invalid password")
	}

	// set permissions based on the isAdmin value
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

func GetUsers() ([]User, error) {
	db := db.GetDBInstance()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	var usersArr []User
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
		user := User{
			UserID:   id,
			Username: username,
			Password: userpassword,
			Email:    email,
			PhoneNum: phoneNum,
			Admin:    isAdmin,
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

func mapContains(mapInput map[int]string, elem int) bool {
	for id, _ := range mapInput {
		if id == elem {
			return true
		}
	}
	return false
}
