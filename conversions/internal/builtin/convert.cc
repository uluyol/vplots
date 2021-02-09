//========================================================================
//
// pdftopng.cc
//
// Copyright 2009 Glyph & Cog, LLC
//
//========================================================================

#include <cstdlib>
#include <string>

#include "GlobalParams.h"
#include "PDFDoc.h"
#include "Splash.h"
#include "SplashBitmap.h"
#include "SplashOutputDev.h"
#include "aconf.h"
#include "config.h"
#include "gtypes.h"
#include "parseargs.h"

extern "C" {
typedef struct {
  uint8_t *rgba_data;
  int32_t width;
  int32_t height;
} OutImage;

void OutImage_Clear(OutImage *im) {
  if (im->rgba_data != nullptr) {
    delete im->rgba_data;
  }
}
}

namespace {

void WriteToOutput(SplashBitmap *bitmap, OutImage *out) {
  out->width = bitmap->getWidth();
  out->height = bitmap->getHeight();
  out->rgba_data = new uint8_t[out->width * out->height * 4]();

  uint8_t *datap = out->rgba_data;

  Guchar *raw_datap = bitmap->getDataPtr();
  Guchar *raw_alphap = bitmap->getAlphaPtr();

  for (int j = 0; j < out->height; ++j) {
    for (int i = 0; i < out->width; ++i) {
      *datap++ = *raw_datap++;   // R
      *datap++ = *raw_datap++;   // G
      *datap++ = *raw_datap++;   // B
      *datap++ = *raw_alphap++;  // A
    }
  }
}

const char *RenderPDF(char *filename, int width, OutImage *out) {
  PDFDoc doc(filename);
  if (!doc.isOk()) {
    return "failed to open PDF document";
  }

  double resolution = 300;
  if (width != 0) {
    resolution = static_cast<double>(width) / (doc.getPageMediaWidth(1) / 72.0);
  }

  SplashColor paper_color;
  paper_color[0] = 0xff;
  paper_color[1] = 0xff;
  paper_color[2] = 0xff;

  SplashOutputDev out_dev(splashModeRGB8, 1, gFalse, paper_color);
  out_dev.setNoComposite(gTrue);
  out_dev.startDoc(doc.getXRef());
  doc.displayPage(&out_dev, 1, resolution, resolution, 0, gFalse, gFalse,
                  gFalse);

  WriteToOutput(out_dev.getBitmap(), out);
  return nullptr;
}

}  // namespace

extern "C" {
void XPDFConverter_GlobalInit() {
  globalParams = new GlobalParams("");
  globalParams->setupBaseFonts(nullptr);
  globalParams->setEnableFreeType((char *)"yes");
  globalParams->setAntialias((char *)"yes");
  globalParams->setVectorAntialias((char *)"yes");
  globalParams->setErrQuiet(gFalse);
}

const char *XPDFConverter_RenderPDF(char *filename, int width, OutImage *out) {
  return RenderPDF(filename, width, out);
}
}
