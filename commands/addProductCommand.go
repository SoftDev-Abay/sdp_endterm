package commands

import (
	"shop/db_f"
	"shop/products"
)

type AddProductCommand struct {
	Product products.Product
}

func (c *AddProductCommand) Execute() error {
	_, err := db.InsertProduct(c.Product)
	return err
}
