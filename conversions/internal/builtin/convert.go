package builtin

/*

#include <stdint.h>
#include <stdlib.h>

typedef struct {
	uint8_t *data;
	int32_t width;
	int32_t height;
} OutImage;

void OutImage_Clear(OutImage *im);

typedef void XPDFConverter;

void XPDFConverter_GlobalInit();

const char *XPDFConverter_RenderPDF(
	char *filename, int width, OutImage *out);

#cgo CXXFLAGS: -I${SRCDIR}/xpdf_darwin_amd64/include -std=c++11
#cgo LDFLAGS: -L${SRCDIR}/xpdf_darwin_amd64/lib -lxpdf

*/
import "C"
import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"unsafe"
)

func init() {
	C.XPDFConverter_GlobalInit()
}

type Converter struct{}

func (c *Converter) ToPNG(pdfPath string, width int) ([]byte, error) {
	cPDFPath := C.CString(pdfPath)
	defer C.free(unsafe.Pointer(cPDFPath))

	var rawIm C.OutImage
	rawIm.data = nil

	errMesg := C.XPDFConverter_RenderPDF(cPDFPath, C.int(width), &rawIm)
	if errMesg != nil {
		return nil, errors.New(C.GoString(errMesg))
	}
	defer C.OutImage_Clear(&rawIm)

	// Otherwise buffer is valid
	im := image.NewRGBA(image.Rect(0, 0, int(rawIm.width), int(rawIm.height)))
	for i := uintptr(0); i < uintptr(len(im.Pix)); i++ {
		im.Pix[i] = *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(rawIm.data)) + i))
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, im)

	return buf.Bytes(), err
}
