package notifications

import (
	"fmt"
	db "shop/db_f"
	"shop/new_users"
)

type Observer interface {
	Update() error
}

type Subject interface {
	RegisterObserver(observer Observer)
	RemoveObserver(observer Observer)
	NotifyObservers()
}

type UserSubject struct {
	observers []Observer
}

func (u *UserSubject) GetObservers() []Observer {
	return u.observers
}

func (u *UserSubject) RegisterAllUsers() error {
	db := db.GetDBInstance()
	rows, err := db.Query("SELECT user_id, email, username FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var username string
		if err := rows.Scan(&id, &email, &username); err != nil {
			return err
		}
		u.RegisterObserver(&new_users.User{UserID: id, Email: email, Username: username})
	}
	return nil
}

func (u *UserSubject) RegisterObserver(observer Observer) {
	u.observers = append(u.observers, observer)
}

func (u *UserSubject) RemoveObserver(observer Observer) {
	for i, obs := range u.observers {
		if obs == observer {
			u.observers = append(u.observers[:i], u.observers[i+1:]...)
			break
		}
	}
}

func (u *UserSubject) NotifyObservers() {
	for _, observer := range u.observers {
		user := observer.(*new_users.User) // type assertion to get the user object from the observer interface
		err := user.Update()               // call the update method on the user object
		if err != nil {
			fmt.Println("Error updating user:", err)
		}
	}
}
