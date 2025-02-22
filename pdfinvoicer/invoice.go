package pdfinvoicer

import "github.com/thdelmas/pdf-invoicer/models"

func NewInvoice() models.Invoice {
	return models.Invoice{
		Issuer: NewIssuer(),
		Client: NewClient(),
	}
}
