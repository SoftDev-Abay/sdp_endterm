package new_users

type iPermissionStrategy interface {
	CanAddProduct() bool
	CanDeleteProduct() bool
	CanBuyProduct() bool
	CanViewProduct() bool
}
