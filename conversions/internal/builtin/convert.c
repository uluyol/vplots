// Forked from pdf2svg project (https://github.com/dawbarton/pdf2svg)

// Copyright (C) 2007-2013 David Barton (davebarton@cityinthesky.co.uk)
// <http://www.cityinthesky.co.uk/>

// Copyright (C) 2007 Matthew Flaschen (matthew.flaschen@gatech.edu)
// Updated to allow conversion of all pages at once.

//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 2 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Definitely works with Cairo v1.2.6 and Poppler 0.5.4

#include <assert.h>
#include <cairo-svg.h>
#include <cairo.h>
#include <glib.h>
#include <poppler-document.h>
#include <poppler-page.h>
#include <poppler.h>
#include <stdlib.h>
#include <string.h>

// Begin theft from ePDFview (Copyright (C) 2006 Emma's Software) under the GPL
static gchar *GetAbsoluteFileName(const gchar *filename) {
  gchar *abs_filename = NULL;
  if (g_path_is_absolute(filename)) {
    abs_filename = g_strdup(filename);
  } else {
    gchar *cur_dir = g_get_current_dir();
    abs_filename = g_build_filename(cur_dir, filename, NULL);
    g_free(cur_dir);
  }
  return abs_filename;
}
// End theft from ePDFview

struct OutBuf {
  char *data;
  int len;
  int cap;
};

static void InitOutBuf(struct OutBuf *buf) {
  buf->data = NULL;
  buf->len = 0;
  buf->cap = 0;
}

static cairo_status_t WriteToOutBuf(void *closure, const unsigned char *data,
                                    unsigned int length) {
  struct OutBuf *buf = (struct OutBuf *)closure;
  if ((buf->len + length) > buf->cap) {
    int new_size = length;
    if (buf->cap * 2 > new_size) {
      new_size = buf->cap * 2;
    }
    // Grow buf->data
    buf->data = realloc(buf->data, new_size);
  }
  memmove(buf->data + buf->len, data, length);
  buf->len += length;
  return CAIRO_STATUS_SUCCESS;
}

#define WANT_DPI 300

static int HighDPIPix(double width) { return (width * WANT_DPI) / 72; }

static char *ConvertPageToSVG(PopplerPage *page, struct OutBuf *buf,
                              double width, double height) {
  cairo_surface_t *surface = cairo_svg_surface_create_for_stream(
      WriteToOutBuf, (void *)buf, width, height);
  cairo_t *drawcontext = cairo_create(surface);

  poppler_page_render_for_printing(page, drawcontext);
  cairo_show_page(drawcontext);

  cairo_destroy(drawcontext);
  cairo_surface_destroy(surface);

  return NULL;
}

static char *ConvertPageToPNG(PopplerPage *page, struct OutBuf *buf,
                              double width, double height) {
  cairo_surface_t *surface = cairo_image_surface_create(
      CAIRO_FORMAT_RGB24, HighDPIPix(width), HighDPIPix(height));
  cairo_t *drawcontext = cairo_create(surface);

  poppler_page_render_for_printing(page, drawcontext);
  cairo_show_page(drawcontext);

  cairo_surface_write_to_png_stream(surface, WriteToOutBuf, buf);

  cairo_destroy(drawcontext);
  cairo_surface_destroy(surface);

  return NULL;
}

char *ConvertPDF(char *filename, struct OutBuf *output, int format) {
  gchar *abs_filename = GetAbsoluteFileName(filename);
  gchar *filename_uri = g_filename_to_uri(abs_filename, NULL, NULL);
  g_free(abs_filename);

  // Open the PDF file
  PopplerDocument *pdffile =
      poppler_document_new_from_file(filename_uri, NULL, NULL);
  g_free(filename_uri);
  if (pdffile == NULL) {
    return "failed to open file";
  }

  PopplerPage *page = poppler_document_get_page(pdffile, 0);
  char *conversion_error = NULL;
  if (page == NULL) {
    conversion_error = "page does not exist";
    goto CleanupPDF;
  }
  InitOutBuf(output);
  double width, height;
  poppler_page_get_size(page, &width, &height);
  if (format == 0) {
    conversion_error = ConvertPageToSVG(page, output, width, height);
  } else if (format == 1) {
    conversion_error = ConvertPageToPNG(page, output, width, height);
  } else {
    conversion_error = "invalid format";
  }

  g_object_unref(page);
CleanupPDF:
  g_object_unref(pdffile);
  return conversion_error;
}
