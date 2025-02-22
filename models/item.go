package models

type Item struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	VATRate     float64
	VATAmount   float64
	Total       float64
}
