package datastructures

import (
	"testing"
)

const lenght uint32 = 100000

type Point[T any] struct {
	x T
	y T
}

func BenchmarkPointers(b *testing.B) {
	arr := make([]*uint64, 0, lenght)
	for i := range uint64(lenght) {
		arr = append(arr, &i)
	}
	result := make([]uint64, lenght)

	for b.Loop() {
		for i := range arr {
			result = append(result, *arr[i])
		}
	}
}

func BenchmarkIndexes(b *testing.B) {
	arr := make([]Point[uint64], 0, lenght)
	for i := range uint64(lenght) {
		arr = append(arr, Point[uint64]{x: i, y: i})
	}

	result := make([]uint64, lenght)
	for b.Loop() {
		for i := range arr {
			result = append(result, arr[i].x)
		}
	}
}

func BenchmarkIndexesHashMap(b *testing.B) {
	// hashMap := make(map[uint32]Data, lenght)
	hashMap := make(map[uint32]Point[uint64])
	for i := range lenght {
		hashMap[i] = Point[uint64]{x: uint64(i), y: uint64(i)}
	}

	result := make([]uint64, lenght)
	for b.Loop() {
		for i := range lenght {
			result = append(result, hashMap[i].x)
		}
	}
}

// type Data[T any] struct {
// 	sync.RWMutex
// 	points []T
// }

// func BenchmarkIndexesWithRWLock(b *testing.B) {
// 	arr := make([]Point, 0, lenght)
// 	for i := range lenght {
// 		arr = append(arr, Point{x: i, y: i})
// 	}
// 	data := Data[Point]{points: arr}
//
// 	result := make([]uint32, lenght)
// 	for b.Loop() {
// 		for i := range arr {
// 			data.RLock()
// 			result = append(result, data.points[i].x)
// 			data.RUnlock()
// 		}
// 	}
// }
//
// func BenchmarkSyncMap(b *testing.B) {
// 	var arr datastructures.SyncMap[uint32, Point]
// 	for i := range lenght {
// 		arr.Store(i, Point{x: i, y: i})
// 	}
//
// 	result := make([]uint32, lenght)
// 	for b.Loop() {
// 		for i := range lenght {
// 			data, ok := arr.Load(i)
// 			if ok {
// 				result[i] = data.x
// 			}
// 		}
// 	}
// }
