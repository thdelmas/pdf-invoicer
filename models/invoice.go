package models

import "time"

type Invoice struct {
	Number        string
	EmissionDate  time.Time
	OperationDate time.Time
	DueDate       time.Time

	Issuer Issuer
	Client Client

	Items []Item

	TotalNet   float64
	TotalVAT   float64
	TotalGross float64

	PaymentMethod string
	Paid          bool

	Notes     string
	Reference string
}
