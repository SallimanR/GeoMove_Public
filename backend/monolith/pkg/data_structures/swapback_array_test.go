package datastructures

import (
	"reflect"
	"testing"
)

func TestRemoveIndexBySwap(t *testing.T) {
	size := 1000

	arr := make([]int, size)
	for i := range arr {
		arr[i] = i
	}
	arr[0] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]

	correctArr := make([]int, size)
	for i := range correctArr {
		correctArr[i] = i
	}
	correctArr[0] = correctArr[size-1]
	correctArr = correctArr[:len(correctArr)-1]

	// log.Println("arr: ", arr)
	// log.Println("correctArr: ", correctArr)
	if !reflect.DeepEqual(arr, correctArr) {
		t.Fail()
	}
}

var n = 20000

func BenchmarkRemoveIndexBySwap(b *testing.B) {
	b.StopTimer()
	// arr := *New[int](1000)
	arr := make([]int, n/2)
	// arrLen := len(arr)

	for i := range arr {
		arr[i] = i
	}

	b.StartTimer()
	for i := n / 2; i < n; i++ {
		// arr.RemoveIndexBySwap(0)
		// arr = append(arr, i)
		arr[0] = arr[len(arr)-1]
		arr = arr[:len(arr)-1]
		arr = append(arr, i)
	}

	// log.Println("BySwap: ", arr[0], arr[len(arr)-1])
	// log.Println("BySwap: ", arr[1], arr[len(arr)-2])
}

func BenchmarkRemoveIndex(b *testing.B) {
	b.StopTimer()
	arr := make([]int, n/2)
	// arrLen := len(arr)

	for i := range arr {
		arr[i] = i
	}

	b.StartTimer()
	for i := n / 2; i < n; i++ {
		arr = append(arr[:0], arr[1:]...)
		arr = append(arr, i)
	}

	// log.Println("ByRegular: ", arr[0], arr[len(arr)-1])
	// log.Println("ByRegular: ", arr[1], arr[len(arr)-2])
}
