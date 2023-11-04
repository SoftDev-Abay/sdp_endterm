package new_users

type AdminPermissions struct{}

func (ap *AdminPermissions) CanAddProduct() bool {
	return true
}

func (ap *AdminPermissions) CanDeleteProduct() bool {
	return true
}

func (ap *AdminPermissions) CanBuyProduct() bool {
	return false
}

func (ap *AdminPermissions) CanViewProduct() bool {
	return true
}
