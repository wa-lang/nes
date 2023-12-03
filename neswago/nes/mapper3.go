package nes

type Mapper3 struct {
	chrBank  int
	prgBank1 int
	prgBank2 int
}

func NewMapper3() Mapper {
	prgBanks := len(Cartridge_PRG) / 0x4000
	return &Mapper3{0, 0, prgBanks - 1}
}

func (this *Mapper3) Step() {
}

func (this *Mapper3) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		index := this.chrBank*0x2000 + int(address)
		return Cartridge_CHR[index]
	case address >= 0xC000:
		index := this.prgBank2*0x4000 + int(address-0xC000)
		return Cartridge_PRG[index]
	case address >= 0x8000:
		index := this.prgBank1*0x4000 + int(address-0x8000)
		return Cartridge_PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return Cartridge_SRAM[index]
	default:
		log_Fatalf("unhandled mapper3 read at address: 0x%04X", address)
	}
	return 0
}

func (this *Mapper3) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		index := this.chrBank*0x2000 + int(address)
		Cartridge_CHR[index] = value
	case address >= 0x8000:
		this.chrBank = int(value & 3)
	case address >= 0x6000:
		index := int(address) - 0x6000
		Cartridge_SRAM[index] = value
	default:
		log_Fatalf("unhandled mapper3 write at address: 0x%04X", address)
	}
}
