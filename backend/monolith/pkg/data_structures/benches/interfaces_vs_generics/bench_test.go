package interfacesvsgenerics

import "testing"

type Doer interface {
	Do() int
}

type myDoer struct {
	val int
}

func (m myDoer) Do() int {
	return m.val
}

var Sink int // prevent dead code elimination

func BenchmarkGeneric(b *testing.B) {
	items := make([]myDoer, 1000)
	for i := range items {
		items[i] = myDoer{val: i}
	}
	b.ResetTimer()
	for b.Loop() {
		Sink = ProcessGeneric(items)
	}
}

func BenchmarkInterface(b *testing.B) {
	items := make([]Doer, 1000)
	for i := range items {
		items[i] = myDoer{val: i}
	}
	b.ResetTimer()
	for b.Loop() {
		Sink = ProcessInterface(items)
	}
}

func ProcessGeneric[T Doer](items []T) int {
	var sum int
	for _, item := range items {
		sum += item.Do()
	}
	return sum
}

func ProcessInterface(items []Doer) int {
	var sum int
	for _, item := range items {
		sum += item.Do()
	}
	return sum
}

// type Shape struct {
// 	x float64
// 	y float64
// }
//
// type ShapeMethods interface {
// 	Area() float64
// }
//
// func (s *Shape) Area() float64 {
// 	return s.x * s.y
// }
//
// func NewShape(x, y float64) ShapeMethods {
// 	return &Shape{x: x, y: y}
// }
//
// func BenchmarkInterface(b *testing.B) {
// 	shape := NewShape(35.78, 31.12)
// 	for b.Loop() {
// 		shape.Area()
// 	}
// }
//
// type ShapeGeneric[T any] struct {
// 	x float64
// 	y float64
// }
//
// func (s *ShapeGeneric[T]) Area() float64 {
// 	return s.x * s.y
// }
//
// func BenchmarkGeneric(b *testing.B) {
// 	shape := ShapeGeneric[float64]{x: 35.78, y: 31.12}
// 	for b.Loop() {
// 		shape.Area()
// 	}
// }
