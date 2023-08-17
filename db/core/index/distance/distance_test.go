package distance

import (
	"testing"
)

func BenchmarkEuclideanDistance(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	dims := []int{128, 256, 512, 1024, 2048}

	b.Run("fractions", func(b *testing.B) {
		for _, dim := range dims {
			for i := 0; i < b.N; i++ {
				v1, v2 := randomVector(dim, 0.5), randomVector(dim, -0.233)
				EuclideanDistance(v1, v2)
			}
		}
	})

	b.Run("whole", func(b *testing.B) {
		for _, dim := range dims {
			for i := 0; i < b.N; i++ {
				v1, v2 := randomVector(dim, 23), randomVector(dim, -14)
				EuclideanDistance(v1, v2)
			}
		}
	})
}

func BenchmarkDot(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	dims := []int{128, 256, 512, 1024, 2048}

	b.Run("fractions", func(b *testing.B) {
		for _, dim := range dims {
			for i := 0; i < b.N; i++ {
				v1, v2 := randomVector(dim, 0.5), randomVector(dim, -0.233)
				Dot(v1, v2)
			}
		}
	})

	b.Run("whole", func(b *testing.B) {
		for _, dim := range dims {
			for i := 0; i < b.N; i++ {
				v1, v2 := randomVector(dim, 23), randomVector(dim, -14)
				Dot(v1, v2)
			}
		}
	})
}

func randomVector(dim int, n float32) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = n
	}

	return vec
}
