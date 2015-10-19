package main

import(
	"strconv"
	"os"
	"fmt"
	"math/rand"
	"math"
)

func intToBinaryString(i int) string {
	//create header
	// var value = int64(i)
	// binary := strconv.FormatInt(value, 2)
	// var diff = 8 - len(binary)
	// for i := 0; i < diff; i++ {
	//     binary = "0" + binary;
	// }
	// return binary

	var header = strconv.FormatInt(int64(i), 16)

	var diff = 5 - len(header)
	for i := 0; i < diff; i++ {
		header = "0" + header
	}

	return header

}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

//get angle between 2 vectors
func getAngleBetween2Vectors(p1 Vector2, p2 Vector2) float64 {
	var deltaY = p2.y - p1.y
	var deltaX = p2.x - p1.x
	var angleInDegrees = math.Atan2(deltaY, deltaX) * 180 / math.Pi
	return angleInDegrees
}

func distance(p1 Vector2, p2 Vector2) float64 {
	dist := math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
	return dist
}

func randomFloat(min, max float64) float64 {
	//return rand.Intn(max - min) + min
	return min + (rand.Float64() * ((max - min) + 1))
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func randFloatInRange(min float64, max float64) float64 {
	return (rand.Float64() * (max - min)) + min
}

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}