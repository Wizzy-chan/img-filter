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
	"flag"
)

func UNUSED(x ...interface{}) {}

type ImageFormat int

const (
	UNKNOWN ImageFormat = iota
	PNG
	JPEG
)

func Usage(w io.Writer) {
	fmt.Fprintf(w, "Usage: img-filter [flags] path/to/image\n")
	flag.PrintDefaults()
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

type ErrorUnknownFormat string

func (err ErrorUnknownFormat) Error() string {
	return string(err)
}

func GetImageFormat(format, inPath string) (ImageFormat, error) {
	if format == "png" {
		return PNG, nil
	}
	if format == "jpeg" || format == "jpg" {
		return JPEG, nil
	}

	if format != "" {
		return UNKNOWN, ErrorUnknownFormat("unknown format " + format)
	}

	if strings.HasSuffix(inPath, ".png") {
		return PNG, nil
	}
	if strings.HasSuffix(inPath, ".jpg") || strings.HasSuffix(inPath, ".jpeg") {
		return JPEG, nil
	}

	return UNKNOWN, ErrorUnknownFormat("could not work out format from file extension")
}

func main() {
	outPath := flag.String("o", "", "The file`path` to write the resulting image to")
	format := flag.String("f", "", "The `format` to use for encoding and decoding")
	help := flag.Bool("help", false, "Show this help message")
	flag.BoolVar(help, "h", false, "")
	
	flag.Parse()

	if *help {
		Usage(os.Stdout)
		os.Exit(0)
	}
	
	inPath := flag.Arg(0)
	if inPath == "" {
		fmt.Fprintf(os.Stderr, "ERROR: No image path provided\n")
		os.Exit(1)
	}
		
	imageFormat, err := GetImageFormat(*format, inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	
	imgFile, err := os.Open(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open image %s: %s\n", inPath, err.Error())
		os.Exit(1)
	}
	
	var img image.Image
	if imageFormat == PNG {
		img, err = png.Decode(imgFile)
	}
	if imageFormat == JPEG {
		img, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to decode image %s: %s\n", inPath, err.Error())
		os.Exit(1)
	}
	
	grayscaleImg := Grayscale(img)

	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open output file %s: %s\n", *outPath, err.Error())
		os.Exit(1)
	}

	if imageFormat == PNG {
		err = png.Encode(outFile, grayscaleImg)
	}
	if imageFormat == JPEG {
		err = jpeg.Encode(outFile, grayscaleImg, nil)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to encode output file: %s\n", err.Error())
		os.Exit(1)
	}
}
