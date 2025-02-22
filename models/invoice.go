package models

import "time"

type Invoice struct {
	Number        string
	Date          time.Time
	DueDate       time.Time
	Issuer        Issuer
	Client        Client
	Items         []Item
	TotalNet      float64
	TotalVAT      float64
	TotalGross    float64
	PaymentMethod string
	Notes         string
	InvoiceType   string
	Reference     string
}
