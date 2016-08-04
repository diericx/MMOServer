package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/vova616/chipmunk"
	//"github.com/vova616/chipmunk/vect"
	"log"
	"math"
	//"math/rand"
	"os"
	"runtime"
	"time"
)

// drawCircle draws a circle for the specified radius, rotation angle, and the specified number of sides
func drawCircle(radius float64, sides int) {
	gl.Begin(gl.LINE_LOOP)
	for a := 0.0; a < 2*math.Pi; a += (2 * math.Pi / float64(sides)) {
		gl.Vertex2d(math.Sin(a)*radius, math.Cos(a)*radius)
	}
	gl.Vertex3f(0, 0, 0)
	gl.End()
}

func drawEntity(e *Entity) {
	gl.Begin(gl.LINE_LOOP)
	//pos := e.body.Position()
	size := e.size

	gl.Vertex2d(float64(-(size.X / 2)), float64((size.Y / 2)))
	gl.Vertex2d(float64((size.X / 2)), float64((size.Y / 2)))
	gl.Vertex2d(float64((size.X / 2)), float64(-(size.Y / 2)))
	gl.Vertex2d(float64(-(size.X / 2)), float64(-(size.Y / 2)))

	//gl.Vertex3f(0, 0, 0)
	gl.End()
}

// OpenGL draw function
func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.LoadIdentity()

	gl.Begin(gl.LINES)
	gl.Color3f(.2, .2, .2)
	gl.End()

	gl.Color4f(.3, .3, 1, .8)
	// draw players
	for _, e := range entities {
		gl.PushMatrix()
		pos := e.body.Position()
		rot := e.body.Angle() * chipmunk.DegreeConst
		gl.Translatef(float32(pos.X), float32(pos.Y), 0.0)
		gl.Rotatef(float32(rot), 0, 0, 1)

		//drawCircle(float64(ballRadius), 60)
		drawEntity(&e)
		gl.PopMatrix()
	}
}

// onResize sets up a simple 2d ortho context based on the window size
func onResize(window *glfw.Window, w, h int) {
	w, h = window.GetSize() // query window to get screen pixels
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(1, 1, 1, 1)
}

func InitializeDisplay() {
	runtime.LockOSThread()

	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize GLFW: ", err)
	}
	defer glfw.Terminate()

	// create window
	window, err := glfw.CreateWindow(600, 600, os.Args[0], nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.SetFramebufferSizeCallback(onResize)
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	// set up opengl context
	onResize(window, 600, 600)

	// set up physics
	//createBodies()

	runtime.LockOSThread()
	glfw.SwapInterval(1)

	//ticksToNextBall := 10
	ticker := time.NewTicker(time.Second / 60)
	for !window.ShouldClose() {
		draw()
		window.SwapBuffers()
		glfw.PollEvents()

		<-ticker.C // wait up to 1/60th of a second
	}
}
