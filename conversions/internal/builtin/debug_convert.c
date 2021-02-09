// +build ignore

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

typedef struct {
  uint8_t *data;
  int32_t width;
  int32_t height;
} OutImage;

void OutImage_Clear(OutImage *im);

typedef void XPDFConverter;

void XPDFConverter_GlobalInit();

const char *XPDFConverter_RenderPDF(char *filename, int width, OutImage *out);

int main(int argc, char **argv) {
  if (argc != 2) {
    fprintf(stderr, "usage: debug_convert input.pdf > out.bitmap\n");
    return 2;
  }

  XPDFConverter_GlobalInit();

  OutImage out;
  out.data = NULL;
  const char *err = XPDFConverter_RenderPDF(argv[1], 100, &out);
  if (err != NULL) {
    fprintf(stderr, "failed: %s\n", err);
    return 1;
  }
  write(1, out.data, out.width * out.height);
  OutImage_Clear(&out);
}
