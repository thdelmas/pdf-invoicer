package invoice

import (
	"time"

	"github.com/thdelmas/pdf-invoicer/models"
)

type Invoice struct {
	models.Invoice
}

func NewInvoice(number string) *Invoice {
	return &Invoice{
		Invoice: models.Invoice{
			Number: number,
			Date:   time.Now(),
			Items:  make([]models.Item, 0),
		},
	}
}

func (i *Invoice) AddItem(item models.Item) {
	i.Items = append(i.Items, item)
	i.calculateTotals()
}

func (i *Invoice) calculateTotals() {
	var totalNet, totalVAT float64

	for _, item := range i.Items {
		totalNet += item.UnitPrice * item.Quantity
		totalVAT += item.VATAmount
	}

	i.TotalNet = totalNet
	i.TotalVAT = totalVAT
	i.TotalGross = totalNet + totalVAT
}
