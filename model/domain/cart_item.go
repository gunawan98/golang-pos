package domain

type CartItem struct {
	Id, CartID, ProductID, Quantity, UnitPrice, TotalPrice int
}

type CartItemWithProduct struct {
	Id, CartID, ProductID, Quantity, UnitPrice, TotalPrice int
	ProductName                                            string
}
