package step_2

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

func MyFunction() []byte {

	data := make([]byte, 1024)

	rand.Read(data)

	sum := sha256.Sum256(data)
	return sum[:]
}

func BenchmarkMyFunction(b *testing.B) {

	for i := 0; i < b.N; i++ {
		MyFunction()
	}

}

/*
Run the Benchmark
- go test -bench

Profile the CPU usage
- go test -bench . -cpuprofile cpu.prof
- go tool pprof cpu.prof

Profile the memory usage
- go test -bench . -memprofile mem.prof
- go tool pprof mem.prof

*/
