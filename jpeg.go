package main

// jpeg.go
// author: Eric Tatchell
// for DCT blocks we are using 8x8

import (
	"fmt"
	"image"
	"math"
)

func ConvertRGBToYCbCr(image image.Image, Y *[][]float64, Cb *[][]float64, Cr *[][]float64) {
	height := image.Bounds().Dy()
	width := image.Bounds().Dx()

	// go across the whole image
	fmt.Println("Starting YCbCr conversion...")
	for y := range height {
		for x := range width {
			color := image.At(x, y)
			r, g, b, _ := color.RGBA()

			// these formulas are from my multimedia textbook i cant remember which one oops
			(*Y)[y][x] = RGBtoYCbCr[0][0]*float64(r) +
				RGBtoYCbCr[0][1]*float64(g) +
				RGBtoYCbCr[0][2]*float64(b)

			(*Cb)[y][x] = RGBtoYCbCr[1][0]*float64(r) +
				RGBtoYCbCr[1][1]*float64(g) +
				RGBtoYCbCr[1][2]*float64(b) + 128

			(*Cr)[y][x] = RGBtoYCbCr[2][0]*float64(r) +
				RGBtoYCbCr[2][1]*float64(g) +
				RGBtoYCbCr[2][2]*float64(b) + 128
		}
	}
	fmt.Println("Finished YCbCr conversion")
}

// the array that is generated is laid out as follows:
// [width, height, Y, Y, Y, Y, .... Cb, Cb, Cb, Cb, .... Cr, Cr, Cr, Cr]
func GetByteArray(image image.Image) ([]byte, error) {
	var width int = image.Bounds().Dx()
	var height int = image.Bounds().Dy()

	fmt.Printf("Dimensions: %d x %d\n", width, height)

	numPixels := width * height

	// without subsampling this is 3, for 3 channels (Y, Cb, Cr)
	// we are using 4:2:0 chroma subsampling, which reduces Cb and Cr by a factor of 1/4 each
	// therefore, 1 + 0.25 + 0.25
	subsampleOffset := 1.5

	// + 4 for storing width and height (least & most significant byte)
	var ycbcr = make([]byte, int(float64(numPixels)*subsampleOffset+2+4))
	fmt.Printf("Created byte array of size %d\n", len(ycbcr))

	// for some reason this is how to initialize 2D arrays? gross!
	Y := make([][]float64, height)
	Cb := make([][]float64, height)
	Cr := make([][]float64, height)
	for i := range Y {
		Y[i] = make([]float64, width)
		Cb[i] = make([]float64, width)
		Cr[i] = make([]float64, width)
	}

	// filling the matrices
	ConvertRGBToYCbCr(image, &Y, &Cb, &Cr)

	// big compression tings
	ChromaSubsample(width, height, &Cb)
	ChromaSubsample(width, height, &Cr)

	// storing the most and least significant byte for width and height to prevent overflow
	ycbcr[0] = byte(width >> 8)
	ycbcr[1] = byte(width & 0xFF)
	ycbcr[2] = byte(height >> 8)
	ycbcr[3] = byte(height & 0xFF)

	i := 4

	// [..., Y, Y, Y, ..., Cb, Cb, Cb, ..., Cr, Cr, Cr, ...]
	PutChannelInByteArray(&i, width, height, &Y, &ycbcr)
	PutChannelInByteArray(&i, width/2, height/2, &Cb, &ycbcr)
	PutChannelInByteArray(&i, width/2, height/2, &Cr, &ycbcr)
	fmt.Println("Filled the byte array")

	var blocks []Block = GetBlocks(ycbcr)

	// this is just so beautiful i love Go
	fmt.Println("Starting DCT and Quantization")
	for _, block := range blocks {
		if len(block.Matrix) != 0 {
			DCT(&block.Matrix)
			Quantize(block.channel, &block.Matrix)
		}
	}

	fmt.Println("Finished DCT and Quantization")
	return ycbcr, nil
}

func Quantize(channel string, block *[][]float64) {
	for y := range 8 {
		for x := range 8 {
			if channel == "Y" {
				(*block)[x][y] = math.Round((*block)[x][y] / Luminance[x][y])
			} else {
				(*block)[x][y] = math.Round((*block)[x][y] / Chrominance[x][y])
			}
		}
	}
}

func GetBlocks(ycbcr []byte) []Block {
	// 4 for the width/height, / 64 for block size
	size := int((len(ycbcr))/64 + 4)
	fmt.Printf("Creating %d 8x8 blocks from the byte array\n", size)

	var blocks = make([]Block, size)
	i := 4
	for j := range blocks {
		blocks[j] = GetBlock(&i, ycbcr, "Y")
	}

	fmt.Println("Finished creating blocks")
	return blocks
}

func GetBlock(i *int, ycbcr []byte, channel string) Block {
	// this weird bitshift stuff is left over logic from the original C# version
	// it was definitely ChatGPT
	// we are keeping it until i figure out another way to do this (never)
	width := int(ycbcr[0]<<8 | ycbcr[1])
	height := int(ycbcr[2]<<8 | ycbcr[3])
	size := width * height

	var block *Block = createEmptyBlock(channel)

	for y := range 8 {
		for x := range 8 {
			if (*i) >= len(ycbcr) || ((*i) >= size && channel == "Y") {
				block.Matrix[x][y] = 0
			} else { // cbcr
				// 99% sure this conversion from byte to float64 doesn't work but we'll find out later.!!
				block.Matrix[x][y] = float64(ycbcr[(*i)])
				(*i)++
			}
		}
	}

	return *block
}

func PutChannelInByteArray(i *int, width int, height int, channel *[][]float64, ycbcr *[]byte) {
	for y := range height {
		for x := range width {
			// kill me
			(*ycbcr)[*i] = byte((*channel)[y][x])
			*i++
		}
	}
}

// 4:2:0
func ChromaSubsample(width int, height int, channel *[][]float64) {
	fmt.Println("Starting subsample (4:2:0)...")
	temp := make([][]float64, height/2)
	for i := range temp {
		temp[i] = make([]float64, width)
	}
	for y := range width / 2 {
		for x := range height / 2 {
			temp[x][y] = (*channel)[x*2][y*2]
		}
	}
	fmt.Println("Finished subsample")
	*channel = temp
}
