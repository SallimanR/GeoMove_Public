package array

import (
	"testing"
)

const lenght = 10000

func BenchmarkArrayIteration(b *testing.B) {
	arr := make([]*uint32, lenght)
	for i := range uint32(lenght) {
		arr[i] = &i
	}

	result := make([]*uint32, lenght)
	for b.Loop() {
		for i := range arr {
			result[i] = arr[i]
		}
	}
	// log.Println(result)
}

func BenchmarkArrayIterationUnrolled(b *testing.B) {
	arr := make([]*uint32, lenght)
	for i := range uint32(lenght) {
		arr[i] = &i
	}

	result := make([]*uint32, lenght)
	for b.Loop() {
		for i := 0; i < lenght; i += 5 {
			result[i] = arr[i]
			result[i+1] = arr[i+1]
			result[i+2] = arr[i+2]
			result[i+3] = arr[i+3]
			result[i+4] = arr[i+4]
		}
	}
	// log.Println(result)
}

func BenchmarkArrayIterationWithIF(b *testing.B) {
	arr := make([]*uint32, lenght)
	for i := range uint32(lenght) {
		arr[i] = &i
	}

	result := make([]*uint32, lenght)
	for b.Loop() {
		for i := range arr {
			if arr[i] == nil {
				continue
			}
			result[i] = arr[i]
		}
	}
	// log.Println(result)
}
