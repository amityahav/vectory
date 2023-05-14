package distance

const DotProduct = "dot_product"

var dotProductImplementation = func(v1 []float32, v2 []float32) float32 {
	var sum float32

	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}

	return sum
}

func init() {
	//if cpu.X86.HasAVX2 {
	//	dotProductImplementation = asm.Dot
	//}
}

func Dot(v1, v2 []float32) float32 {
	return dotProductImplementation(v1, v2)
}
