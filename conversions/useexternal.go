// +build !cgo useexternal !linux,!darwin !amd64

package conversions

type mConverter struct{}

func (c *Converter) mInit()  {}
func (c *Converter) mClose() {}

func (c *Converter) ToSVG(pdfPath string) ([]byte, error) {
	return externalToSVG(pdfPath)
}

func (c *Converter) ToPNG(pdfPath string, width int, logAttr string) ([]byte, error) {
	return c.externalToPNG(pdfPath, width, logAttr)
}

func AppMain(mainFunc func()) { mainFunc() }
