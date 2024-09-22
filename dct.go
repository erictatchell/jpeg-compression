package main

import "math"

func C(x int) float64 {
	if x == 0 {
		return 1 / math.Sqrt(2)
	} else {
		return 1
	}
}

// todo: test this
func DCT(F *[][]float64) [][]float64 {
	var result *Block = createEmptyBlock()
	for u := range 8 {
		for v := range 8 {
			var sum float64 = 0
			for x := range 8 {
				for y := range 8 {
					sum += math.Cos((float64(2*x+1)*float64(u)*math.Pi)/(2*8)) * math.Cos((float64(2*y+1)*float64(v)*math.Pi)/(2*8)) * (*F)[x][y]
				}
			}
			sum *= C(u) * C(v) * 0.25 // this is usually (2 / math.Sqrt(n * m)) but its always 8x8 blocks
			result.Matrix[u][v] = sum
		}
	}
	return result.Matrix
}

func IDCT(H [][]float64) [][]float64 {

	result := make([][]float64, 8)
	for i := range 8 {
		result[i] = make([]float64, 8)
	}

	for x := range 8 {
		for y := range 8 {
			var sum float64 = 0
			for u := range 8 {
				for v := range 8 {
					sum += 2 * ((C(u) * C(v)) / 8) * math.Cos((float64(2*x+1)*float64(u)*math.Pi)/(2*8)) * math.Cos((float64(2*y+1)*float64(v)*math.Pi)/float64(2*8)) * H[u][v]
				}
			}
			result[x][y] = sum
		}
	}
	return result

}
