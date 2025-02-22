package pdfinvoicer

import (
	"encoding/json"
	"fmt"
	"os"
)

func (i *Invoice) GeneratePDF(outputPath string) error {
	if err := i.validate(); err != nil {
		return fmt.Errorf("invalid invoice data: %v", err)
	}

	invoiceData, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize invoice: %v", err)
	}

	if err := os.WriteFile(outputPath, invoiceData, 0644); err != nil {
		return fmt.Errorf("failed to write invoice file: %v", err)
	}

	fmt.Printf("Invoice %s generated successfully at %s\n", i.Number, outputPath)
	return nil
}
