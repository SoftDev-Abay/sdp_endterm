package products

type IProduct interface {
	SetName(name string)
	SetPrice(price string)

	GetDetails() string
}
