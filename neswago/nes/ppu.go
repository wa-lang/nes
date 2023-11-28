package nes

import (
	"image"
)

type PPU struct {
	Memory           // memory interface
	console *Console // reference to parent object

	Cycle    int    // 0-340
	ScanLine int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64 // frame counter

	// storage variables
	paletteData   [32]byte
	nameTableData [2048]byte
	oamData       [256]byte
	front         *image.RGBA
	back          *image.RGBA

	// PPU registers
	v uint16 // current vram address (15 bit)
	t uint16 // temporary vram address (15 bit)
	x byte   // fine x scroll (3 bit)
	w byte   // write toggle (1 bit)
	f byte   // even/odd frame flag (1 bit)

	register byte

	// NMI flags
	nmiOccurred bool
	nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	// background temporary variables
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
	tileData           uint64

	// sprite temporary variables
	spriteCount      int
	spritePatterns   [8]uint32
	spritePositions  [8]byte
	spritePriorities [8]byte
	spriteIndexes    [8]byte

	// $2000 PPUCTRL
	flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	flagIncrement       byte // 0: add 1; 1: add 32
	flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
	flagBackgroundTable byte // 0: $0000; 1: $1000
	flagSpriteSize      byte // 0: 8x8; 1: 8x16
	flagMasterSlave     byte // 0: read EXT; 1: write EXT

	// $2001 PPUMASK
	flagGrayscale          byte // 0: color; 1: grayscale
	flagShowLeftBackground byte // 0: hide; 1: show
	flagShowLeftSprites    byte // 0: hide; 1: show
	flagShowBackground     byte // 0: hide; 1: show
	flagShowSprites        byte // 0: hide; 1: show
	flagRedTint            byte // 0: normal; 1: emphasized
	flagGreenTint          byte // 0: normal; 1: emphasized
	flagBlueTint           byte // 0: normal; 1: emphasized

	// $2002 PPUSTATUS
	flagSpriteZeroHit  byte
	flagSpriteOverflow byte

	// $2003 OAMADDR
	oamAddress byte

	// $2007 PPUDATA
	bufferedData byte // for buffered reads
}

func NewPPU(console *Console) *PPU {
	this := &PPU{Memory: NewPPUMemory(console), console: console}
	this.front = image.NewRGBA(image.Rect(0, 0, 256, 240))
	this.back = image.NewRGBA(image.Rect(0, 0, 256, 240))
	this.Reset()
	return this
}

func (this *PPU) Reset() {
	this.Cycle = 340
	this.ScanLine = 240
	this.Frame = 0
	this.writeControl(0)
	this.writeMask(0)
	this.writeOAMAddress(0)
}

func (this *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return this.paletteData[address]
}

func (this *PPU) writePalette(address uint16, value byte) {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	this.paletteData[address] = value
}

func (this *PPU) readRegister(address uint16) byte {
	switch address {
	case 0x2002:
		return this.readStatus()
	case 0x2004:
		return this.readOAMData()
	case 0x2007:
		return this.readData()
	}
	return 0
}

func (this *PPU) writeRegister(address uint16, value byte) {
	this.register = value
	switch address {
	case 0x2000:
		this.writeControl(value)
	case 0x2001:
		this.writeMask(value)
	case 0x2003:
		this.writeOAMAddress(value)
	case 0x2004:
		this.writeOAMData(value)
	case 0x2005:
		this.writeScroll(value)
	case 0x2006:
		this.writeAddress(value)
	case 0x2007:
		this.writeData(value)
	case 0x4014:
		this.writeDMA(value)
	}
}

// $2000: PPUCTRL
func (this *PPU) writeControl(value byte) {
	this.flagNameTable = (value >> 0) & 3
	this.flagIncrement = (value >> 2) & 1
	this.flagSpriteTable = (value >> 3) & 1
	this.flagBackgroundTable = (value >> 4) & 1
	this.flagSpriteSize = (value >> 5) & 1
	this.flagMasterSlave = (value >> 6) & 1
	this.nmiOutput = (value>>7)&1 == 1
	this.nmiChange()
	// t: ....BA.. ........ = d: ......BA
	this.t = (this.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
}

// $2001: PPUMASK
func (this *PPU) writeMask(value byte) {
	this.flagGrayscale = (value >> 0) & 1
	this.flagShowLeftBackground = (value >> 1) & 1
	this.flagShowLeftSprites = (value >> 2) & 1
	this.flagShowBackground = (value >> 3) & 1
	this.flagShowSprites = (value >> 4) & 1
	this.flagRedTint = (value >> 5) & 1
	this.flagGreenTint = (value >> 6) & 1
	this.flagBlueTint = (value >> 7) & 1
}

// $2002: PPUSTATUS
func (this *PPU) readStatus() byte {
	result := this.register & 0x1F
	result |= this.flagSpriteOverflow << 5
	result |= this.flagSpriteZeroHit << 6
	if this.nmiOccurred {
		result |= 1 << 7
	}
	this.nmiOccurred = false
	this.nmiChange()
	// w:                   = 0
	this.w = 0
	return result
}

// $2003: OAMADDR
func (this *PPU) writeOAMAddress(value byte) {
	this.oamAddress = value
}

// $2004: OAMDATA (read)
func (this *PPU) readOAMData() byte {
	data := this.oamData[this.oamAddress]
	if (this.oamAddress & 0x03) == 0x02 {
		data = data & 0xE3
	}
	return data
}

// $2004: OAMDATA (write)
func (this *PPU) writeOAMData(value byte) {
	this.oamData[this.oamAddress] = value
	this.oamAddress++
}

// $2005: PPUSCROLL
func (this *PPU) writeScroll(value byte) {
	if this.w == 0 {
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		this.t = (this.t & 0xFFE0) | (uint16(value) >> 3)
		this.x = value & 0x07
		this.w = 1
	} else {
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		this.t = (this.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
		this.t = (this.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		this.w = 0
	}
}

// $2006: PPUADDR
func (this *PPU) writeAddress(value byte) {
	if this.w == 0 {
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		this.t = (this.t & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		this.w = 1
	} else {
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		this.t = (this.t & 0xFF00) | uint16(value)
		this.v = this.t
		this.w = 0
	}
}

// $2007: PPUDATA (read)
func (this *PPU) readData() byte {
	value := this.Read(this.v)
	// emulate buffered reads
	if this.v%0x4000 < 0x3F00 {
		buffered := this.bufferedData
		this.bufferedData = value
		value = buffered
	} else {
		this.bufferedData = this.Read(this.v - 0x1000)
	}
	// increment address
	if this.flagIncrement == 0 {
		this.v += 1
	} else {
		this.v += 32
	}
	return value
}

// $2007: PPUDATA (write)
func (this *PPU) writeData(value byte) {
	this.Write(this.v, value)
	if this.flagIncrement == 0 {
		this.v += 1
	} else {
		this.v += 32
	}
}

// $4014: OAMDMA
func (this *PPU) writeDMA(value byte) {
	cpu := this.console.CPU
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		this.oamData[this.oamAddress] = cpu.Read(address)
		this.oamAddress++
		address++
	}
	cpu.stall += 513
	if cpu.Cycles%2 == 1 {
		cpu.stall++
	}
}

// NTSC Timing Helper Functions

func (this *PPU) incrementX() {
	// increment hori(v)
	// if coarse X == 31
	if this.v&0x001F == 31 {
		// coarse X = 0
		this.v &= 0xFFE0
		// switch horizontal nametable
		this.v ^= 0x0400
	} else {
		// increment coarse X
		this.v++
	}
}

func (this *PPU) incrementY() {
	// increment vert(v)
	// if fine Y < 7
	if this.v&0x7000 != 0x7000 {
		// increment fine Y
		this.v += 0x1000
	} else {
		// fine Y = 0
		this.v &= 0x8FFF
		// let y = coarse Y
		y := (this.v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			this.v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		this.v = (this.v & 0xFC1F) | (y << 5)
	}
}

func (this *PPU) copyX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	this.v = (this.v & 0xFBE0) | (this.t & 0x041F)
}

func (this *PPU) copyY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	this.v = (this.v & 0x841F) | (this.t & 0x7BE0)
}

func (this *PPU) nmiChange() {
	nmi := this.nmiOutput && this.nmiOccurred
	if nmi && !this.nmiPrevious {
		// TODO: this fixes some games but the delay shouldn't have to be so
		// long, so the timings are off somewhere
		this.nmiDelay = 15
	}
	this.nmiPrevious = nmi
}

func (this *PPU) setVerticalBlank() {
	this.front, this.back = this.back, this.front
	this.nmiOccurred = true
	this.nmiChange()
}

func (this *PPU) clearVerticalBlank() {
	this.nmiOccurred = false
	this.nmiChange()
}

func (this *PPU) fetchNameTableByte() {
	v := this.v
	address := 0x2000 | (v & 0x0FFF)
	this.nameTableByte = this.Read(address)
}

func (this *PPU) fetchAttributeTableByte() {
	v := this.v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	this.attributeTableByte = ((this.Read(address) >> shift) & 3) << 2
}

func (this *PPU) fetchLowTileByte() {
	fineY := (this.v >> 12) & 7
	table := this.flagBackgroundTable
	tile := this.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	this.lowTileByte = this.Read(address)
}

func (this *PPU) fetchHighTileByte() {
	fineY := (this.v >> 12) & 7
	table := this.flagBackgroundTable
	tile := this.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	this.highTileByte = this.Read(address + 8)
}

func (this *PPU) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		a := this.attributeTableByte
		p1 := (this.lowTileByte & 0x80) >> 7
		p2 := (this.highTileByte & 0x80) >> 6
		this.lowTileByte <<= 1
		this.highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	this.tileData |= uint64(data)
}

func (this *PPU) fetchTileData() uint32 {
	return uint32(this.tileData >> 32)
}

func (this *PPU) backgroundPixel() byte {
	if this.flagShowBackground == 0 {
		return 0
	}
	data := this.fetchTileData() >> ((7 - this.x) * 4)
	return byte(data & 0x0F)
}

func (this *PPU) spritePixel() (byte, byte) {
	if this.flagShowSprites == 0 {
		return 0, 0
	}
	for i := 0; i < this.spriteCount; i++ {
		offset := (this.Cycle - 1) - int(this.spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((this.spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func (this *PPU) renderPixel() {
	x := this.Cycle - 1
	y := this.ScanLine
	background := this.backgroundPixel()
	i, sprite := this.spritePixel()
	if x < 8 && this.flagShowLeftBackground == 0 {
		background = 0
	}
	if x < 8 && this.flagShowLeftSprites == 0 {
		sprite = 0
	}
	b := background%4 != 0
	s := sprite%4 != 0
	var color byte
	if !b && !s {
		color = 0
	} else if !b && s {
		color = sprite | 0x10
	} else if b && !s {
		color = background
	} else {
		if this.spriteIndexes[i] == 0 && x < 255 {
			this.flagSpriteZeroHit = 1
		}
		if this.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = background
		}
	}
	c := Palette[this.readPalette(uint16(color))%64]
	this.back.SetRGBA(x, y, c)
}

func (this *PPU) fetchSpritePattern(i, row int) uint32 {
	tile := this.oamData[i*4+1]
	attributes := this.oamData[i*4+2]
	var address uint16
	if this.flagSpriteSize == 0 {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		table := this.flagSpriteTable
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := this.Read(address)
	highTileByte := this.Read(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (this *PPU) evaluateSprites() {
	var h int
	if this.flagSpriteSize == 0 {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := this.oamData[i*4+0]
		a := this.oamData[i*4+2]
		x := this.oamData[i*4+3]
		row := this.ScanLine - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			this.spritePatterns[count] = this.fetchSpritePattern(i, row)
			this.spritePositions[count] = x
			this.spritePriorities[count] = (a >> 5) & 1
			this.spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		this.flagSpriteOverflow = 1
	}
	this.spriteCount = count
}

// tick updates Cycle, ScanLine and Frame counters
func (this *PPU) tick() {
	if this.nmiDelay > 0 {
		this.nmiDelay--
		if this.nmiDelay == 0 && this.nmiOutput && this.nmiOccurred {
			this.console.CPU.triggerNMI()
		}
	}

	if this.flagShowBackground != 0 || this.flagShowSprites != 0 {
		if this.f == 1 && this.ScanLine == 261 && this.Cycle == 339 {
			this.Cycle = 0
			this.ScanLine = 0
			this.Frame++
			this.f ^= 1
			return
		}
	}
	this.Cycle++
	if this.Cycle > 340 {
		this.Cycle = 0
		this.ScanLine++
		if this.ScanLine > 261 {
			this.ScanLine = 0
			this.Frame++
			this.f ^= 1
		}
	}
}

// Step executes a single PPU cycle
func (this *PPU) Step() {
	this.tick()

	renderingEnabled := this.flagShowBackground != 0 || this.flagShowSprites != 0
	preLine := this.ScanLine == 261
	visibleLine := this.ScanLine < 240
	// postLine := this.ScanLine == 240
	renderLine := preLine || visibleLine
	preFetchCycle := this.Cycle >= 321 && this.Cycle <= 336
	visibleCycle := this.Cycle >= 1 && this.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	// background logic
	if renderingEnabled {
		if visibleLine && visibleCycle {
			this.renderPixel()
		}
		if renderLine && fetchCycle {
			this.tileData <<= 4
			switch this.Cycle % 8 {
			case 1:
				this.fetchNameTableByte()
			case 3:
				this.fetchAttributeTableByte()
			case 5:
				this.fetchLowTileByte()
			case 7:
				this.fetchHighTileByte()
			case 0:
				this.storeTileData()
			}
		}
		if preLine && this.Cycle >= 280 && this.Cycle <= 304 {
			this.copyY()
		}
		if renderLine {
			if fetchCycle && this.Cycle%8 == 0 {
				this.incrementX()
			}
			if this.Cycle == 256 {
				this.incrementY()
			}
			if this.Cycle == 257 {
				this.copyX()
			}
		}
	}

	// sprite logic
	if renderingEnabled {
		if this.Cycle == 257 {
			if visibleLine {
				this.evaluateSprites()
			} else {
				this.spriteCount = 0
			}
		}
	}

	// vblank logic
	if this.ScanLine == 241 && this.Cycle == 1 {
		this.setVerticalBlank()
	}
	if preLine && this.Cycle == 1 {
		this.clearVerticalBlank()
		this.flagSpriteZeroHit = 0
		this.flagSpriteOverflow = 0
	}
}
