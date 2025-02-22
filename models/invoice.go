package models

import (
	"time"
)

// Issuer represents the entity issuing the invoice (your platform or freelancer)
type Issuer struct {
	Name         string
	Street       string
	StreetNumber string
	ZipCode      string
	City         string
	Country      string
	NIF          string // Tax Identification Number
	IBAN         string // Required for bank transfers
	Email        string // Useful for digital invoices
	Phone        string // Optional
}

// Client represents the recipient of the invoice
type Client struct {
	Name         string
	Street       string
	StreetNumber string
	ZipCode      string
	City         string
	Country      string
	NIF          string // Required for B2B invoices
	Email        string // For sending invoices
}

// Item represents an invoice line item
type Item struct {
	Description string  // Service/Product description
	Quantity    float64 // Can be hours, units, etc.
	UnitPrice   float64 // Price per unit before VAT
	VATRate     float64 // VAT percentage (e.g., 21.0 for 21%)
	VATAmount   float64 // VAT amount for this item
	Total       float64 // Total amount (UnitPrice * Quantity + VAT)
}

// Invoice represents the full invoice document
type Invoice struct {
	Number        string    // Unique invoice number
	Date          time.Time // Invoice issue date
	DueDate       time.Time // Payment due date
	Issuer        Issuer    // Issuing entity
	Client        Client    // Recipient
	Items         []Item    // List of invoice items
	TotalNet      float64   // Total before VAT
	TotalVAT      float64   // Total VAT amount
	TotalGross    float64   // Total after VAT
	PaymentMethod string    // Bank Transfer, PayPal, etc.
	Notes         string    // Additional information
	InvoiceType   string    // "Standard", "Rectifying", etc.
	Reference     string    // Reference to previous invoice if rectifying
}
