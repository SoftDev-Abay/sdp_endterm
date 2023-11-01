package products

type Product struct {
	Id       int
	Name     string
	Desc     string
	Price    int
	Category map[int]string
}

func (p *Product) SetName(name string) {
	p.Name = name
}

func (p *Product) SetPrice(price int) {
	p.Price = price
}

func (p *Product) GetDetails() (string, int, string, map[int]string) {
	return p.Name, p.Price, p.Desc, p.Category
}
