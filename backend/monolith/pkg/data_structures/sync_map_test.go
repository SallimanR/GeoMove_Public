package datastructures

import (
	"log"
	"sync"
	"testing"
)

var mapSize uint32 = 100

func BenchmarkStdMap(b *testing.B) {
	gm := make(map[uint32]uint32)

	var i uint32
	for i = range mapSize {
		gm[i] = i
	}

	mapValues := make([]uint32, mapSize)
	for b.Loop() {
		for i := range mapSize {
			mapValues[i] = gm[i]
		}
	}
}

func BenchmarkMapPreallocated(b *testing.B) {
	// NOTE: preallocation of wrong size may improve or degrade performance
	gm := make(map[uint32]uint32, mapSize/100)

	// var i uint32
	for i := range mapSize {
		gm[i] = i
	}

	mapValues := make([]uint32, mapSize)
	for b.Loop() {
		for i := range mapSize {
			mapValues[i] = gm[i]
		}
	}
}

func BenchmarkSyncMap(b *testing.B) {
	gm := sync.Map{}

	var i uint32
	for i = range mapSize {
		gm.Store(i, i)
	}

	mapValues := make([]any, mapSize)
	for b.Loop() {
		for i := range mapSize {
			mapValues[i], _ = gm.Load(i)
		}
	}
}

func BenchmarkSyncMapGeneric(b *testing.B) {
	gm := SyncMap[uint32, uint32]{}

	var i uint32
	for i = range mapSize {
		gm.Store(i, i)
	}
	mapValues := make([]uint32, mapSize)

	b.RunParallel(func(pb *testing.PB) {
		// for b.Loop() {
		for pb.Next() {
			for i := range mapSize {
				mapValues[i], _ = gm.Load(i)
			}
		}
		// }
	})
}

func workerSyncMap(syncMap *sync.Map, mapRange int, results chan<- any) {
	for i := range mapRange {
		val, _ := syncMap.Load(i)
		results <- val
	}
}

func workerSyncMapGeneric(syncMap *SyncMap[uint32, uint32], start, end uint32, results chan<- uint32, wg *sync.WaitGroup) {
	for i := start; i < end; i++ {
		syncMap.Load(i)
		val, _ := syncMap.Load(i)
		results <- val
	}
	wg.Done()
}

func BenchmarkSyncMapGenericWithWorkers(b *testing.B) {
	// b.StopTimer()
	gm := SyncMap[uint32, uint32]{}
	mapValues := make([]uint32, mapSize)
	var i uint32
	for i = range mapSize {
		gm.Store(i, i)
	}

	results := make(chan uint32, mapSize)
	workerCount := 10
	var wg sync.WaitGroup
	wg.Add(workerCount)
	workerRange := mapSize / uint32(workerCount)

	// for b.Loop() {
	var start uint32 = 0
	var end uint32 = 0

	// b.StartTimer()
	for range workerCount {
		end += workerRange
		go workerSyncMapGeneric(&gm, start, end, results, &wg)
		start += workerRange
	}
	// b.StopTimer()

	wg.Wait()
	// log.Println("Completed")
	for range results {
		// mapValues[i] = <-results
		<-results
	}
	close(results)
	log.Println("mapValues: ", mapValues)
	// }
}
