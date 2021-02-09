package builtin

import (
	"errors"
	"unsafe"
)

/*

#include <stdint.h>
#include <stdlib.h>

struct OutBuf {
  char *data;
  int len;
  int cap;
};

char *ConvertPDF(char *filename, struct OutBuf *output, int format);

void ClearOutBuf(struct OutBuf *buf) {
	free(buf->data);
}

*/
import "C"

func convert(pdfPath string, format C.int) ([]byte, error) {
	cPDFPath := C.CString(pdfPath)
	defer C.free(unsafe.Pointer(cPDFPath))

	var buffer C.struct_OutBuf
	errMesg := C.ConvertPDF(cPDFPath, &buffer, format)
	if errMesg != nil {
		return nil, errors.New(C.GoString(errMesg))
	}

	// Otherwise buffer is valid
	data := C.GoBytes(unsafe.Pointer(buffer.data), buffer.len)
	C.ClearOutBuf(&buffer)

	return data, nil
}

func ToSVG(pdfPath string) ([]byte, error) {
	return convert(pdfPath, 0)
}

func ToPNG(pdfPath string) ([]byte, error) {
	return convert(pdfPath, 1)
}
