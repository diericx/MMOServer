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
	pos        Vect2
	targetPos  Vect2
	size       Vect2
	vel        Vect2
	continuous bool
	angle      float64
	points     []Vect2
}

func (b Body) Size(x float64, y float64) {
	b.size.x = x
	b.size.y = x
}
