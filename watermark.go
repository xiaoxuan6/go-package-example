package main

import (
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"image"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "/../../hexo/mkdocs/docs/static/images/")

	_ = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() == true {
			return nil
		}

		watermark(path)
		return nil
	})
}

func watermark(img string) {
	i, err := mergi.Import(impexp.NewFileImporter(img))
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

	_ = mergi.Export(impexp.NewFileExporter(res, img))
}
