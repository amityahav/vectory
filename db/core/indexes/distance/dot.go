package distance

import (
	"golang.org/x/sys/cpu"
)

const DotProduct = "dot_product"

var dotProductImplementation = func(v1 []float32, v2 []float32) float32 {
	var sum float32

	for i := 0; i < len(v1); i += 8 {
		sum += v1[i] * v2[i]
		sum += v1[i+1] * v2[i+1]
		sum += v1[i+2] * v2[i+2]
		sum += v1[i+3] * v2[i+3]
		sum += v1[i+4] * v2[i+4]
		sum += v1[i+5] * v2[i+5]
		sum += v1[i+6] * v2[i+6]
		sum += v1[i+7] * v2[i+7]
	}

	//for i := 0; i < len(v1); i++ {
	//	sum += v1[i] * v2[i]
	//}

	return sum
}

func init() {
	if cpu.X86.HasAVX2 {
		//dotProductImplementation = asm.Dot
	}
}

func Dot(v1, v2 []float32) float32 {
	return dotProductImplementation(v1, v2)
}
