package main

// jpeg.go
// author: Eric Tatchell
// for DCT blocks we are using 8x8

import (
	"fmt"
	"image"
)

func ConvertRGBToYCbCr(image image.Image, Y *[][]float64, Cb *[][]float64, Cr *[][]float64) {
	height := image.Bounds().Dy()
	width := image.Bounds().Dx()

	// go across the whole image
	for y := range height {
		for x := range width {
			color := image.At(x, y)
			r, g, b, _ := color.RGBA()

			// these formulas are from my multimedia textbook i cant remember which one oops
			(*Y)[x][y] = RGBtoYCbCr[0][0]*float64(r) +
				RGBtoYCbCr[0][1]*float64(g) +
				RGBtoYCbCr[0][2]*float64(b)

			(*Cb)[x][y] = RGBtoYCbCr[1][0]*float64(r) +
				RGBtoYCbCr[1][1]*float64(g) +
				RGBtoYCbCr[1][2]*float64(b) + 128

			(*Cr)[x][y] = RGBtoYCbCr[2][0]*float64(r) +
				RGBtoYCbCr[2][1]*float64(g) +
				RGBtoYCbCr[2][2]*float64(b) + 128
		}
	}
}

// the array that is generated is laid out as follows:
// [width, height, Y, Y, Y, Y, .... Cb, Cb, Cb, Cb, .... Cr, Cr, Cr, Cr]
func GetByteArray(image image.Image) []byte {
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

	var blocks []Block = GetBlocks(ycbcr)

	fmt.Printf("%d", len(blocks))

	return ycbcr
}

type Block struct {
	Matrix [][]float64
}

func GetBlocks(ycbcr []byte) []Block {
	// 4 for the width/height, / 64 for block size
	size := int((len(ycbcr))/64 + 4)
	var blocks = make([]Block, size)
	j := 0
	i := 4
	for i < len(ycbcr) {
		blocks[j] = GetBlock(&i, ycbcr, "y")
		j++
	}
	return blocks
}

func createEmptyBlock() *Block {
	block := Block{Matrix: make([][]float64, 8)}
	for j := range 8 {
		block.Matrix[j] = make([]float64, 8)
	}
	return &block
}

func GetBlock(i *int, ycbcr []byte, channel string) Block {
	// this weird bitshift stuff is left over logic from the original C# version
	// it was definitely ChatGPT
	// we are keeping it until i figure out another way to do this (never)
	width := int(ycbcr[0]<<8 | ycbcr[1])
	height := int(ycbcr[2]<<8 | ycbcr[3])
	size := width * height

	var block *Block = createEmptyBlock()

	for y := range 8 {
		for x := range 8 {
			if (*i) >= len(ycbcr) || ((*i) >= size && channel == "Y") {
				block.Matrix[x][y] = 0
			} else {
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
			(*ycbcr)[*i] = byte((*channel)[x][y])
			*i++
		}
	}
}

// 4:2:0
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
