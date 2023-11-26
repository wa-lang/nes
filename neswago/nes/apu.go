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
	filterChain FilterChain
}

func NewAPU(console *Console) *APU {
	apu := APU{}
	apu.console = console
	apu.noise.shiftRegister = 1
	apu.pulse1.channel = 1
	apu.pulse2.channel = 2
	apu.framePeriod = 4
	apu.dmc.cpu = console.CPU
	return &apu
}

func (apu *APU) Step() {
	cycle1 := apu.cycle
	apu.cycle++
	cycle2 := apu.cycle
	apu.stepTimer()
	f1 := int(float64(cycle1) / frameCounterRate)
	f2 := int(float64(cycle2) / frameCounterRate)
	if f1 != f2 {
		apu.stepFrameCounter()
	}
	s1 := int(float64(cycle1) / apu.sampleRate)
	s2 := int(float64(cycle2) / apu.sampleRate)
	if s1 != s2 {
		// todo 发送数据
	}
}

// mode 0:    mode 1:       function
// ---------  -----------  -----------------------------
//   - - - f    - - - - -    IRQ (if bit 6 is clear)
//   - l - l    l - l - -    Length counter and sweep
//     e e e e    e e e e -    Envelope and linear counter
func (apu *APU) stepFrameCounter() {
	switch apu.framePeriod {
	case 4:
		apu.frameValue = (apu.frameValue + 1) % 4
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		case 3:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
			apu.fireIRQ()
		}
	case 5:
		apu.frameValue = (apu.frameValue + 1) % 5
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1, 3:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		}
	}
}

func (apu *APU) stepTimer() {
	if apu.cycle%2 == 0 {
		apu.pulse1.stepTimer()
		apu.pulse2.stepTimer()
		apu.noise.stepTimer()
		apu.dmc.stepTimer()
	}
	apu.triangle.stepTimer()
}

func (apu *APU) stepEnvelope() {
	apu.pulse1.stepEnvelope()
	apu.pulse2.stepEnvelope()
	apu.triangle.stepCounter()
	apu.noise.stepEnvelope()
}

func (apu *APU) stepSweep() {
	apu.pulse1.stepSweep()
	apu.pulse2.stepSweep()
}

func (apu *APU) stepLength() {
	apu.pulse1.stepLength()
	apu.pulse2.stepLength()
	apu.triangle.stepLength()
	apu.noise.stepLength()
}

func (apu *APU) fireIRQ() {
	if apu.frameIRQ {
		apu.console.CPU.triggerIRQ()
	}
}

func (apu *APU) readRegister(address uint16) byte {
	switch address {
	case 0x4015:
		return apu.readStatus()
		// default:
		// 	log.Fatalf("unhandled apu register read at address: 0x%04X", address)
	}
	return 0
}

func (apu *APU) writeRegister(address uint16, value byte) {
	switch address {
	case 0x4000:
		apu.pulse1.writeControl(value)
	case 0x4001:
		apu.pulse1.writeSweep(value)
	case 0x4002:
		apu.pulse1.writeTimerLow(value)
	case 0x4003:
		apu.pulse1.writeTimerHigh(value)
	case 0x4004:
		apu.pulse2.writeControl(value)
	case 0x4005:
		apu.pulse2.writeSweep(value)
	case 0x4006:
		apu.pulse2.writeTimerLow(value)
	case 0x4007:
		apu.pulse2.writeTimerHigh(value)
	case 0x4008:
		apu.triangle.writeControl(value)
	case 0x4009:
	case 0x4010:
		apu.dmc.writeControl(value)
	case 0x4011:
		apu.dmc.writeValue(value)
	case 0x4012:
		apu.dmc.writeAddress(value)
	case 0x4013:
		apu.dmc.writeLength(value)
	case 0x400A:
		apu.triangle.writeTimerLow(value)
	case 0x400B:
		apu.triangle.writeTimerHigh(value)
	case 0x400C:
		apu.noise.writeControl(value)
	case 0x400D:
	case 0x400E:
		apu.noise.writePeriod(value)
	case 0x400F:
		apu.noise.writeLength(value)
	case 0x4015:
		apu.writeControl(value)
	case 0x4017:
		apu.writeFrameCounter(value)
		// default:
		// 	log.Fatalf("unhandled apu register write at address: 0x%04X", address)
	}
}

func (apu *APU) readStatus() byte {
	var result byte
	if apu.pulse1.lengthValue > 0 {
		result |= 1
	}
	if apu.pulse2.lengthValue > 0 {
		result |= 2
	}
	if apu.triangle.lengthValue > 0 {
		result |= 4
	}
	if apu.noise.lengthValue > 0 {
		result |= 8
	}
	if apu.dmc.currentLength > 0 {
		result |= 16
	}
	return result
}

func (apu *APU) writeControl(value byte) {
	apu.pulse1.enabled = value&1 == 1
	apu.pulse2.enabled = value&2 == 2
	apu.triangle.enabled = value&4 == 4
	apu.noise.enabled = value&8 == 8
	apu.dmc.enabled = value&16 == 16
	if !apu.pulse1.enabled {
		apu.pulse1.lengthValue = 0
	}
	if !apu.pulse2.enabled {
		apu.pulse2.lengthValue = 0
	}
	if !apu.triangle.enabled {
		apu.triangle.lengthValue = 0
	}
	if !apu.noise.enabled {
		apu.noise.lengthValue = 0
	}
	if !apu.dmc.enabled {
		apu.dmc.currentLength = 0
	} else {
		if apu.dmc.currentLength == 0 {
			apu.dmc.restart()
		}
	}
}

func (apu *APU) writeFrameCounter(value byte) {
	apu.framePeriod = 4 + (value>>7)&1
	apu.frameIRQ = (value>>6)&1 == 0
	// apu.frameValue = 0
	if apu.framePeriod == 5 {
		apu.stepEnvelope()
		apu.stepSweep()
		apu.stepLength()
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

func (t *Triangle) writeControl(value byte) {
	t.lengthEnabled = (value>>7)&1 == 0
	t.counterPeriod = value & 0x7F
}

func (t *Triangle) writeTimerLow(value byte) {
	t.timerPeriod = (t.timerPeriod & 0xFF00) | uint16(value)
}

func (t *Triangle) writeTimerHigh(value byte) {
	t.lengthValue = lengthTable[value>>3]
	t.timerPeriod = (t.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
	t.timerValue = t.timerPeriod
	t.counterReload = true
}

func (t *Triangle) stepTimer() {
	if t.timerValue == 0 {
		t.timerValue = t.timerPeriod
		if t.lengthValue > 0 && t.counterValue > 0 {
			t.dutyValue = (t.dutyValue + 1) % 32
		}
	} else {
		t.timerValue--
	}
}

func (t *Triangle) stepLength() {
	if t.lengthEnabled && t.lengthValue > 0 {
		t.lengthValue--
	}
}

func (t *Triangle) stepCounter() {
	if t.counterReload {
		t.counterValue = t.counterPeriod
	} else if t.counterValue > 0 {
		t.counterValue--
	}
	if t.lengthEnabled {
		t.counterReload = false
	}
}

func (t *Triangle) output() byte {
	if !t.enabled {
		return 0
	}
	if t.timerPeriod < 3 {
		return 0
	}
	if t.lengthValue == 0 {
		return 0
	}
	if t.counterValue == 0 {
		return 0
	}
	return triangleTable[t.dutyValue]
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

func (n *Noise) writeControl(value byte) {
	n.lengthEnabled = (value>>5)&1 == 0
	n.envelopeLoop = (value>>5)&1 == 1
	n.envelopeEnabled = (value>>4)&1 == 0
	n.envelopePeriod = value & 15
	n.constantVolume = value & 15
	n.envelopeStart = true
}

func (n *Noise) writePeriod(value byte) {
	n.mode = value&0x80 == 0x80
	n.timerPeriod = noiseTable[value&0x0F]
}

func (n *Noise) writeLength(value byte) {
	n.lengthValue = lengthTable[value>>3]
	n.envelopeStart = true
}

func (n *Noise) stepTimer() {
	if n.timerValue == 0 {
		n.timerValue = n.timerPeriod
		var shift byte
		if n.mode {
			shift = 6
		} else {
			shift = 1
		}
		b1 := n.shiftRegister & 1
		b2 := (n.shiftRegister >> shift) & 1
		n.shiftRegister >>= 1
		n.shiftRegister |= (b1 ^ b2) << 14
	} else {
		n.timerValue--
	}
}

func (n *Noise) stepEnvelope() {
	if n.envelopeStart {
		n.envelopeVolume = 15
		n.envelopeValue = n.envelopePeriod
		n.envelopeStart = false
	} else if n.envelopeValue > 0 {
		n.envelopeValue--
	} else {
		if n.envelopeVolume > 0 {
			n.envelopeVolume--
		} else if n.envelopeLoop {
			n.envelopeVolume = 15
		}
		n.envelopeValue = n.envelopePeriod
	}
}

func (n *Noise) stepLength() {
	if n.lengthEnabled && n.lengthValue > 0 {
		n.lengthValue--
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

func (d *DMC) writeControl(value byte) {
	d.irq = value&0x80 == 0x80
	d.loop = value&0x40 == 0x40
	d.tickPeriod = dmcTable[value&0x0F]
}

func (d *DMC) writeValue(value byte) {
	d.value = value & 0x7F
}

func (d *DMC) writeAddress(value byte) {
	// Sample address = %11AAAAAA.AA000000
	d.sampleAddress = 0xC000 | (uint16(value) << 6)
}

func (d *DMC) writeLength(value byte) {
	// Sample length = %0000LLLL.LLLL0001
	d.sampleLength = (uint16(value) << 4) | 1
}

func (d *DMC) restart() {
	d.currentAddress = d.sampleAddress
	d.currentLength = d.sampleLength
}

func (d *DMC) stepTimer() {
	if !d.enabled {
		return
	}
	d.stepReader()
	if d.tickValue == 0 {
		d.tickValue = d.tickPeriod
		d.stepShifter()
	} else {
		d.tickValue--
	}
}

func (d *DMC) stepReader() {
	if d.currentLength > 0 && d.bitCount == 0 {
		d.cpu.stall += 4
		d.shiftRegister = d.cpu.Read(d.currentAddress)
		d.bitCount = 8
		d.currentAddress++
		if d.currentAddress == 0 {
			d.currentAddress = 0x8000
		}
		d.currentLength--
		if d.currentLength == 0 && d.loop {
			d.restart()
		}
	}
}

func (d *DMC) stepShifter() {
	if d.bitCount == 0 {
		return
	}
	if d.shiftRegister&1 == 1 {
		if d.value <= 125 {
			d.value += 2
		}
	} else {
		if d.value >= 2 {
			d.value -= 2
		}
	}
	d.shiftRegister >>= 1
	d.bitCount--
}

func (d *DMC) output() byte {
	return d.value
}
