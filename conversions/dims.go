package conversions

import (
	"sync"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

type dim = struct {
	wInches float64
	hInches float64
}

type Converter struct {
	mu       sync.Mutex
	dimCache map[string]dim
}

func NewConverter() *Converter {
	return &Converter{
		dimCache: make(map[string]dim),
	}
}

func (c *Converter) getDim(pdf string) dim {
	c.mu.Lock()
	v, ok := c.dimCache[pdf]
	c.mu.Unlock()
	if ok {
		return v
	} else {
		dims, err := pdfcpuapi.PageDimsFile(pdf)
		if err == nil && len(dims) > 0 {
			gotDim := dims[0].ToInches()
			d := dim{
				wInches: gotDim.Width,
				hInches: gotDim.Height,
			}
			c.mu.Lock()
			c.dimCache[pdf] = d
			c.mu.Unlock()
			return d
		}
	}
	return dim{1, 1} // just use some default
}
