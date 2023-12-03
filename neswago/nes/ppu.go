package nes

import (
	"image"
)

var PPU_Cycle int    // 0-340
var PPU_ScanLine int // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
var PPU_Frame uint64 // frame counter
var PPU_paletteData [32]byte

// storage variables
var PPU_nameTableData [2048]byte
var PPU_oamData [256]byte
var PPU_front *image.RGBA
var PPU_back *image.RGBA

// PPU registers
var PPU_v uint16 // current vram address (15 bit)
var PPU_t uint16 // temporary vram address (15 bit)
var PPU_x byte   // fine x scroll (3 bit)
var PPU_w byte   // write toggle (1 bit)
var PPU_f byte   // even/odd frame flag (1 bit)

var PPU_register byte

// NMI flags
var PPU_nmiOccurred bool
var PPU_nmiOutput bool
var PPU_nmiPrevious bool
var PPU_nmiDelay byte

// background temporary variables
var PPU_nameTableByte byte
var PPU_attributeTableByte byte
var PPU_lowTileByte byte
var PPU_highTileByte byte
var PPU_tileData uint64

// sprite temporary variables
var PPU_spriteCount int
var PPU_spritePatterns [8]uint32
var PPU_spritePositions [8]byte
var PPU_spritePriorities [8]byte
var PPU_spriteIndexes [8]byte

// $2000 PPUCTRL
var PPU_flagNameTable byte       // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
var PPU_flagIncrement byte       // 0: add 1; 1: add 32
var PPU_flagSpriteTable byte     // 0: $0000; 1: $1000; ignored in 8x16 mode
var PPU_flagBackgroundTable byte // 0: $0000; 1: $1000
var PPU_flagSpriteSize byte      // 0: 8x8; 1: 8x16
var PPU_flagMasterSlave byte     // 0: read EXT; 1: write EXT

// $2001 PPUMASK
var PPU_flagGrayscale byte          // 0: color; 1: grayscale
var PPU_flagShowLeftBackground byte // 0: hide; 1: show
var PPU_flagShowLeftSprites byte    // 0: hide; 1: show
var PPU_flagShowBackground byte     // 0: hide; 1: show
var PPU_flagShowSprites byte        // 0: hide; 1: show
var PPU_flagRedTint byte            // 0: normal; 1: emphasized
var PPU_flagGreenTint byte          // 0: normal; 1: emphasized
var PPU_flagBlueTint byte           // 0: normal; 1: emphasized

// $2002 PPUSTATUS
var PPU_flagSpriteZeroHit byte
var PPU_flagSpriteOverflow byte

// $2003 OAMADDR
var PPU_oamAddress byte

// $2007 PPUDATA
var PPU_bufferedData byte // for buffered reads

//var file_log *os.File

func PPU_initNesPPU() {
	PPU_front = image.NewRGBA(image.Rect(0, 0, 256, 240))
	PPU_back = image.NewRGBA(image.Rect(0, 0, 256, 240))
	PPU_Reset()
}

func PPU_Reset() {
	PPU_Cycle = 340
	PPU_ScanLine = 240
	PPU_Frame = 0
	PPU_writeControl(0)
	PPU_writeMask(0)
	PPU_writeOAMAddress(0)
}

func PPU_readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return PPU_paletteData[address]
}

func PPU_writePalette(address uint16, value byte) {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	PPU_paletteData[address] = value
}

func PPU_readRegister(address uint16) byte {
	switch address {
	case 0x2002:
		return PPU_readStatus()
	case 0x2004:
		return PPU_readOAMData()
	case 0x2007:
		return PPU_readData()
	}
	return 0
}

func PPU_writeRegister(address uint16, value byte) {
	PPU_register = value
	switch address {
	case 0x2000:
		PPU_writeControl(value)
	case 0x2001:
		PPU_writeMask(value)
	case 0x2003:
		PPU_writeOAMAddress(value)
	case 0x2004:
		PPU_writeOAMData(value)
	case 0x2005:
		PPU_writeScroll(value)
	case 0x2006:
		PPU_writeAddress(value)
	case 0x2007:
		PPU_writeData(value)
	case 0x4014:
		PPU_writeDMA(value)
	}
}

// $2000: PPUCTRL
func PPU_writeControl(value byte) {
	PPU_flagNameTable = (value >> 0) & 3
	PPU_flagIncrement = (value >> 2) & 1
	PPU_flagSpriteTable = (value >> 3) & 1
	PPU_flagBackgroundTable = (value >> 4) & 1
	PPU_flagSpriteSize = (value >> 5) & 1
	PPU_flagMasterSlave = (value >> 6) & 1

	PPU_nmiOutput = (value>>7)&1 == 1

	PPU_nmiChange()
	// t: ....BA.. ........ = d: ......BA
	PPU_t = (PPU_t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
}

// $2001: PPUMASK
func PPU_writeMask(value byte) {
	PPU_flagGrayscale = (value >> 0) & 1
	PPU_flagShowLeftBackground = (value >> 1) & 1
	PPU_flagShowLeftSprites = (value >> 2) & 1
	PPU_flagShowBackground = (value >> 3) & 1
	PPU_flagShowSprites = (value >> 4) & 1
	PPU_flagRedTint = (value >> 5) & 1
	PPU_flagGreenTint = (value >> 6) & 1
	PPU_flagBlueTint = (value >> 7) & 1
}

// $2002: PPUSTATUS
func PPU_readStatus() byte {
	result := PPU_register & 0x1F
	result |= PPU_flagSpriteOverflow << 5
	result |= PPU_flagSpriteZeroHit << 6
	if PPU_nmiOccurred {
		result |= 1 << 7
	}
	PPU_nmiOccurred = false
	PPU_nmiChange()
	// w:                   = 0
	PPU_w = 0
	return result
}

// $2003: OAMADDR
func PPU_writeOAMAddress(value byte) {
	PPU_oamAddress = value
}

// $2004: OAMDATA (read)
func PPU_readOAMData() byte {
	data := PPU_oamData[PPU_oamAddress]
	if (PPU_oamAddress & 0x03) == 0x02 {
		data = data & 0xE3
	}
	return data
}

// $2004: OAMDATA (write)
func PPU_writeOAMData(value byte) {
	PPU_oamData[PPU_oamAddress] = value
	PPU_oamAddress++
}

// $2005: PPUSCROLL
func PPU_writeScroll(value byte) {
	if PPU_w == 0 {
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		PPU_t = (PPU_t & 0xFFE0) | (uint16(value) >> 3)
		PPU_x = value & 0x07
		PPU_w = 1
	} else {
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		PPU_t = (PPU_t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
		PPU_t = (PPU_t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		PPU_w = 0
	}
}

// $2006: PPUADDR
func PPU_writeAddress(value byte) {
	if PPU_w == 0 {
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		PPU_t = (PPU_t & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		PPU_w = 1
	} else {
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		PPU_t = (PPU_t & 0xFF00) | uint16(value)
		PPU_v = PPU_t
		PPU_w = 0
	}
}

// $2007: PPUDATA (read)
func PPU_readData() byte {
	value := PPUMemory_Read(PPU_v)
	// emulate buffered reads
	if PPU_v%0x4000 < 0x3F00 {
		buffered := PPU_bufferedData
		PPU_bufferedData = value
		value = buffered
	} else {
		PPU_bufferedData = PPUMemory_Read(PPU_v - 0x1000)
	}
	// increment address
	if PPU_flagIncrement == 0 {
		PPU_v += 1
	} else {
		PPU_v += 32
	}
	return value
}

// $2007: PPUDATA (write)
func PPU_writeData(value byte) {
	PPUMemory_Write(PPU_v, value)
	if PPU_flagIncrement == 0 {
		PPU_v += 1
	} else {
		PPU_v += 32
	}
}

// $4014: OAMDMA
func PPU_writeDMA(value byte) {
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		PPU_oamData[PPU_oamAddress] = CPUMemory_Read(address)
		PPU_oamAddress++
		address++
	}
	CPU_stall += 513
	if CPU_Cycles%2 == 1 {
		CPU_stall++
	}
}

// NTSC Timing Helper Functions

func PPU_incrementX() {
	// increment hori(v)
	// if coarse X == 31
	if PPU_v&0x001F == 31 {
		// coarse X = 0
		PPU_v &= 0xFFE0
		// switch horizontal nametable
		PPU_v ^= 0x0400
	} else {
		// increment coarse X
		PPU_v++
	}
}

func PPU_incrementY() {
	// increment vert(v)
	// if fine Y < 7
	if PPU_v&0x7000 != 0x7000 {
		// increment fine Y
		PPU_v += 0x1000
	} else {
		// fine Y = 0
		PPU_v &= 0x8FFF
		// let y = coarse Y
		y := (PPU_v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			PPU_v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		PPU_v = (PPU_v & 0xFC1F) | (y << 5)
	}
}

func PPU_copyX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	PPU_v = (PPU_v & 0xFBE0) | (PPU_t & 0x041F)
}

func PPU_copyY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	PPU_v = (PPU_v & 0x841F) | (PPU_t & 0x7BE0)
}

func PPU_nmiChange() {
	nmi := PPU_nmiOutput && PPU_nmiOccurred
	if nmi && !PPU_nmiPrevious {
		// TODO: this fixes some games but the delay shouldn't have to be so
		// long, so the timings are off somewhere
		PPU_nmiDelay = 15
	}
	PPU_nmiPrevious = nmi
}

func PPU_setVerticalBlank() {
	PPU_front, PPU_back = PPU_back, PPU_front
	PPU_nmiOccurred = true
	PPU_nmiChange()
}

func PPU_clearVerticalBlank() {
	PPU_nmiOccurred = false
	PPU_nmiChange()
}

func PPU_fetchNameTableByte() {
	v := PPU_v
	address := 0x2000 | (v & 0x0FFF)
	PPU_nameTableByte = PPUMemory_Read(address)
}

func PPU_fetchAttributeTableByte() {
	v := PPU_v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	PPU_attributeTableByte = ((PPUMemory_Read(address) >> shift) & 3) << 2
}

func PPU_fetchLowTileByte() {
	fineY := (PPU_v >> 12) & 7
	table := PPU_flagBackgroundTable
	tile := PPU_nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	PPU_lowTileByte = PPUMemory_Read(address)
}

func PPU_fetchHighTileByte() {
	fineY := (PPU_v >> 12) & 7
	table := PPU_flagBackgroundTable
	tile := PPU_nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	PPU_highTileByte = PPUMemory_Read(address + 8)
}

func PPU_storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		a := PPU_attributeTableByte
		p1 := (PPU_lowTileByte & 0x80) >> 7
		p2 := (PPU_highTileByte & 0x80) >> 6
		PPU_lowTileByte <<= 1
		PPU_highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	PPU_tileData |= uint64(data)
}

func PPU_fetchTileData() uint32 {
	return uint32(PPU_tileData >> 32)
}

func PPU_backgroundPixel() byte {
	if PPU_flagShowBackground == 0 {
		return 0
	}
	data := PPU_fetchTileData() >> ((7 - PPU_x) * 4)
	return byte(data & 0x0F)
}

func PPU_spritePixel() (byte, byte) {
	if PPU_flagShowSprites == 0 {
		return 0, 0
	}
	for i := 0; i < PPU_spriteCount; i++ {
		offset := (PPU_Cycle - 1) - int(PPU_spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((PPU_spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func PPU_renderPixel() {
	x := PPU_Cycle - 1
	y := PPU_ScanLine
	background := PPU_backgroundPixel()
	i, sprite := PPU_spritePixel()
	if x < 8 && PPU_flagShowLeftBackground == 0 {
		background = 0
	}
	if x < 8 && PPU_flagShowLeftSprites == 0 {
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
		if PPU_spriteIndexes[i] == 0 && x < 255 {
			PPU_flagSpriteZeroHit = 1
		}
		if PPU_spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = background
		}
	}
	c := Palette[PPU_readPalette(uint16(color))%64]
	PPU_back.SetRGBA(x, y, c)
}

func PPU_fetchSpritePattern(i, row int) uint32 {
	tile := PPU_oamData[i*4+1]
	attributes := PPU_oamData[i*4+2]
	var address uint16
	if PPU_flagSpriteSize == 0 {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		table := PPU_flagSpriteTable
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
	lowTileByte := PPUMemory_Read(address)
	highTileByte := PPUMemory_Read(address + 8)
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

func PPU_evaluateSprites() {
	var h int
	if PPU_flagSpriteSize == 0 {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := PPU_oamData[i*4+0]
		a := PPU_oamData[i*4+2]
		x := PPU_oamData[i*4+3]
		row := PPU_ScanLine - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			PPU_spritePatterns[count] = PPU_fetchSpritePattern(i, row)
			PPU_spritePositions[count] = x
			PPU_spritePriorities[count] = (a >> 5) & 1
			PPU_spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		PPU_flagSpriteOverflow = 1
	}
	PPU_spriteCount = count
}

// tick updates Cycle, ScanLine and Frame counters
func PPU_tick() {
	if PPU_nmiDelay > 0 {
		PPU_nmiDelay--
		if PPU_nmiDelay == 0 && PPU_nmiOutput && PPU_nmiOccurred {
			CPU_triggerNMI()
		}
	}

	if PPU_flagShowBackground != 0 || PPU_flagShowSprites != 0 {
		if PPU_f == 1 && PPU_ScanLine == 261 && PPU_Cycle == 339 {
			PPU_Cycle = 0
			PPU_ScanLine = 0
			PPU_Frame++
			PPU_f ^= 1
			return
		}
	}
	PPU_Cycle++
	if PPU_Cycle > 340 {
		PPU_Cycle = 0
		PPU_ScanLine++
		if PPU_ScanLine > 261 {
			PPU_ScanLine = 0
			PPU_Frame++
			PPU_f ^= 1
		}
	}
}

// Step executes a single PPU cycle
func PPU_Step() {
	PPU_tick()

	renderingEnabled := PPU_flagShowBackground != 0 || PPU_flagShowSprites != 0
	preLine := PPU_ScanLine == 261
	visibleLine := PPU_ScanLine < 240
	// postLine := PPU_ScanLine == 240
	renderLine := preLine || visibleLine
	preFetchCycle := PPU_Cycle >= 321 && PPU_Cycle <= 336
	visibleCycle := PPU_Cycle >= 1 && PPU_Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	// background logic
	if renderingEnabled {
		if visibleLine && visibleCycle {
			PPU_renderPixel()
		}
		if renderLine && fetchCycle {
			PPU_tileData <<= 4
			switch PPU_Cycle % 8 {
			case 1:
				PPU_fetchNameTableByte()
			case 3:
				PPU_fetchAttributeTableByte()
			case 5:
				PPU_fetchLowTileByte()
			case 7:
				PPU_fetchHighTileByte()
			case 0:
				PPU_storeTileData()
			}
		}
		if preLine && PPU_Cycle >= 280 && PPU_Cycle <= 304 {
			PPU_copyY()
		}
		if renderLine {
			if fetchCycle && PPU_Cycle%8 == 0 {
				PPU_incrementX()
			}
			if PPU_Cycle == 256 {
				PPU_incrementY()
			}
			if PPU_Cycle == 257 {
				PPU_copyX()
			}
		}
	}

	// sprite logic
	if renderingEnabled {
		if PPU_Cycle == 257 {
			if visibleLine {
				PPU_evaluateSprites()
			} else {
				PPU_spriteCount = 0
			}
		}
	}

	// vblank logic
	if PPU_ScanLine == 241 && PPU_Cycle == 1 {
		PPU_setVerticalBlank()
	}
	if preLine && PPU_Cycle == 1 {
		PPU_clearVerticalBlank()
		PPU_flagSpriteZeroHit = 0
		PPU_flagSpriteOverflow = 0
	}
}
