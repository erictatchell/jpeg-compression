// block.go
// 8x8 blocks for dct and quantize

package main

import "fmt"

type Block struct {
	channel string
	Matrix  [][]float64
}

func createEmptyBlock(channel string) *Block {
	block := Block{Matrix: make([][]float64, 8)}
	for j := range 8 {
		block.Matrix[j] = make([]float64, 8)
	}
	block.channel = channel
	return &block
}

func GetBlocks(channel string, ycbcr []byte) []Block {
	// 4 for the width/height, / 64 for block size
	arraySize := int((len(ycbcr) - 4) / 64)
	fmt.Println("Array size: ", arraySize)
	fmt.Println("YCBCR length: ", len(ycbcr))

	// this weird bitshift stuff is left over logic from the original C# version it was definitely ChatGPT
	// we are keeping it until i figure out another way to do this (never)
	width := int(ycbcr[0]<<8 | ycbcr[1])
	height := int(ycbcr[2]<<8 | ycbcr[3])
	idx := 4
	imageSize := width * height

	if channel == "CbCr" {
		// if CbCr, move the array pointer past the width/height and Y's
		idx = idx + imageSize

		// account for subsampling
		imageSize = imageSize / 4
		arraySize = arraySize / 4
	}

	var blocks = make([]Block, arraySize)

	for j := range arraySize {
		blocks[j] = GetBlock(&idx, ycbcr, channel, imageSize)
	}

	fmt.Println("Finished creating blocks")
	return blocks
}

func GetBlock(i *int, ycbcr []byte, channel string, size int) Block {
	var block *Block = createEmptyBlock(channel)
	idx := *i

	for r := range 8 {
		for c := range 8 {
			if idx >= len(ycbcr) || (idx >= size && channel == "Y") {
				block.Matrix[r][c] = 0
			} else {
				// 99% sure this conversion from byte to float64 doesn't work but we'll find out later.!!
				block.Matrix[r][c] = float64(ycbcr[(*i)])
				(*i)++
			}
		}
	}

	return *block
}
