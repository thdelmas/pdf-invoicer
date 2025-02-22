package pdfinvoicer

import "fmt"

func (i *Invoice) validate() error {
	if i.Number == "" {
		return fmt.Errorf("invoice number is required")
	}
	if i.Date.IsZero() {
		return fmt.Errorf("invoice date is required")
	}
	if i.Issuer.Name == "" {
		return fmt.Errorf("issuer name is required")
	}
	if i.Client.Name == "" {
		return fmt.Errorf("client name is required")
	}
	if len(i.Items) == 0 {
		return fmt.Errorf("invoice must have at least one item")
	}
	return nil
}
