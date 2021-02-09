// +build !usebuiltin

package conversions

func (c *Converter) ToSVG(pdfPath string) ([]byte, error) {
	return externalToSVG(pdfPath)
}

func (c *Converter) ToPNG(pdfPath string, width int, logAttr string) ([]byte, error) {
	return c.externalToPNG(pdfPath, width, logAttr)
}
