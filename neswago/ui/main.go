package ui

import (
	"log"
	"runtime"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

func Main(name string, romBytes []byte) {
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	app := NewAppWindow(name, 256*3, 240*3)
	app.Run(name, romBytes)
}

type AppWindow struct {
	window    *glfw.Window
	console   *nes.Console
	texture   uint32
	timestamp float64
}

func NewAppWindow(name string, width, height int) *AppWindow {
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	window.MakeContextCurrent()
	window.SetTitle(name)

	// initialize gl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	this := &AppWindow{}
	this.window = window
	return this
}

func (d *AppWindow) Run(name string, romBytes []byte) {
	console, err := nes.NewConsole(romBytes)
	if err != nil {
		log.Fatalln(err)
	}

	d.console = console
	d.texture = createTexture()

	d.window.SetKeyCallback(d.OnKey)

	defer func() {
		d.window.SetKeyCallback(nil)
	}()

	gl.ClearColor(0, 0, 0, 1)
	for !d.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		timestamp := glfw.GetTime()
		dt := timestamp - d.timestamp
		d.timestamp = timestamp

		d.Update(timestamp, dt)

		d.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (d *AppWindow) OnKey(
	window *glfw.Window,
	key glfw.Key, scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if action == glfw.Press {
		if key >= glfw.Key0 && key <= glfw.Key9 {
			// snapshot := int(key - glfw.Key0)
			if mods&glfw.ModShift == 0 {
				// view.load(snapshot)
			} else {
				// view.save(snapshot)
			}
		}
		switch key {
		case glfw.KeySpace:
			// screenshot(view.console.Buffer())
		case glfw.KeyR:
			d.console.Reset()
		case glfw.KeyTab:
			//
		}
	}
}

func (d *AppWindow) Update(t, dt float64) {
	if dt > 1 {
		dt = 0
	}
	window := d.window

	if joystickReset(glfw.Joystick1) {
		// view.director.ShowMenu()
	}
	if joystickReset(glfw.Joystick2) {
		// view.director.ShowMenu()
	}
	if readKey(window, glfw.KeyEscape) {
		// view.director.ShowMenu()
	}

	// 处理键盘信息
	{
		turbo := d.console.PPU.Frame%6 < 3
		k1 := readKeys(window, turbo)
		j1 := readJoystick(glfw.Joystick1, turbo)
		j2 := readJoystick(glfw.Joystick2, turbo)
		d.console.SetButtons1(combineButtons(k1, j1))
		d.console.SetButtons2(j2)
	}

	// 执行 NES
	d.console.StepSeconds(dt)

	// 重要: 复制NES图形缓冲并显示!!!

	gl.BindTexture(gl.TEXTURE_2D, d.texture)
	setTexture(d.console.Buffer())
	d.drawBuffer()
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (d *AppWindow) drawBuffer() {
	w, h := d.window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}
