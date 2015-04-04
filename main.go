package main

import (
	"flag"
	"fmt"
	"github.com/andviro/go-qrcode/qrcode"
)

var image *string

//go build -o test -ldflags "-linkmode external -extldflags -static"
func init() {
	image = flag.String("i", "", "image path")
}

func main() {

	flag.Parse()

	results, err := qrcode.ScanTIFF(*image)
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Printf("Symbol Type: %s, Data %s", result.SymbolType, result.Data)
	}
}
