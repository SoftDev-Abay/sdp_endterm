package commands

import "shop/db_f"

type DeleteProductCommand struct {
	ProductID int
}

func (c *DeleteProductCommand) Execute() error {
	return db.DeleteProduct(c.ProductID)
}
