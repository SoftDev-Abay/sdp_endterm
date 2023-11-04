package new_users

import (
	"shop/db_f"
)

// CartItem represents an item in the user's shopping cart.
type CartItem struct {
	ProductID   int
	ProductName string
	Quantity    int
}

// ViewCart retrieves the items in the user's cart from the database.
func ViewCart(userID int) ([]CartItem, error) {
	db := db.GetDBInstance() // Get your DB instance.

	// Use the correct placeholder syntax for PostgreSQL.
	query := `
        SELECT p.id, p.name, c.quantity
        FROM cart c
        JOIN products p ON c.product_id = p.id
        WHERE c.user_id = $1
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}
