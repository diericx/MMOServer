package main

import (
	"math"
)

func compareRects(objRect rectangle, bulletRect rectangle) bool {

	var pRectRot rectangle = objRect
	rotateRectsPoints(pRectRot, (float64(objRect.rotation)/180.0)*3.14159)

	var bRectRot rectangle = bulletRect
	rotateRectsPoints(bRectRot, (float64(bulletRect.rotation)/180.0)*3.14159)

	//CHECK X

	var v point
	v.x = pRectRot.width / 2
	v.y = 0
	v = rotatePoint(v, (float64(objRect.rotation)/180.0)*3.14159)

	v = normalize(v)

	var a point = pRectRot.points[2]

	var av float64 = dotProduct(a, v)

	var rv float64 = 0

	for i := 0; i < 4; i++ {
		var poop float64 = dotProduct(bRectRot.points[i], v)
		if poop > rv {
			rv = poop
		}
	}

	movePoints(pRectRot)
	movePoints(bRectRot)

	var c point
	c.x = math.Abs(pRectRot.x - bRectRot.x)
	c.y = math.Abs(pRectRot.y - bRectRot.y)

	var cv float64 = dotProduct(c, v)
	var result bool = cv > av+rv

	//fmt.Printf("%t\n", result)

	//CHECK Y
	pRectRot = objRect
	rotateRectsPoints(pRectRot, (float64(objRect.rotation)/180.0)*3.14159)

	bRectRot = bulletRect
	rotateRectsPoints(bRectRot, (float64(bulletRect.rotation)/180.0)*3.14159)

	var w point
	w.y = pRectRot.height / 2
	w.x = 0
	w = rotatePoint(w, (float64(objRect.rotation)/180.0)*3.14159)

	w = normalize(w)

	var a2 = pRectRot.points[2]

	var aw float64 = dotProduct(a2, w)

	var rw float64 = 0

	for i := 0; i < 4; i++ {
		var poop float64 = dotProduct(bRectRot.points[i], w)
		if poop > rw {
			rw = poop
		}
	}

	var cw float64 = dotProduct(c, w)
	var result2 bool = cw > aw+rw

	// fmt.Printf("%v %v \n", cw, aw + rw)

	// fmt.Printf("%t\n", result2)

	return !result && !result2
}

func rotateRectsPoints(r rectangle, angle float64) {
	for _, p := range r.points {
		p = rotatePoint(p, angle)
	}
}

func dotProduct(pointA point, pointB point) float64 {
	return math.Abs(pointA.x*pointB.x) + math.Abs(pointA.y*pointB.y)
}

func movePoints(rect rectangle) {
	for _, p := range rect.points {
		p.x = p.x + rect.x
		p.y = p.y + rect.y
	}
}
