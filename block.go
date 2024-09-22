package main

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
