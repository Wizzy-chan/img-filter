package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"image/jpeg"
	"os"
	"io"
	"flag"
)

func UNUSED(x ...interface{}) {}

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

func main() {
	outPath := flag.String("o", "", "The file`path` to write the resulting image to")
	outFormat := flag.String("f", "", "The `format` to use for encoding the resulting image")
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

	imgFile, err := os.Open(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open image %s: %s\n", inPath, err.Error())
		os.Exit(1)
	}

	_, format, err := image.DecodeConfig(imgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to decode image config: %s\n", err.Error())
		os.Exit(1)
	}

	imgFile.Seek(0, io.SeekStart)
	
	var img image.Image
	if format == "png" {
		img, err = png.Decode(imgFile)
	}
	if format == "jpeg" {
		img, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to decode image %s: %s\n", inPath, err.Error())
		os.Exit(1)
	}
	
	grayscaleImg := Grayscale(img)

	// TODO: Construct new output filepath if no output is provided with -o
	if *outPath == "" {
		fmt.Fprintf(os.Stderr, "ERROR: No output filepath provided.\n")
		os.Exit(1)
	}

	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open output file %s: %s\n", *outPath, err.Error())
		os.Exit(1)
	}

	if *outFormat != "" {
		format = *outFormat
	}
	if format == "png" {
		err = png.Encode(outFile, grayscaleImg)
	} else if format == "jpeg" {
		err = jpeg.Encode(outFile, grayscaleImg, nil)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: Unknown encoding format %s.\n", format)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to encode output file: %s\n", err.Error())
		os.Exit(1)
	}
}
