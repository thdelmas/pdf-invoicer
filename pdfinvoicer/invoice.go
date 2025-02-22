package pdfinvoicer

import (
	"time"

	"github.com/thdelmas/pdf-invoicer/invoicertypes"
)

type Invoice struct {
	invoicertypes.Invoice
}

func NewInvoice(number string) *Invoice {
	return &Invoice{
		Invoice: invoicertypes.Invoice{
			Number: number,
			Date:   time.Now(),
			Items:  make([]invoicertypes.Item, 0),
		},
	}
}

func (i *Invoice) AddItem(item invoicertypes.Item) {
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
