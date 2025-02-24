package pdfinvoicer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-pdf/fpdf"
)

type Address struct {
	Street       string
	StreetNumber string
	Stairs       string
	Floor        string
	Door         string
	ZipCode      string
	City         string
	State        string
	Country      string
}

type Issuer struct {
	Name    string
	Address Address
	NIF     string
	IBAN    string
	Email   string
	Phone   string
}

type Client struct {
	Name    string
	Address Address
	NIF     string
}

type Item struct {
	Description string  // Description of the service or product
	Quantity    float64 // Quantity of items/services (if applicable)
	UnitPrice   float64 // Price per unit (if applicable)
	VATRate     float64 // VAT rate applied (21%, 10%, 4%, or 0 for exempt)
	VATAmount   float64 // Calculated VAT amount
	Total       float64 // Total amount including VAT
}

type Invoice struct {
	Number   string    // Unique invoice number
	EmitDate time.Time // Issuance date
	OpDate   time.Time // Operation date
	DueDate  time.Time // Due date

	Issuer Issuer // Issuing company (Sociedad Limitada)
	Client Client // Client receiving the invoice

	Items []Item // List of items in the invoice

	BaseAmount float64 // Taxable base amount before VAT
	VATAmount  float64 // Total VAT amount
	Total      float64 // Total payable amount (BaseAmount + VATAmount)

	Paid bool // Payment status (true if paid, false otherwise)

	Notes    string // Additional notes or comments
	Ref      string // Reference or internal tracking number
	LogoPath string // Path to the logo to include in the invoice
}

func checkAddress(address Address) error {
	if address.Street == "" {
		err := errors.New("Address street is mandatory")
		return err
	}

	if address.StreetNumber == "" {
		err := errors.New("Address street number is mandatory")
		return err
	}

	if address.ZipCode == "" {
		err := errors.New("Address zip code is mandatory")
		return err
	}

	if address.City == "" {
		err := errors.New("Address city is mandatory")
		return err
	}

	if address.Country == "" {
		err := errors.New("Address country is mandatory")
		return err
	}

	return nil
}

func checkIssuer(issuer Issuer) error {
	if issuer.Name == "" {
		err := errors.New("Issuer name is mandatory")
		return err
	}

	if issuer.NIF == "" {
		err := errors.New("Issuer NIF is mandatory")
		return err
	}

	// Check Address
	err := checkAddress(issuer.Address)
	if err != nil {
		return err
	}

	return nil
}

func checkClient(client Client) error {
	if client.Name == "" {
		err := errors.New("Client name is mandatory")
		return err
	}

	if client.NIF == "" {
		err := errors.New("Client NIF is mandatory")
		return err
	}

	// Check Address
	err := checkAddress(client.Address)
	if err != nil {
		return err
	}

	return nil
}

func NewAddress(street, streetNumber, stairs, floor, door, zipCode, city, state, country string) (Address, error) {
	address := Address{
		Street:       street,
		StreetNumber: streetNumber,
		Stairs:       stairs,
		Floor:        floor,
		Door:         door,
		ZipCode:      zipCode,
		City:         city,
		State:        state,
		Country:      country,
	}

	err := checkAddress(address)
	if err != nil {
		return Address{}, err
	}

	return address, nil
}

func NewIssuer(name string, address Address, nif, iban, email, phone string) (Issuer, error) {
	issuer := Issuer{
		Name:    name,
		Address: address,
		NIF:     nif,
		IBAN:    iban,
		Email:   email,
		Phone:   phone,
	}

	err := checkIssuer(issuer)
	if err != nil {
		return Issuer{}, err
	}

	return issuer, nil
}

func NewClient(name string, address Address, nif string) (Client, error) {
	client := Client{
		Name:    name,
		Address: address,
		NIF:     nif,
	}

	err := checkClient(client)
	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func NewItem(description string, quantity, unitPrice, vatRate float64) (Item, error) {
	if description == "" {
		err := errors.New("Item description is mandatory")
		return Item{}, err
	}

	if quantity < 0 {
		err := errors.New("Item quantity must be positive")
		return Item{}, err
	}

	if unitPrice < 0 {
		err := errors.New("Item unit price must be positive")
		return Item{}, err
	}

	if vatRate < 0 {
		err := errors.New("Item VAT rate must be positive")
		return Item{}, err
	}

	item := Item{
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		VATRate:     vatRate,
	}

	item.VATAmount = item.Quantity * item.UnitPrice * item.VATRate
	item.Total = item.Quantity*item.UnitPrice + item.VATAmount

	return item, nil
}

func NewInvoice(number string, emitDate, opDate, dueDate time.Time, issuer Issuer, client Client, items []Item, notes, ref string) (Invoice, error) {
	if number == "" {
		err := errors.New("Invoice number is mandatory")
		return Invoice{}, err
	}

	// check dates
	if emitDate.Before(opDate) {
		err := errors.New("emission date must not be after operation date")
		return Invoice{}, err
	}

	if opDate.After(dueDate) {
		err := errors.New("operation date must be before due date")
		return Invoice{}, err
	}

	// Check Items
	if len(items) == 0 {
		err := errors.New("Invoice must have at least one item")
		return Invoice{}, err
	}

	// Check Issuer
	err := checkIssuer(issuer)
	if err != nil {
		return Invoice{}, err
	}

	// Check Client
	err = checkClient(client)
	if err != nil {
		return Invoice{}, err
	}

	// Check Items
	for _, item := range items {
		_, err := NewItem(item.Description, item.Quantity, item.UnitPrice, item.VATRate)
		if err != nil {
			return Invoice{}, err
		}
	}

	baseAmount := 0.0
	vatAmount := 0.0
	total := 0.0

	for _, item := range items {
		baseAmount += item.Quantity * item.UnitPrice
		vatAmount += item.VATAmount
		total += item.Total
	}

	invoice := Invoice{
		Number:   number,
		EmitDate: emitDate,
		OpDate:   opDate,
		DueDate:  dueDate,

		Issuer: issuer,
		Client: client,

		Items: items,

		BaseAmount: baseAmount,
		VATAmount:  vatAmount,
		Total:      total,

		Paid: false,

		Notes: notes,
		Ref:   ref,
	}

	return invoice, nil
}

func formatAddress(address Address) string {
	var formatted string

	formatted += fmt.Sprintf("%s %s\n", address.Street, address.StreetNumber)
	if address.Stairs != "" {
		formatted += fmt.Sprintf("Stairs: %s\n", address.Stairs)
	}
	if address.Floor != "" {
		formatted += fmt.Sprintf("Floor: %s\n", address.Floor)
	}
	if address.Door != "" {
		formatted += fmt.Sprintf("Door: %s\n", address.Door)
	}
	formatted += fmt.Sprintf("%s %s\n", address.ZipCode, address.City)
	if address.State != "" {
		formatted += fmt.Sprintf("%s\n", address.State)
	}
	formatted += fmt.Sprintf("%s\n", address.Country)

	return formatted
}

/*func addHeader(pdf *fpdf.Fpdf, invoice Invoice) {
	// Add logo if path is provided
	if invoice.LogoPath != "" {
		pdf.ImageOptions(
			invoice.LogoPath,
			10, // x position
			10, // y position
			30, // width
			0,  // height (0 = auto-calculated)
			false,
			fpdf.ImageOptions{ImageType: "", ReadDpi: true},
			0,
			"",
		)
	}

	// Company name in header (to the right of logo)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(45, 15)
	pdf.CellFormat(100, 10, invoice.Issuer.Name, "0", 0, "L", false, 0, "")

	// Add contact info in smaller font
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(45, 25)
	pdf.CellFormat(100, 5, fmt.Sprintf("%s %s, %s %s",
		invoice.Issuer.Address.Street,
		invoice.Issuer.Address.StreetNumber,
		invoice.Issuer.Address.ZipCode,
		invoice.Issuer.Address.City),
		"0", 1, "L", false, 0, "")

	// Add horizontal line under header
	pdf.SetLineWidth(0.5)
	pdf.Line(10, 40, 200, 40)

	// Reset position for rest of document
	pdf.SetY(50)
}*/

func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f EUR", amount)
}

func (i *Invoice) GeneratePDF(outputPath string) error {
	log.Println("Generating PDF for invoice", i.Number)

	log.Println("Invoice number:", i.Number)
	log.Println("Emit date:", i.EmitDate)
	log.Println("Operation date:", i.OpDate)
	log.Println("Due date:", i.DueDate)

	log.Println("Issuer name:", i.Issuer.Name)
	log.Println("Issuer NIF:", i.Issuer.NIF)
	log.Println("Issuer IBAN:", i.Issuer.IBAN)
	log.Println("Issuer email:", i.Issuer.Email)
	log.Println("Issuer phone:", i.Issuer.Phone)

	log.Println("Issuer address:")
	log.Println("  Street:", i.Issuer.Address.Street)
	log.Println("  Street number:", i.Issuer.Address.StreetNumber)
	log.Println("  Stairs:", i.Issuer.Address.Stairs)
	log.Println("  Floor:", i.Issuer.Address.Floor)
	log.Println("  Door:", i.Issuer.Address.Door)
	log.Println("  Zip code:", i.Issuer.Address.ZipCode)
	log.Println("  City:", i.Issuer.Address.City)
	log.Println("  State:", i.Issuer.Address.State)
	log.Println("  Country:", i.Issuer.Address.Country)

	log.Println("Client name:", i.Client.Name)
	log.Println("Client NIF:", i.Client.NIF)
	log.Println("Client address:")
	log.Println("  Street:", i.Client.Address.Street)
	log.Println("  Street number:", i.Client.Address.StreetNumber)
	log.Println("  Stairs:", i.Client.Address.Stairs)
	log.Println("  Floor:", i.Client.Address.Floor)
	log.Println("  Door:", i.Client.Address.Door)
	log.Println("  Zip code:", i.Client.Address.ZipCode)
	log.Println("  City:", i.Client.Address.City)
	log.Println("  State:", i.Client.Address.State)
	log.Println("  Country:", i.Client.Address.Country)

	log.Println("Items:")
	for _, item := range i.Items {
		log.Println("  Description:", item.Description)
		log.Println("  Quantity:", item.Quantity)
		log.Println("  Unit price:", item.UnitPrice)
		log.Println("  VAT rate:", item.VATRate)
		log.Println("  VAT amount:", item.VATAmount)
		log.Println("  Total:", item.Total)
	}

	log.Println("Base amount:", i.BaseAmount)
	log.Println("VAT amount:", i.VATAmount)
	log.Println("Total:", i.Total)

	log.Println("Paid:", i.Paid)

	log.Println("Notes:", i.Notes)
	log.Println("Reference:", i.Ref)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Use a font that supports UTF-8
	pdf.AddUTF8Font("Arial", "", "font.ttf")
	pdf.SetFont("Arial", "", 12)

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "INVOICE", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Issuer and Client information
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(95, 7, "From:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, "To:", "0", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	issuerAddress := formatAddress(i.Issuer.Address)
	clientAddress := formatAddress(i.Client.Address)

	// Using MultiCell for addresses to handle multiple lines
	currentY := pdf.GetY()
	pdf.MultiCell(95, 5, issuerAddress, "0", "L", false)
	pdf.SetXY(105, currentY)
	pdf.MultiCell(95, 5, clientAddress, "0", "R", false)

	pdf.Ln(10)

	// Invoice details
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 7, fmt.Sprintf("Invoice Number: %s", i.Number), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("Date: %s", i.EmitDate), "0", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Amount breakdown table
	pdf.SetFont("Arial", "B", 10)
	cols := []float64{70, 40, 40, 40}
	headers := []string{"Description", "Base Amount", "VAT (21%)", "Total"}

	for i, header := range headers {
		pdf.CellFormat(cols[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Calculate amounts
	total := i.Total

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(cols[0], 7, "Professional Services", "1", 0, "L", false, 0, "")
	pdf.CellFormat(cols[1], 7, formatCurrency(i.BaseAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(cols[2], 7, formatCurrency(i.VATAmount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(cols[3], 7, formatCurrency(total), "1", 1, "R", false, 0, "")
	pdf.Ln(10)

	// Total amount
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 7, fmt.Sprintf("Total Amount: %s", formatCurrency(total)), "0", 1, "R", false, 0, "")

	// Payment terms
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 5, "Payment Terms: Due within 30 days\nPlease include invoice number in payment reference", "0", "L", false)

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		panic(err)
	}

	return nil
}
