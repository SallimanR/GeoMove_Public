package main

const lenght uint32 = 100000

type Point[T any] struct {
	x T
	y T
}

func main() {
	arr := make([]Point[uint64], 0, lenght)
	for i := range uint64(lenght) {
		arr = append(arr, Point[uint64]{x: i, y: i})
	}

	result := make([]uint64, lenght)
	for i := range arr {
		result = append(result, arr[i].x)
	}
	// log.Println(result)
}
