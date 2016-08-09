package main

import "math"

func (e *Entity) detectCollisions() bool {
	var keys = e.getNearbyKeys(2)
	e.body.updatePoints()

	for _, key := range keys {
		for _, otherE := range m[key] {
			if otherE == e {
				continue
			}
			otherE.body.updatePoints()
			if e.isCollidingWith(otherE) {
				//call on collide functions
				if e.onCollide != nil {
					e.onCollide(otherE)
				}
				if otherE.onCollide != nil {
					otherE.onCollide(e)
				}
				//return true
				return true
			}
		}
	}

	return false
}

func (e *Entity) isCollidingWith(e2 *Entity) bool {
	//e2 = bullet
	//e = player
	val := compareRects(e.body, e2.body)
	if val == true {
		//println((firstValue != v.origin && v != firstValue.origin) && (firstValue.value != 0 && v.value != 0) && (firstValue.health != 0 && v.health != 0))
		if (e != e2.origin && e2 != e.origin) && (e.active && e2.active) && (e.Health() != 0 && e2.Health() != 0) {
			return true
		}
	}
	return false
}

func compareRects(a Body, b Body) bool {
	polygons := [2]Body{a, b}
	inf := float64(9999999)

	for _, polygon := range polygons {
		for i1 := 0; i1 < len(polygon.points); i1++ {
			i2 := (i1 + 1) % len(polygon.points)
			p1 := polygon.points[i1]
			p2 := polygon.points[i2]

			normal := Vect2{x: p1.x - p2.x, y: p2.y - p1.y}

			var minA float64 = inf
			var maxA float64 = inf

			for _, p := range a.points {
				projected := normal.x*p.x + normal.y*p.y
				if minA == inf || projected < minA {
					minA = projected
				}
				if maxA == inf || projected > maxA {
					maxA = projected

				}
			}

			var minB float64 = inf
			var maxB float64 = inf
			for _, p := range b.points {
				projected := normal.x*p.x + normal.y*p.y
				if minB == inf || projected < minB {
					minB = projected
				}
				if maxB == inf || projected > maxB {
					maxB = projected
				}
			}

			if maxA < minB || maxB < minA {
				return false
			}

		}
	}

	return true
}

func rotatePoint(p Vect2, angle float64) Vect2 {
	var newP Vect2
	newP.x = (p.x * (math.Cos(angle))) - (p.y * (math.Sin(angle)))
	newP.y = (p.x * (math.Sin(angle))) - (p.y * (math.Cos(angle)))

	return newP
}

func (b *Body) updatePoints() {
	var w2 = b.size.x / 2
	var h2 = b.size.y / 2

	b.points[0].x = -w2
	b.points[0].y = -h2
	b.points[0].rotate(b.angle)
	b.points[0].move(b.pos.x, b.pos.y)

	b.points[1].x = w2
	b.points[1].y = -h2
	b.points[1].rotate(b.angle)
	b.points[1].move(b.pos.x, b.pos.y)

	b.points[2].x = w2
	b.points[2].y = h2
	b.points[2].rotate(b.angle)
	b.points[2].move(b.pos.x, b.pos.y)

	b.points[3].x = -w2
	b.points[3].y = h2
	b.points[3].rotate(b.angle)
	b.points[3].move(b.pos.x, b.pos.y)
}

func (p *Vect2) move(x float64, y float64) {
	p.x += x
	p.y += y
}

func (p *Vect2) rotate(angle float64) {
	x := p.x
	y := p.y
	p.x = (x * math.Cos(angle)) - (y * math.Sin(angle))
	p.y = (x * math.Sin(angle)) - (y * math.Cos(angle))
}
