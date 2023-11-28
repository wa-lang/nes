package nes

const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Controller struct {
	buttons [8]bool
	index   byte
	strobe  byte
}

func NewController() *Controller {
	return &Controller{}
}

func (this *Controller) SetButtons(buttons [8]bool) {
	this.buttons = buttons
}

func (this *Controller) Read() byte {
	value := byte(0)
	if this.index < 8 && this.buttons[this.index] {
		value = 1
	}
	this.index++
	if this.strobe&1 == 1 {
		this.index = 0
	}
	return value
}

func (this *Controller) Write(value byte) {
	this.strobe = value
	if this.strobe&1 == 1 {
		this.index = 0
	}
}
