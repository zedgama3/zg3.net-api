package model

// Product represents a product in the system
type Product struct {
	ID    int
	Name  string
	Price float64
}

// NewProduct creates and returns a new Product instance
func NewProduct(id int, name string, price float64) *Product {
	return &Product{ID: id, Name: name, Price: price}
}
