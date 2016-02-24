#include <stdlib.h>
#include <tiffio.h>
#include "tiffread.h"

typedef union {
    struct {
        unsigned char r,g,b,a;
    } rgba;
    uint32 t;
} pixel;
 
int tiffread (const char *name, int *width, int *height, void **raw)
{
 
  TIFF* tif = TIFFOpen(name, "r");
  if (tif) {
    size_t npixels;
    pixel *tmp;
    int i;
 
    TIFFGetField(tif, TIFFTAG_IMAGEWIDTH, width);
    TIFFGetField(tif, TIFFTAG_IMAGELENGTH, height);
    npixels = *width * *height;
 
    tmp = (pixel *)malloc(npixels * sizeof (pixel));
    if (tmp != NULL) {
      if (!TIFFReadRGBAImage(tif, *width, *height, (uint32*) tmp, 0)) {
            free(tmp);
            TIFFClose(tif);
            return 2;
        } else {
            *raw = malloc(npixels);
            if (!*raw) {
                free(tmp);
                TIFFClose(tif);
                return 3;
            }
            for (i=0; i < npixels; i++) {
                /*((unsigned char*) *raw)[i] = (unsigned char) (tmp[i].rgba.r * 0.2989 + tmp[i].rgba.g * 0.5870+ tmp[i].rgba.b * 0.1140);*/
                ((unsigned char*) *raw)[i] = tmp[i].rgba.g; // XXX: faster
            }
        }
        free(tmp);
      }
      TIFFClose(tif);
    } else {
        return 5;
    }
  return 0;
}

int test_dir(const char *name)
{
    TIFF* tif = TIFFOpen(name, "r");
    int dircount = -1, width, height;
    if (tif) {
        dircount = 0;
        do {
            TIFFGetField(tif, TIFFTAG_IMAGEWIDTH, &width);
            TIFFGetField(tif, TIFFTAG_IMAGELENGTH, &height);
            dircount++;
            printf("%i x %i\n", width, height);
        } while (TIFFReadDirectory(tif));
        TIFFClose(tif);
    }
    return dircount;
}
