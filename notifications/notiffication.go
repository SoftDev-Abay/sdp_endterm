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

type UserNotificationSubject struct {
	observers []Observer
}

func (u *UserNotificationSubject) AddNotification(text string) error {
	err := db.CreateNotification(text)
	if err != nil {
		return err
	}
	u.NotifyObservers()
	return nil
}

func (u *UserNotificationSubject) GetObservers() []Observer {
	return u.observers
}

func (u *UserNotificationSubject) RegisterAllUsers() error {
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

func (u *UserNotificationSubject) RegisterObserver(observer Observer) {
	u.observers = append(u.observers, observer)
}

func (u *UserNotificationSubject) RemoveObserver(observer Observer) {
	for i, obs := range u.observers {
		if obs == observer {
			u.observers = append(u.observers[:i], u.observers[i+1:]...)
			break
		}
	}
}

func (u *UserNotificationSubject) NotifyObservers() {
	for _, observer := range u.observers {
		user := observer.(*new_users.User)
		err := user.Update()
		if err != nil {
			fmt.Println("Error updating user:", err)
		}
	}
}
