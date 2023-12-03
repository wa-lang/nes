package nes

import (
	"image"
)

var Consol_Mapper Mapper
var nes_RAM []byte

var Nes_Controller1 *Controller
var Nes_Controller2 *Controller

func InitConsole(romBytes []byte) error {

	//file_log, _ = os.Create("C:\\Users\\Ending\\Desktop\\g.txt")

	var err error
	err = LoadNESFile(romBytes)
	if err != nil {
		return err
	}
	nes_RAM = make([]byte, 2048)
	Nes_Controller1 = NewController()
	Nes_Controller2 = NewController()
	mapper, err := NewMapper()
	if err != nil {
		return err
	}
	Consol_Mapper = mapper
	CPU_Init()
	PPU_initNesPPU()
	return nil
}

func Consloe_Reset() {
	CPU_Reset()
	PPU_Reset()
}

func Consloe_Step() int {
	cpuCycles := CPU_Step()
	ppuCycles := cpuCycles * 3
	for i := 0; i < ppuCycles; i++ {
		PPU_Step()
		Consol_Mapper.Step()
	}
	return cpuCycles
}

func Consloe_StepFrame() int {
	cpuCycles := 0
	frame := PPU_Frame
	for frame == PPU_Frame {
		cpuCycles += Consloe_Step()
	}
	return cpuCycles
}

func Consloe_StepSeconds(seconds float64) {
	cycles := int(CPUFrequency * seconds)
	for cycles > 0 {
		cycles -= Consloe_Step()
	}
}

func Consloe_Buffer() *image.RGBA {
	return PPU_front
}

func Consloe_SetButtons1(buttons [8]bool) {
	Nes_Controller1.SetButtons(buttons)
}

func Consloe_SetButtons2(buttons [8]bool) {
	Nes_Controller2.SetButtons(buttons)
}
