package datastructures

import (
	"log"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	testData := []uint32{0, 1, 2, 3, 4, 5}
	testArr := NewOptimizedVec(testData)

	log.Println(testArr.data)
	if !reflect.DeepEqual(testArr.data, testData) {
		t.Fail()
	}
	log.Println(testArr.index)
	log.Println(testArr.ids)
}

func TestDelete(t *testing.T) {
	testData := []uint32{0, 1, 2, 3, 4, 5}
	testArr := NewOptimizedVec(testData)
	testArr.Delete(3)
	log.Println(testArr.data)
	log.Println(testArr.index)
	log.Println(testArr.ids)
}
