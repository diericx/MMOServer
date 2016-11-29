package main

type Vect2 struct {
	x float64
	y float64
}

type Vect3 struct {
	x float64
	y float64
	z float64
}

type Body struct {
	pos        Vect3
	targetPos  Vect3
	size       Vect3
	vel        Vect3
	continuous bool
	angle      float64
	points     []Vect2
}

func (b Body) Size(x float64, y float64) {
	b.size.x = x
	b.size.y = x
}
