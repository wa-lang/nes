package nes

const frameCounterRate = CPUFrequency / 240.0

var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

var dutyTable = [][]byte{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

var triangleTable = []byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

var noiseTable = []uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}

var dmcTable = []byte{
	214, 190, 170, 160, 143, 127, 113, 107, 95, 80, 71, 64, 53, 42, 36, 27,
}

var pulseTable [31]float32
var tndTable [203]float32

func init() {
	for i := 0; i < 31; i++ {
		pulseTable[i] = 95.52 / (8128.0/float32(i) + 100)
	}
	for i := 0; i < 203; i++ {
		tndTable[i] = 163.67 / (24329.0/float32(i) + 100)
	}
}

// APU

type APU struct {
	console     *Console
	sampleRate  float64
	pulse1      Pulse
	pulse2      Pulse
	triangle    Triangle
	noise       Noise
	dmc         DMC
	cycle       uint64
	framePeriod byte
	frameValue  byte
	frameIRQ    bool
}

func NewAPU(console *Console) *APU {
	this := &APU{}
	this.console = console
	this.noise.shiftRegister = 1
	this.pulse1.channel = 1
	this.pulse2.channel = 2
	this.framePeriod = 4
	this.dmc.cpu = console.CPU
	return this
}

func (this *APU) Step() {
	return
	cycle1 := this.cycle
	this.cycle++
	cycle2 := this.cycle
	this.stepTimer()
	f1 := int(float64(cycle1) / frameCounterRate)
	f2 := int(float64(cycle2) / frameCounterRate)
	if f1 != f2 {
		this.stepFrameCounter()
	}
	s1 := int(float64(cycle1) / this.sampleRate)
	s2 := int(float64(cycle2) / this.sampleRate)
	if s1 != s2 {
		// todo 发送数据
	}
}

// mode 0:    mode 1:       function
// ---------  -----------  -----------------------------
//   - - - f    - - - - -    IRQ (if bit 6 is clear)
//   - l - l    l - l - -    Length counter and sweep
//     e e e e    e e e e -    Envelope and linear counter
func (this *APU) stepFrameCounter() {
	switch this.framePeriod {
	case 4:
		this.frameValue = (this.frameValue + 1) % 4
		switch this.frameValue {
		case 0, 2:
			this.stepEnvelope()
		case 1:
			this.stepEnvelope()
			this.stepSweep()
			this.stepLength()
		case 3:
			this.stepEnvelope()
			this.stepSweep()
			this.stepLength()
			this.fireIRQ()
		}
	case 5:
		this.frameValue = (this.frameValue + 1) % 5
		switch this.frameValue {
		case 0, 2:
			this.stepEnvelope()
		case 1, 3:
			this.stepEnvelope()
			this.stepSweep()
			this.stepLength()
		}
	}
}

func (this *APU) stepTimer() {
	if this.cycle%2 == 0 {
		this.pulse1.stepTimer()
		this.pulse2.stepTimer()
		this.noise.stepTimer()
		this.dmc.stepTimer()
	}
	this.triangle.stepTimer()
}

func (this *APU) stepEnvelope() {
	this.pulse1.stepEnvelope()
	this.pulse2.stepEnvelope()
	this.triangle.stepCounter()
	this.noise.stepEnvelope()
}

func (this *APU) stepSweep() {
	this.pulse1.stepSweep()
	this.pulse2.stepSweep()
}

func (this *APU) stepLength() {
	this.pulse1.stepLength()
	this.pulse2.stepLength()
	this.triangle.stepLength()
	this.noise.stepLength()
}

func (this *APU) fireIRQ() {
	if this.frameIRQ {
		this.console.CPU.triggerIRQ()
	}
}

func (this *APU) readRegister(address uint16) byte {
	switch address {
	case 0x4015:
		return this.readStatus()
		// default:
		// 	log.Fatalf("unhandled apu register read at address: 0x%04X", address)
	}
	return 0
}

func (this *APU) writeRegister(address uint16, value byte) {
	switch address {
	case 0x4000:
		this.pulse1.writeControl(value)
	case 0x4001:
		this.pulse1.writeSweep(value)
	case 0x4002:
		this.pulse1.writeTimerLow(value)
	case 0x4003:
		this.pulse1.writeTimerHigh(value)
	case 0x4004:
		this.pulse2.writeControl(value)
	case 0x4005:
		this.pulse2.writeSweep(value)
	case 0x4006:
		this.pulse2.writeTimerLow(value)
	case 0x4007:
		this.pulse2.writeTimerHigh(value)
	case 0x4008:
		this.triangle.writeControl(value)
	case 0x4009:
	case 0x4010:
		this.dmc.writeControl(value)
	case 0x4011:
		this.dmc.writeValue(value)
	case 0x4012:
		this.dmc.writeAddress(value)
	case 0x4013:
		this.dmc.writeLength(value)
	case 0x400A:
		this.triangle.writeTimerLow(value)
	case 0x400B:
		this.triangle.writeTimerHigh(value)
	case 0x400C:
		this.noise.writeControl(value)
	case 0x400D:
	case 0x400E:
		this.noise.writePeriod(value)
	case 0x400F:
		this.noise.writeLength(value)
	case 0x4015:
		this.writeControl(value)
	case 0x4017:
		this.writeFrameCounter(value)
		// default:
		// 	log.Fatalf("unhandled apu register write at address: 0x%04X", address)
	}
}

func (this *APU) readStatus() byte {
	var result byte
	if this.pulse1.lengthValue > 0 {
		result |= 1
	}
	if this.pulse2.lengthValue > 0 {
		result |= 2
	}
	if this.triangle.lengthValue > 0 {
		result |= 4
	}
	if this.noise.lengthValue > 0 {
		result |= 8
	}
	if this.dmc.currentLength > 0 {
		result |= 16
	}
	return result
}

func (this *APU) writeControl(value byte) {
	this.pulse1.enabled = value&1 == 1
	this.pulse2.enabled = value&2 == 2
	this.triangle.enabled = value&4 == 4
	this.noise.enabled = value&8 == 8
	this.dmc.enabled = value&16 == 16
	if !this.pulse1.enabled {
		this.pulse1.lengthValue = 0
	}
	if !this.pulse2.enabled {
		this.pulse2.lengthValue = 0
	}
	if !this.triangle.enabled {
		this.triangle.lengthValue = 0
	}
	if !this.noise.enabled {
		this.noise.lengthValue = 0
	}
	if !this.dmc.enabled {
		this.dmc.currentLength = 0
	} else {
		if this.dmc.currentLength == 0 {
			this.dmc.restart()
		}
	}
}

func (this *APU) writeFrameCounter(value byte) {
	this.framePeriod = 4 + (value>>7)&1
	this.frameIRQ = (value>>6)&1 == 0
	// this.frameValue = 0
	if this.framePeriod == 5 {
		this.stepEnvelope()
		this.stepSweep()
		this.stepLength()
	}
}

// Pulse

type Pulse struct {
	enabled         bool
	channel         byte
	lengthEnabled   bool
	lengthValue     byte
	timerPeriod     uint16
	timerValue      uint16
	dutyMode        byte
	dutyValue       byte
	sweepReload     bool
	sweepEnabled    bool
	sweepNegate     bool
	sweepShift      byte
	sweepPeriod     byte
	sweepValue      byte
	envelopeEnabled bool
	envelopeLoop    bool
	envelopeStart   bool
	envelopePeriod  byte
	envelopeValue   byte
	envelopeVolume  byte
	constantVolume  byte
}

func (p *Pulse) writeControl(value byte) {
	p.dutyMode = (value >> 6) & 3
	p.lengthEnabled = (value>>5)&1 == 0
	p.envelopeLoop = (value>>5)&1 == 1
	p.envelopeEnabled = (value>>4)&1 == 0
	p.envelopePeriod = value & 15
	p.constantVolume = value & 15
	p.envelopeStart = true
}

func (p *Pulse) writeSweep(value byte) {
	p.sweepEnabled = (value>>7)&1 == 1
	p.sweepPeriod = (value>>4)&7 + 1
	p.sweepNegate = (value>>3)&1 == 1
	p.sweepShift = value & 7
	p.sweepReload = true
}

func (p *Pulse) writeTimerLow(value byte) {
	p.timerPeriod = (p.timerPeriod & 0xFF00) | uint16(value)
}

func (p *Pulse) writeTimerHigh(value byte) {
	p.lengthValue = lengthTable[value>>3]
	p.timerPeriod = (p.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
	p.envelopeStart = true
	p.dutyValue = 0
}

func (p *Pulse) stepTimer() {
	if p.timerValue == 0 {
		p.timerValue = p.timerPeriod
		p.dutyValue = (p.dutyValue + 1) % 8
	} else {
		p.timerValue--
	}
}

func (p *Pulse) stepEnvelope() {
	if p.envelopeStart {
		p.envelopeVolume = 15
		p.envelopeValue = p.envelopePeriod
		p.envelopeStart = false
	} else if p.envelopeValue > 0 {
		p.envelopeValue--
	} else {
		if p.envelopeVolume > 0 {
			p.envelopeVolume--
		} else if p.envelopeLoop {
			p.envelopeVolume = 15
		}
		p.envelopeValue = p.envelopePeriod
	}
}

func (p *Pulse) stepSweep() {
	if p.sweepReload {
		if p.sweepEnabled && p.sweepValue == 0 {
			p.sweep()
		}
		p.sweepValue = p.sweepPeriod
		p.sweepReload = false
	} else if p.sweepValue > 0 {
		p.sweepValue--
	} else {
		if p.sweepEnabled {
			p.sweep()
		}
		p.sweepValue = p.sweepPeriod
	}
}

func (p *Pulse) stepLength() {
	if p.lengthEnabled && p.lengthValue > 0 {
		p.lengthValue--
	}
}

func (p *Pulse) sweep() {
	delta := p.timerPeriod >> p.sweepShift
	if p.sweepNegate {
		p.timerPeriod -= delta
		if p.channel == 1 {
			p.timerPeriod--
		}
	} else {
		p.timerPeriod += delta
	}
}

// Triangle

type Triangle struct {
	enabled       bool
	lengthEnabled bool
	lengthValue   byte
	timerPeriod   uint16
	timerValue    uint16
	dutyValue     byte
	counterPeriod byte
	counterValue  byte
	counterReload bool
}

func (this *Triangle) writeControl(value byte) {
	this.lengthEnabled = (value>>7)&1 == 0
	this.counterPeriod = value & 0x7F
}

func (this *Triangle) writeTimerLow(value byte) {
	this.timerPeriod = (this.timerPeriod & 0xFF00) | uint16(value)
}

func (this *Triangle) writeTimerHigh(value byte) {
	this.lengthValue = lengthTable[value>>3]
	this.timerPeriod = (this.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
	this.timerValue = this.timerPeriod
	this.counterReload = true
}

func (this *Triangle) stepTimer() {
	if this.timerValue == 0 {
		this.timerValue = this.timerPeriod
		if this.lengthValue > 0 && this.counterValue > 0 {
			this.dutyValue = (this.dutyValue + 1) % 32
		}
	} else {
		this.timerValue--
	}
}

func (this *Triangle) stepLength() {
	if this.lengthEnabled && this.lengthValue > 0 {
		this.lengthValue--
	}
}

func (this *Triangle) stepCounter() {
	if this.counterReload {
		this.counterValue = this.counterPeriod
	} else if this.counterValue > 0 {
		this.counterValue--
	}
	if this.lengthEnabled {
		this.counterReload = false
	}
}

func (this *Triangle) output() byte {
	if !this.enabled {
		return 0
	}
	if this.timerPeriod < 3 {
		return 0
	}
	if this.lengthValue == 0 {
		return 0
	}
	if this.counterValue == 0 {
		return 0
	}
	return triangleTable[this.dutyValue]
}

// Noise

type Noise struct {
	enabled         bool
	mode            bool
	shiftRegister   uint16
	lengthEnabled   bool
	lengthValue     byte
	timerPeriod     uint16
	timerValue      uint16
	envelopeEnabled bool
	envelopeLoop    bool
	envelopeStart   bool
	envelopePeriod  byte
	envelopeValue   byte
	envelopeVolume  byte
	constantVolume  byte
}

func (this *Noise) writeControl(value byte) {
	this.lengthEnabled = (value>>5)&1 == 0
	this.envelopeLoop = (value>>5)&1 == 1
	this.envelopeEnabled = (value>>4)&1 == 0
	this.envelopePeriod = value & 15
	this.constantVolume = value & 15
	this.envelopeStart = true
}

func (this *Noise) writePeriod(value byte) {
	this.mode = value&0x80 == 0x80
	this.timerPeriod = noiseTable[value&0x0F]
}

func (this *Noise) writeLength(value byte) {
	this.lengthValue = lengthTable[value>>3]
	this.envelopeStart = true
}

func (this *Noise) stepTimer() {
	if this.timerValue == 0 {
		this.timerValue = this.timerPeriod
		var shift byte
		if this.mode {
			shift = 6
		} else {
			shift = 1
		}
		b1 := this.shiftRegister & 1
		b2 := (this.shiftRegister >> shift) & 1
		this.shiftRegister >>= 1
		this.shiftRegister |= (b1 ^ b2) << 14
	} else {
		this.timerValue--
	}
}

func (this *Noise) stepEnvelope() {
	if this.envelopeStart {
		this.envelopeVolume = 15
		this.envelopeValue = this.envelopePeriod
		this.envelopeStart = false
	} else if this.envelopeValue > 0 {
		this.envelopeValue--
	} else {
		if this.envelopeVolume > 0 {
			this.envelopeVolume--
		} else if this.envelopeLoop {
			this.envelopeVolume = 15
		}
		this.envelopeValue = this.envelopePeriod
	}
}

func (this *Noise) stepLength() {
	if this.lengthEnabled && this.lengthValue > 0 {
		this.lengthValue--
	}
}

// DMC

type DMC struct {
	cpu            *CPU
	enabled        bool
	value          byte
	sampleAddress  uint16
	sampleLength   uint16
	currentAddress uint16
	currentLength  uint16
	shiftRegister  byte
	bitCount       byte
	tickPeriod     byte
	tickValue      byte
	loop           bool
	irq            bool
}

func (this *DMC) writeControl(value byte) {
	this.irq = value&0x80 == 0x80
	this.loop = value&0x40 == 0x40
	this.tickPeriod = dmcTable[value&0x0F]
}

func (this *DMC) writeValue(value byte) {
	this.value = value & 0x7F
}

func (this *DMC) writeAddress(value byte) {
	// Sample address = %11AAAAAA.AA000000
	this.sampleAddress = 0xC000 | (uint16(value) << 6)
}

func (this *DMC) writeLength(value byte) {
	// Sample length = %0000LLLL.LLLL0001
	this.sampleLength = (uint16(value) << 4) | 1
}

func (this *DMC) restart() {
	this.currentAddress = this.sampleAddress
	this.currentLength = this.sampleLength
}

func (this *DMC) stepTimer() {
	if !this.enabled {
		return
	}
	this.stepReader()
	if this.tickValue == 0 {
		this.tickValue = this.tickPeriod
		this.stepShifter()
	} else {
		this.tickValue--
	}
}

func (this *DMC) stepReader() {
	if this.currentLength > 0 && this.bitCount == 0 {
		this.cpu.stall += 4
		this.shiftRegister = this.cpu.Read(this.currentAddress)
		this.bitCount = 8
		this.currentAddress++
		if this.currentAddress == 0 {
			this.currentAddress = 0x8000
		}
		this.currentLength--
		if this.currentLength == 0 && this.loop {
			this.restart()
		}
	}
}

func (this *DMC) stepShifter() {
	if this.bitCount == 0 {
		return
	}
	if this.shiftRegister&1 == 1 {
		if this.value <= 125 {
			this.value += 2
		}
	} else {
		if this.value >= 2 {
			this.value -= 2
		}
	}
	this.shiftRegister >>= 1
	this.bitCount--
}

func (this *DMC) output() byte {
	return this.value
}
