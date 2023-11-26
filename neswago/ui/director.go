package ui

import (
	"log"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Director struct {
	window    *glfw.Window
	gameView  *GameView
	timestamp float64
}

func NewDirector(window *glfw.Window) *Director {
	director := Director{}
	director.window = window
	return &director
}

func (d *Director) SetTitle(title string) {
	d.window.SetTitle(title)
}

func (d *Director) Run(name string, romBytes []byte) {
	d.playGame(name, romBytes)

	for !d.window.ShouldClose() {
		d.stepRun()
		d.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (d *Director) playGame(name string, romBytes []byte) {
	console, err := nes.NewConsole(name, romBytes)
	if err != nil {
		log.Fatalln(err)
	}

	d.gameView = NewGameView(d, console, name)
	d.gameView.Enter()
	d.timestamp = glfw.GetTime()
}

func (d *Director) stepRun() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	timestamp := glfw.GetTime()
	dt := timestamp - d.timestamp
	d.timestamp = timestamp
	if d.gameView != nil {
		d.gameView.Update(timestamp, dt)
	}
}
