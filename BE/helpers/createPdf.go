package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func CreatePdf(text string, code string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Summary")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, text, "", "", false)

	summaryDir := "./summary"
	if _, err := os.Stat(summaryDir); os.IsNotExist(err) {
		err := os.Mkdir(summaryDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	summaryPath := filepath.Join(summaryDir, fmt.Sprintf("%s-summary.pdf", code))
	err := pdf.OutputFileAndClose(summaryPath)
	if err != nil {
		return err
	}

	return nil
}
