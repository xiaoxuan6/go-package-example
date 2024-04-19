package main

import (
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"image"
	"log"
)

func main() {
	i, err := mergi.Import(impexp.NewFileImporter("1712295471.png"))
	if err != nil {
		log.Fatal(err)
	}

	watermarkImage, err := mergi.Import(impexp.NewFileImporter("watermark.png"))
	if err != nil {
		log.Fatal(err)
	}

	res, err := mergi.Watermark(watermarkImage, i, image.Pt(0, 0))
	if err != nil {
		log.Fatal(err)
	}

	_ = mergi.Export(impexp.NewFileExporter(res, "watermark_1.png"))
}
