package datastructures

type SwapbackArray[T any] []T

// func New[T any]() *SwapbackArray[T] {
// 	return &SwapbackArray[T]{
// 		sa: make([]T, 100),
// 	}
// }

func New[T any](capacity int) *SwapbackArray[T] {
	sa := make(SwapbackArray[T], 0, capacity)
	return &sa
}

func (sa *SwapbackArray[T]) RemoveIndexBySwap(index int) bool {
	saLen := len(*sa)
	if index < 0 || index >= saLen {
		return false
	}

	(*sa)[index] = (*sa)[saLen-1]

	var zero T
	(*sa)[saLen-1] = zero

	*sa = (*sa)[:saLen-1]

	return true
}
