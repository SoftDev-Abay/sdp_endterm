package commands

import db "shop/db_f"

type DeleteProductCommand struct {
	ProductID int
}

func (c *DeleteProductCommand) Execute() error {
	err := db.DeleteProduct(c.ProductID)
	return err
}
