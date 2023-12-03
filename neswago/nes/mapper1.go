package nes

type Mapper1 struct {
	shiftRegister byte
	control       byte
	prgMode       byte
	chrMode       byte
	prgBank       byte
	chrBank0      byte
	chrBank1      byte
	prgOffsets    [2]int
	chrOffsets    [2]int
}

func NewMapper1() Mapper {
	this := &Mapper1{}
	this.shiftRegister = 0x10
	this.prgOffsets[1] = this.prgBankOffset(-1)
	return this
}

func (this *Mapper1) Step() {
}

func (this *Mapper1) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		bank := address / 0x1000
		offset := address % 0x1000
		return Cartridge_CHR[this.chrOffsets[bank]+int(offset)]
	case address >= 0x8000:
		address = address - 0x8000
		bank := address / 0x4000
		offset := address % 0x4000
		return Cartridge_PRG[this.prgOffsets[bank]+int(offset)]
	case address >= 0x6000:
		return Cartridge_SRAM[int(address)-0x6000]
	default:
		log_Fatalf("unhandled mapper1 read at address: 0x%04X", address)
	}
	return 0
}

func (this *Mapper1) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		bank := address / 0x1000
		offset := address % 0x1000
		Cartridge_CHR[this.chrOffsets[bank]+int(offset)] = value
	case address >= 0x8000:
		this.loadRegister(address, value)
	case address >= 0x6000:
		Cartridge_SRAM[int(address)-0x6000] = value
	default:
		log_Fatalf("unhandled mapper1 write at address: 0x%04X", address)
	}
}

func (this *Mapper1) loadRegister(address uint16, value byte) {
	if value&0x80 == 0x80 {
		this.shiftRegister = 0x10
		this.writeControl(this.control | 0x0C)
	} else {
		complete := this.shiftRegister&1 == 1
		this.shiftRegister >>= 1
		this.shiftRegister |= (value & 1) << 4
		if complete {
			this.writeRegister(address, this.shiftRegister)
			this.shiftRegister = 0x10
		}
	}
}

func (this *Mapper1) writeRegister(address uint16, value byte) {
	switch {
	case address <= 0x9FFF:
		this.writeControl(value)
	case address <= 0xBFFF:
		this.writeCHRBank0(value)
	case address <= 0xDFFF:
		this.writeCHRBank1(value)
	case address <= 0xFFFF:
		this.writePRGBank(value)
	}
}

// Control (internal, $8000-$9FFF)
func (this *Mapper1) writeControl(value byte) {
	this.control = value
	this.chrMode = (value >> 4) & 1
	this.prgMode = (value >> 2) & 3
	mirror := value & 3
	switch mirror {
	case 0:
		Cartridge_Mirror = MirrorSingle0
	case 1:
		Cartridge_Mirror = MirrorSingle1
	case 2:
		Cartridge_Mirror = MirrorVertical
	case 3:
		Cartridge_Mirror = MirrorHorizontal
	}
	this.updateOffsets()
}

// CHR bank 0 (internal, $A000-$BFFF)
func (this *Mapper1) writeCHRBank0(value byte) {
	this.chrBank0 = value
	this.updateOffsets()
}

// CHR bank 1 (internal, $C000-$DFFF)
func (this *Mapper1) writeCHRBank1(value byte) {
	this.chrBank1 = value
	this.updateOffsets()
}

// PRG bank (internal, $E000-$FFFF)
func (this *Mapper1) writePRGBank(value byte) {
	this.prgBank = value & 0x0F
	this.updateOffsets()
}

func (this *Mapper1) prgBankOffset(index int) int {
	if index >= 0x80 {
		index -= 0x100
	}
	index %= len(Cartridge_PRG) / 0x4000
	offset := index * 0x4000
	if offset < 0 {
		offset += len(Cartridge_PRG)
	}
	return offset
}

func (this *Mapper1) chrBankOffset(index int) int {
	if index >= 0x80 {
		index -= 0x100
	}
	index %= len(Cartridge_CHR) / 0x1000
	offset := index * 0x1000
	if offset < 0 {
		offset += len(Cartridge_CHR)
	}
	return offset
}

// PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
//
//	2: fix first bank at $8000 and switch 16 KB bank at $C000;
//	3: fix last bank at $C000 and switch 16 KB bank at $8000)
//
// CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
func (this *Mapper1) updateOffsets() {
	switch this.prgMode {
	case 0, 1:
		this.prgOffsets[0] = this.prgBankOffset(int(this.prgBank & 0xFE))
		this.prgOffsets[1] = this.prgBankOffset(int(this.prgBank | 0x01))
	case 2:
		this.prgOffsets[0] = 0
		this.prgOffsets[1] = this.prgBankOffset(int(this.prgBank))
	case 3:
		this.prgOffsets[0] = this.prgBankOffset(int(this.prgBank))
		this.prgOffsets[1] = this.prgBankOffset(-1)
	}
	switch this.chrMode {
	case 0:
		this.chrOffsets[0] = this.chrBankOffset(int(this.chrBank0 & 0xFE))
		this.chrOffsets[1] = this.chrBankOffset(int(this.chrBank0 | 0x01))
	case 1:
		this.chrOffsets[0] = this.chrBankOffset(int(this.chrBank0))
		this.chrOffsets[1] = this.chrBankOffset(int(this.chrBank1))
	}
}
