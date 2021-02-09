// +build linux,amd64 darwin,amd64
// +build !useexternal
// +build cgo

package conversions

import (
	"github.com/uluyol/vplots/conversions/internal/builtin"
)

type mConverter struct {
	builtinC builtin.Converter
}

func (c *Converter) mInit()  {}
func (c *Converter) mClose() {}

func (c *Converter) ToSVG(pdfPath string) ([]byte, error) {
	return externalToSVG(pdfPath)
}

func (c *Converter) ToPNG(pdfPath string, width int, logAttr string) ([]byte, error) {
	return c.builtinC.ToPNG(pdfPath, width)
}

func AppMain(mainFunc func()) { builtin.AppMain(mainFunc) }
