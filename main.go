package main

import (
	"fmt"
	"github.com/zyxar/image2ascii/ascii"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	//f, err := os.Open("88888888.jpg")
	//defer f.Close()
	//if err != nil {
	//	_, _ = fmt.Fprintln(os.Stderr, err)
	//	return
	//}
	//
	//img, _, err := image.Decode(f)
	//if err != nil {
	//	_, _ = fmt.Fprintln(os.Stderr, err)
	//	return
	//}
	//bounds := img.Bounds()
	//width := bounds.Max.X / 4
	//height := bounds.Max.Y / 4
	//
	//fmt.Println("图片宽", width)
	//fmt.Println("图片高", height)

	opt := ascii.Options{
		Width:  50,
		Height: 50,
		Color:  true,
		Invert: false,
		Flipx:  false,
		Flipy:  false,
	}

	f1, err := os.Open("88888888.jpg")
	defer f1.Close()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	a, err := ascii.Decode(f1, opt)
	if err != nil {
		fmt.Println("111")
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	_, err = a.WriteTo(os.Stdout)
	if err != nil {
		fmt.Println("222")
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}
