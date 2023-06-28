package distance

import "math"

const Euclidean = "euclidean_distance"

func EuclideanDistance(v1, v2 []float32) float32 {
	var sum float64

	for i := 0; i < len(v1); i++ {
		sum += math.Pow(float64(v1[i]-v2[i]), float64(2))
	}

	return float32(math.Sqrt(sum))
}

func manhatthanDistance(v1, v2 []float32) float32 {
	var sum float64

	for i := 0; i < len(v1); i++ {
		sum += math.Abs(float64(v1[i] - v2[i]))
	}

	return float32(sum)
}

func cosineSimilarity(v1, v2 []float32) float32 {
	dot := Dot(v1, v2)
	m1, m2 := magnitude(v1), magnitude(v2)

	return dot / (m1 * m2)
}

func magnitude(v []float32) float32 {
	var sum float64

	for i := 0; i < len(v); i++ {
		sum += math.Pow(float64(v[i]), 2)
	}

	return float32(math.Sqrt(sum))
}
