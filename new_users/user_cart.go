package new_users

import (
	"shop/db_f"
)

// CartItem represents an item in the user's shopping cart.
type CartItem struct {
	ProductID   int
	ProductName string
	Quantity    int
	TotalPrice  int
}

// ViewCart retrieves the items in the user's cart from the database.
func ViewCart(userID int) ([]CartItem, error) {
	db := db.GetDBInstance()

	query := `
        SELECT p.id, p.name, c.quantity, c.total_price
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
		err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.TotalPrice)
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
