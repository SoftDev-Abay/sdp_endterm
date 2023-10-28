package products

type product struct {
	id       int
	name     string
	price    string
	category []string
}

func (p *product) SetName(name string) {
	p.name = name
}

func (p *product) SetPrice(price string) {
	p.price = price
}

func (p *product) GetDetails() (string, string) {
	return p.name, p.price
}
