package main

import (
	"math"
)

var bulletsToRemove = make([]*Bullet, 0)

func compareRects(objRect Rectangle, bulletRect Rectangle) bool {

	var pRectRot Rectangle = objRect
	rotateRectsPoints(pRectRot, (float64(objRect.rotation)/180.0)*3.14159)

	var bRectRot Rectangle = bulletRect
	rotateRectsPoints(bRectRot, (float64(bulletRect.rotation)/180.0)*3.14159)

	//CHECK X

	var v Point
	v.x = pRectRot.width / 2
	v.y = 0
	v = rotatePoint(v, (float64(objRect.rotation)/180.0)*3.14159)

	v = normalize(v)

	var a Point = pRectRot.points[2]

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

	var c Point
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

	var w Point
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

func addBulletToRemoveList(b *Bullet) {
	bulletsToRemove = append(bulletsToRemove, b)
}

func removeBulletFromList(bullets *[]*Bullet, b *Bullet) {
	var i = 0
	var foundIndex = -1
	var bulletsCopy = *bullets
	for _, bullet := range *bullets {
		if b == bullet {
			foundIndex = i
		}
		i++
	}
	
	if foundIndex != -1 {
		*bullets = append(bulletsCopy[:foundIndex], bulletsCopy[foundIndex+1:]...)
	}
}

func clearBulletRemoveList(bullets *[]*Bullet) {
	for _, bullet := range bulletsToRemove {
		removeBulletFromList(bullets, bullet)
	}
	bulletsToRemove = make([]*Bullet, 0)
}

func rotateRectsPoints(r Rectangle, angle float64) {
	for _, p := range r.points {
		p = rotatePoint(p, angle)
	}
}

func dotProduct(pointA Point, pointB Point) float64 {
	return math.Abs(pointA.x*pointB.x) + math.Abs(pointA.y*pointB.y)
}

func movePoints(rect Rectangle) {
	for _, p := range rect.points {
		p.x = p.x + rect.x
		p.y = p.y + rect.y
	}
}

func normalize(p Point) Point {
	var magnitude = math.Sqrt(p.x*p.x + p.y*p.y)
	if magnitude > 0 {
		p.x = p.x / magnitude
		p.y = p.y / magnitude
	}
	return p
}

func createRect(x float64, y float64, width float64, height float64) Rectangle {

	var newRect Rectangle

	newRect.x = x
	newRect.y = y
	newRect.width = width
	newRect.height = height

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

func rotatePoint(p Point, angle float64) Point {
	var newP Point
	newP.x = (p.x * math.Cos(angle)) - (p.y * math.Sin(angle))
	newP.y = (p.x * math.Sin(angle)) - (p.y * math.Cos(angle))

	return newP
}
