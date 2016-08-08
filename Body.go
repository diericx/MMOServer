package main

type Body struct {
	pos        Vect2
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

func (b Body) Position() Vect2 {
	return b.pos
}
