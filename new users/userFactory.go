package new_users

// User структура с полем для стратегии
type User struct {
	Username    string
	Password    string
	Permissions iPermissionStrategy
}

// Фабрика для создания пользователей с разными правами
func UserInsert(username, password string, admin bool) *User {
	user := &User{Username: username, Password: password}
	if admin {
		user.Permissions = &AdminPermissions{}
	} else {
		user.Permissions = &UserPermissions{}
	}
	return user
}
