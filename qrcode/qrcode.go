package qrcode

// #cgo LDFLAGS: -lzbar -ltiff
// #include <zbar.h>
// #include "tiffread.h"
// typedef zbar_image_cleanup_handler_t *zbar_image_set_data_callback;
import "C"

import (
	"errors"
	"fmt"
	"image"
	"unsafe"
)

type Result struct {
	SymbolType string
	Data       string
}

func ScanTIFF(imgPath string) (results []Result, err error) {
	pth := C.CString(imgPath)
	scanner := C.zbar_image_scanner_create()

	C.zbar_image_scanner_set_config(scanner, 0, C.ZBAR_CFG_ENABLE, 1)
	defer C.zbar_image_scanner_destroy(scanner)

	var width, height C.int = 0, 0
	var raw unsafe.Pointer = nil
	errorCode := C.tiffread(pth, &width, &height, &raw)
	if int(errorCode) != 0 {
		err = errors.New(fmt.Sprintf("Error reading from tiff file. Error code %d", errorCode))
		return
	}

	image := C.zbar_image_create()

	defer C.zbar_image_destroy(image)

	C.zbar_image_set_format(image, C.ulong(808466521))
	C.zbar_image_set_size(image, C.uint(width), C.uint(height))

	f := C.zbar_image_set_data_callback(C.zbar_image_free_data)
	C.zbar_image_set_data(image, raw, C.ulong(width*height), f)

	C.zbar_scan_image(scanner, image)

	symbol := C.zbar_image_first_symbol(image)

	for ; symbol != nil; symbol = C.zbar_symbol_next(symbol) {
		typ := C.zbar_symbol_get_type(symbol)
		data := C.zbar_symbol_get_data(symbol)
		symbolType := C.GoString(C.zbar_get_symbol_name(typ))
		dataString := C.GoString(data)
		results = append(results, Result{symbolType, dataString})
	}

	return
}

func ScanImage(img image.Image) (results []Result, err error) {

	scanner := C.zbar_image_scanner_create()
	C.zbar_image_scanner_set_config(scanner, 0, C.ZBAR_CFG_ENABLE, 1)

	defer C.zbar_image_scanner_destroy(scanner)

	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	gray := image.NewGray(rect)

	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	image := C.zbar_image_create()

	defer C.zbar_image_destroy(image)

	C.zbar_image_set_format(image, C.ulong(808466521))
	C.zbar_image_set_size(image, C.uint(width), C.uint(height))

	C.zbar_image_set_data(image, unsafe.Pointer(&gray.Pix[0]), C.ulong(width*height), nil)

	C.zbar_scan_image(scanner, image)

	symbol := C.zbar_image_first_symbol(image)

	for ; symbol != nil; symbol = C.zbar_symbol_next(symbol) {
		typ := C.zbar_symbol_get_type(symbol)
		data := C.zbar_symbol_get_data(symbol)
		symbolType := C.GoString(C.zbar_get_symbol_name(typ))
		dataString := C.GoString(data)
		results = append(results, Result{symbolType, dataString})
	}

	return
}
