package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"io"
)

func UNUSED(x ...interface{}) {}

func Usage(w io.Writer, prog string) {
	fmt.Fprintf(w, "Usage: %s path/to/image", prog)
}

func PrintErr(err string, args ...string) {
	fmt.Fprintf(os.Stderr, err, args)
}

func GrayscaleColor(clr color.Color) color.Color {
	r,g,b,a := clr.RGBA()
	y := (r + g + b) / 3 >> 8
	return color.NRGBA{ uint8(y), uint8(y), uint8(y), uint8(a >> 8) }
}

func Grayscale(img image.Image) image.Image {
	out := image.NewNRGBA(img.Bounds())
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			out.Set(x, y, GrayscaleColor(img.At(x, y)))
		}
	}
	return out
}

func main() {
	program := os.Args[0]
	if len(os.Args) < 2 {
		PrintErr("ERROR: Image path not provided")
		Usage(os.Stderr, program)
		os.Exit(1)
	}
	imgPath := os.Args[1]
	imgFile, err := os.Open(imgPath)
	if err != nil {
		PrintErr("ERROR: Unable to open image %s: %s\n", imgPath, err.Error())
		os.Exit(1)
	}
	img, err := png.Decode(imgFile)
	if err != nil {
		PrintErr("ERROR: Unable to decode image %s: %s\n", imgPath, err.Error())
		os.Exit(1)
	}
	grayscaleImg := Grayscale(img)
	outFile, err := os.Create(imgPath + "_grayscale")
	if err != nil {
		PrintErr("ERROR: Unable to open output file %s: %s\n", imgPath + "_grayscale", err.Error())
		os.Exit(1)
	}
	err = png.Encode(outFile, grayscaleImg)
	if err != nil {
		PrintErr("ERROR: Unable to encode output png file: %s\n", err.Error())
		os.Exit(1)
	}
}
