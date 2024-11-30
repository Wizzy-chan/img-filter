package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"image/jpeg"
	"os"
	"io"
	"strings"
)

func UNUSED(x ...interface{}) {}

type ImageType int

const (
	UNKNOWN ImageType = iota
	PNG
	JPEG
)

func Usage(w io.Writer, prog string) {
	fmt.Fprintf(w, "Usage: %s path/to/image", prog)
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
	imageType := UNKNOWN
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR: Image path not provided\n")
		Usage(os.Stderr, program)
		os.Exit(1)
	}
	imgPath := os.Args[1]
	imgFile, err := os.Open(imgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open image %s: %s\n", imgPath, err.Error())
		os.Exit(1)
	}
	if strings.HasSuffix(imgPath, ".png") {
		imageType = PNG
	}
	if strings.HasSuffix(imgPath, ".jpg") || strings.HasSuffix(imgPath, ".jpeg") {
		imageType = JPEG
	}
	if imageType == UNKNOWN {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to work out image encoding from file extension.")
		os.Exit(1)
	}
	var img image.Image
	if imageType == PNG {
		img, err = png.Decode(imgFile)
	}
	if imageType == JPEG {
		img, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to decode image %s: %s\n", imgPath, err.Error())
		os.Exit(1)
	}
	grayscaleImg := Grayscale(img)
	outFile, err := os.Create(imgPath + "_grayscale")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open output file %s: %s\n", imgPath + "_grayscale", err.Error())
		os.Exit(1)
	}
	if imageType == PNG {
		err = png.Encode(outFile, grayscaleImg)
	}
	if imageType == JPEG {
		err = jpeg.Encode(outFile, grayscaleImg, nil)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to encode output file: %s\n", err.Error())
		os.Exit(1)
	}
}
