// +build usebuiltin

package conversions

import "github.com/uluyol/vplots/conversions/internal/builtin"

func ToSVG(pdfPath string) ([]byte, error) {
	return builtin.ToSVG(pdfPath)
}

func ToPNG(pdfPath string) ([]byte, error) {
	return builtin.ToPNG(pdfPath)
}
