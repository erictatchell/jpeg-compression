//
// Custom JPEG compressor
// Learning golang - Eric Tatchell
//

package main

import (
	"fmt"
	"image"
	"os"
)

var RGBtoYCbCr [][]float64 = [][]float64{
	{0.299, 0.587, 0.114},
	{-0.168736, -0.331264, 0.5},
	{0.5, -0.418688, -0.081312},
}

var YCbCrtoRGB [][]float64 = [][]float64{
	{1.0, 0.0, 1.402},
	{1.0, -0.34414, -0.71414},
	{1.0, 1.772, 0.0},
}

var Luminance [][]uint8 = [][]uint8{
	{16, 11, 10, 16, 24, 40, 51, 61},
	{12, 12, 14, 19, 26, 58, 60, 55},
	{14, 13, 16, 24, 40, 57, 69, 56},
	{14, 17, 22, 29, 51, 87, 80, 62},
	{18, 22, 37, 56, 68, 109, 103, 77},
	{24, 35, 55, 64, 81, 104, 113, 92},
	{49, 64, 78, 87, 103, 121, 120, 101},
	{72, 92, 95, 98, 112, 100, 103, 99},
}

var Chrominance [][]uint8 = [][]uint8{
	{17, 18, 24, 47, 99, 99, 99, 99},
	{18, 21, 26, 66, 99, 99, 99, 99},
	{24, 26, 56, 99, 99, 99, 99, 99},
	{47, 66, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
	{99, 99, 99, 99, 99, 99, 99, 99},
}

type JPEGImage struct {
	origImage image.Image
	width     int
	height    int
	ycbcr     []byte
	rgb       []byte
}

// with the help of https://stackoverflow.com/questions/49594259/reading-image-in-go
func openImageFromPath(imagePath string) (image.Image, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Unable to open image from path: " + imagePath)
		return nil, err
	}
	defer imageFile.Close()
	image, _, err := image.Decode(imageFile)
	return image, err

}

func convertRGBtoYCbCr(image image.Image) []byte {
	var width int = image.Bounds().Dx()
	var height int = image.Bounds().Dy()
	var numPixels int = width * height
	var ycbcr = make([]byte, numPixels)

	// for some reason this is how to finish initializing 2D arrays? there must be a better way!
	Y := make([][]float64, height)
	Cb := make([][]float64, height)
	Cr := make([][]float64, height)
	for i := range height {
		Y[i] = make([]float64, width)
		Cb[i] = make([]float64, width)
		Cr[i] = make([]float64, width)
	}

	for y := range height {
		for x := range width {
			color := image.At(x, y)
			r, g, b, _ := color.RGBA()
			Y[x][y] = RGBtoYCbCr[0][0]*float64(r) +
				RGBtoYCbCr[0][1]*float64(g) +
				RGBtoYCbCr[0][2]*float64(b)

			Cb[x][y] = RGBtoYCbCr[1][0]*float64(r) +
				RGBtoYCbCr[1][1]*float64(g) +
				RGBtoYCbCr[1][2]*float64(b) + 128

			Cr[x][y] = RGBtoYCbCr[2][0]*float64(r) +
				RGBtoYCbCr[2][1]*float64(g) +
				RGBtoYCbCr[2][2]*float64(b) + 128
		}
	}

	return ycbcr
}

func main() {
	var jpegImage JPEGImage = JPEGImage{}
	var imagePath string = os.Args[1]
	image, err := openImageFromPath(imagePath)
	if err != nil {
		fmt.Println("Error opening/decoding the image file.")
		os.Exit(1)
	}
	jpegImage.origImage = image
	jpegImage.ycbcr = convertRGBtoYCbCr(image)

	fmt.Println(imagePath)
}
