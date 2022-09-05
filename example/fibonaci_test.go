package example

import (
	"testing"
)

func BenchmarkFibonacciRecursion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibonacciRecursion(10)
	}
}

func BenchmarkFibonacciLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibonacciLoop(10)
	}
}

func FibonacciRecursion(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursion(n-1) + FibonacciRecursion(n-2)
}

func FibonacciLoop(n int) int {
	f := make([]int, n+1, n+2)
	if n < 2 {
		f = f[0:2]
	}
	f[0] = 0
	f[1] = 1
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f[n]
}
