package main

import (
	"image"
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

func convertRGBtoYCbCr(image image.Image) []byte {
	var width int = image.Bounds().Dx()
	var height int = image.Bounds().Dy()

	numPixels := width * height

	// without subsampling this is 3, for 3 channels (Y, Cb, Cr)
	// we are using 4:2:0 chroma subsampling, which reduces Cb and Cr by a factor of 1/4 each
	// therefore, 1 + 0.25 + 0.25
	subsampleOffset := 1.5

	// + 4 for storing width and height (least & most significant byte)
	var ycbcr = make([]byte, int(float64(numPixels)*subsampleOffset+2+4))

	// for some reason this is how to initialize 2D arrays? gross!
	Y := make([][]float64, height)
	Cb := make([][]float64, height)
	Cr := make([][]float64, height)
	for i := range Y {
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

	ChromaSubsample(width, height, &Cb)
	ChromaSubsample(width, height, &Cr)

	// storing the most and least significant byte for width and height
	ycbcr[0] = byte(width >> 8)
	ycbcr[1] = byte(width & 0xFF)
	ycbcr[2] = byte(height >> 8)
	ycbcr[3] = byte(height & 0xFF)

	i := 4

	PutChannelInByteArray(&i, width, height, &Y, &ycbcr)
	PutChannelInByteArray(&i, width/2, height/2, &Cb, &ycbcr)
	PutChannelInByteArray(&i, width/2, height/2, &Cr, &ycbcr)
	return ycbcr
}

func PutChannelInByteArray(i *int, width int, height int, channel *[][]float64, ycbcr *[]byte) {
	for y := range height {
		for x := range width {
			(*ycbcr)[*i] = byte((*channel)[x][y])
			*i++
		}
	}

}

func ChromaSubsample(width int, height int, channel *[][]float64) {
	temp := make([][]float64, height/2)
	for i := range temp {
		temp[i] = make([]float64, width)
	}
	for y := range width / 2 {
		for x := range height / 2 {
			temp[x][y] = (*channel)[x*2][y*2]
		}
	}
	*channel = temp
}
