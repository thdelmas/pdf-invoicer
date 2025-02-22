package pdf_invoicer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thdelmas/pdf-invoicer/models"
)

func GenerateInvoice(invoice models.Invoice, outputPath string) error {
	// Ensure the output path is valid
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Convert invoice struct to JSON (can be adapted for PDF generation)
	invoiceData, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize invoice: %v", err)
	}

	// Write to file
	err = os.WriteFile(outputPath, invoiceData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write invoice file: %v", err)
	}

	fmt.Printf("Invoice %s generated successfully at %s\n", invoice.Number, outputPath)
	return nil
}
