package domain

type Product struct {
	Id, Stock, Price int
	Discount         *int
	Name, Barcode    string
}
