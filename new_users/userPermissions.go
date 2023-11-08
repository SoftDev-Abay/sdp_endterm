package new_users

type UserPermissions struct{}

func (up *UserPermissions) CanAddProduct() bool {
	return false
}

func (up *UserPermissions) CanDeleteProduct() bool {
	return false
}

func (up *UserPermissions) CanBuyProduct() bool {
	return true
}

func (up *UserPermissions) CanViewProduct() bool {
	return true
}

func (up *UserPermissions) CanAddNotification() bool {
	return false
}
