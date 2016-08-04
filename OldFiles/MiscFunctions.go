package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

type Vector2 struct {
	x float64
	y float64
}

type Point struct {
	x float64
	y float64
}

type Rectangle struct {
	pos      Vector2
	size     Vector2
	rotation int
	points   []Point
}

func createRect(x float64, y float64, width float64, height float64) Rectangle {

	var newRect Rectangle

	newRect.pos = Vector2{x: x, y: y}
	newRect.size = Vector2{x: width, y: height}

	var w2 = width / 2
	var h2 = height / 2

	var point1 Point
	point1.x = -w2
	point1.y = -h2
	newRect.points = append(newRect.points, point1)

	var point2 Point
	point2.x = w2
	point2.y = -h2
	newRect.points = append(newRect.points, point2)

	var point3 Point
	point3.x = w2
	point3.y = h2
	newRect.points = append(newRect.points, point3)

	var point4 Point
	point4.x = -w2
	point4.y = h2
	newRect.points = append(newRect.points, point4)

	return newRect

}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func distance(p1 Vector2, p2 Vector2) float64 {
	dist := math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
	return dist
}

//get angle between 2 vectors
func getAngleBetween2Vectors(p1 Vector2, p2 Vector2) float64 {
	var deltaY = p2.y - p1.y
	var deltaX = p2.x - p1.x
	var angleInDegrees = math.Atan2(deltaY, deltaX) * 180 / math.Pi
	return angleInDegrees
}

func randFloatInRange(min float64, max float64) float64 {
	return (rand.Float64() * (max - min)) + min
}
