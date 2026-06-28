package datastructures

type OptimizedVec[T any] struct {
	index []uint32
	ids   []uint32
	data  []T
}

func NewOptimizedVec[T any](data []T) *OptimizedVec[T] {
	indexes := make([]uint32, len(data))
	for i := range indexes {
		indexes[i] = uint32(i)
	}
	return &OptimizedVec[T]{
		index: indexes,
		ids:   indexes,
		data:  data,
	}
}

func (ov *OptimizedVec[T]) Find(index uint32) T {
	indx := ov.index[index]
	return ov.data[indx]
}

func (ov *OptimizedVec[T]) Add(val T) {
	lenIndex := len(ov.index) - 1
	lenData := len(ov.data) - 1
	ov.data = append(ov.data, val)
	if lenData < lenIndex {
		return
	}
	ov.index = append(ov.index, uint32(lenIndex))
	ov.ids = append(ov.index, uint32(lenIndex))
}

func (ov *OptimizedVec[T]) Swap(index1, index2 uint32) {
	indx1 := ov.index[index1]
	indx2 := ov.index[index2]
	data1 := ov.data[indx1]
	ov.data[indx1] = ov.data[indx2]
	ov.data[indx2] = data1
}

func (ov *OptimizedVec[T]) Delete(index uint32) {
	lenght := len(ov.data) - 1
	indx := ov.index[index]
	ov.data[indx] = ov.data[lenght]
	// ov.Swap(index, uint32(lenght))

	ov.data = ov.data[:lenght]
}
