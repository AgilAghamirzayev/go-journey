package main

import "fmt"

type Shape interface {
	Draw()
}

type Circle struct{}

func (c *Circle) Draw() {
	fmt.Println("Drawing a Circle")
}

type Square struct{}

func (s *Square) Draw() {
	fmt.Println("Drawing a Square")
}

func ShapeFactory(shapeType string) Shape {
	switch shapeType {
	case "circle":
		return &Circle{}
	case "square":
		return &Square{}
	default:
		return nil
	}
}

func main() {
	shape1 := ShapeFactory("circle")
	shape1.Draw()

	shape2 := ShapeFactory("square")
	shape2.Draw()

	shape3 := ShapeFactory("triangle")
	if shape3 == nil {
		fmt.Println("Unknown shape type")
	}
}
